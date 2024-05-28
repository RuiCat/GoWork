// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

type shaderSources struct {
	Name      string
	SPIRV     []byte
	GLSL100ES []byte
	GLSL150   []byte
	DXBC      []byte
	MetalLibs MetalLibs
	Reflect   Metadata
}

func main() {
	packageName := flag.String("package", "", "specify Go package name")
	workdir := flag.String("work", "", "temporary working directory (default TEMP)")
	shadersDir := flag.String("dir", "shaders", "shaders directory")

	flag.Parse()

	var work WorkDir
	cleanup := func() {}
	if *workdir == "" {
		tempdir, err := os.MkdirTemp("", "shader-convert")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create tempdir: %v\n", err)
			os.Exit(1)
		}
		cleanup = func() { os.RemoveAll(tempdir) }
		defer cleanup()

		work = WorkDir(tempdir)
	} else {
		if abs, err := filepath.Abs(*workdir); err == nil {
			*workdir = abs
		}
		work = WorkDir(*workdir)
	}

	var out bytes.Buffer
	conv := NewConverter(work, *packageName, *shadersDir)
	if err := conv.Run(&out); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		cleanup()
		os.Exit(1)
	}

	if err := os.WriteFile("shaders.go", out.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create shaders: %v\n", err)
		cleanup()
		os.Exit(1)
	}

	cmd := exec.Command("gofmt", "-s", "-w", "shaders.go")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "formatting shaders.go failed: %v\n", err)
		cleanup()
		os.Exit(1)
	}
}

type Converter struct {
	workDir    WorkDir
	shadersDir string

	packageName string

	glslvalidator *GLSLValidator
	spirv         *SPIRVCross
	fxc           *FXC
	msl           *MSL
}

func NewConverter(workDir WorkDir, packageName, shadersDir string) *Converter {
	if abs, err := filepath.Abs(shadersDir); err == nil {
		shadersDir = abs
	}

	conv := &Converter{}
	conv.workDir = workDir
	conv.shadersDir = shadersDir

	conv.packageName = packageName

	conv.glslvalidator = NewGLSLValidator()
	conv.spirv = NewSPIRVCross()
	conv.fxc = NewFXC()
	conv.msl = &MSL{
		WorkDir: workDir.Dir("msl"),
	}

	verifyBinaryPath(&conv.glslvalidator.Bin)
	verifyBinaryPath(&conv.spirv.Bin)
	// We cannot check fxc nor msl since they may depend on wine.

	conv.glslvalidator.WorkDir = workDir.Dir("glslvalidator")
	conv.fxc.WorkDir = workDir.Dir("fxc")
	conv.spirv.WorkDir = workDir.Dir("spirv")

	return conv
}

func verifyBinaryPath(bin *string) {
	new, err := exec.LookPath(*bin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to find %q: %v\n", *bin, err)
	} else {
		*bin = new
	}
}

func (conv *Converter) Run(out io.Writer) error {
	shaders, err := filepath.Glob(filepath.Join(conv.shadersDir, "*"))
	if err != nil {
		return fmt.Errorf("failed to list shaders in %q: %w", conv.shadersDir, err)
	}

	sort.Strings(shaders)

	var workers Workers

	type ShaderResult struct {
		Path    string
		Shaders []shaderSources
		Error   error
	}
	shaderResults := make([]ShaderResult, len(shaders))

	for i, shaderPath := range shaders {
		i, shaderPath := i, shaderPath

		switch filepath.Ext(shaderPath) {
		case ".vert", ".frag":
			workers.Go(func() {
				shaders, err := conv.Shader(shaderPath)
				shaderResults[i] = ShaderResult{
					Path:    shaderPath,
					Shaders: shaders,
					Error:   err,
				}
			})
		case ".comp":
			workers.Go(func() {
				shaders, err := conv.ComputeShader(shaderPath)
				shaderResults[i] = ShaderResult{
					Path:    shaderPath,
					Shaders: shaders,
					Error:   err,
				}
			})
		default:
			continue
		}
	}

	workers.Wait()

	var allErrors string
	for _, r := range shaderResults {
		if r.Error != nil {
			if len(allErrors) > 0 {
				allErrors += "\n\n"
			}
			allErrors += "--- " + r.Path + " --- \n\n" + r.Error.Error() + "\n"
		}
	}
	if len(allErrors) > 0 {
		return errors.New(allErrors)
	}

	fmt.Fprintf(out, "// Code generated by build.go. DO NOT EDIT.\n\n")
	fmt.Fprintf(out, "package %s\n\n", conv.packageName)
	fmt.Fprintf(out, "import (\n")
	fmt.Fprintf(out, "\t%q\n", "runtime")
	fmt.Fprintf(out, "\t_ %q\n", "embed")
	fmt.Fprintln(out)
	fmt.Fprintf(out, "\t%q\n", "ui/util/x/shader")
	fmt.Fprintf(out, ")\n\n")

	fmt.Fprintf(out, "var (\n")

	var genErr error
	for _, r := range shaderResults {
		if len(r.Shaders) == 0 {
			continue
		}

		name := filepath.Base(r.Path)
		name = strings.ReplaceAll(name, ".", "_")
		fmt.Fprintf(out, "\tShader_%s = ", name)

		multiVariant := len(r.Shaders) > 1
		if multiVariant {
			fmt.Fprintf(out, "[...]shader.Sources{\n")
		}

		writeGenerated := func(src []byte, prefix, path string, idx int) {
			if len(src) == 0 || genErr != nil {
				return
			}
			base := fmt.Sprintf("z%s.%d.%s", filepath.Base(path), idx, strings.ToLower(prefix))
			p := filepath.Join(filepath.Dir(path), base)
			genErr = os.WriteFile(p, src, 0o644)
		}
		for i, src := range r.Shaders {
			fmt.Fprintf(out, "shader.Sources{\n")
			fmt.Fprintf(out, "Name: %#v,\n", src.Name)
			if inp := src.Reflect.Inputs; len(inp) > 0 {
				fmt.Fprintf(out, "Inputs: %#v,\n", inp)
			}
			if u := src.Reflect.Uniforms; u.Size > 0 {
				fmt.Fprintf(out, "Uniforms: shader.UniformsReflection{\n")
				fmt.Fprintf(out, "Locations: %#v,\n", u.Locations)
				fmt.Fprintf(out, "Size: %d,\n", u.Size)
				fmt.Fprintf(out, "},\n")
			}
			if tex := src.Reflect.Textures; len(tex) > 0 {
				fmt.Fprintf(out, "Textures: %#v,\n", tex)
			}
			if imgs := src.Reflect.Images; len(imgs) > 0 {
				fmt.Fprintf(out, "Images: %#v,\n", imgs)
			}
			if bufs := src.Reflect.StorageBuffers; len(bufs) > 0 {
				fmt.Fprintf(out, "StorageBuffers: %#v,\n", bufs)
			}
			if wg := src.Reflect.WorkgroupSize; wg != [3]int{} {
				fmt.Fprintf(out, "WorkgroupSize: %#v,\n", wg)
			}
			writeGenerated(src.SPIRV, "SPIRV", r.Path, i)
			writeGenerated(src.GLSL100ES, "GLSL100ES", r.Path, i)
			writeGenerated(src.GLSL150, "GLSL150", r.Path, i)
			writeGenerated(src.DXBC, "DXBC", r.Path, i)
			writeGenerated(src.MetalLibs.MacOS, "MetalLibMacOS", r.Path, i)
			writeGenerated(src.MetalLibs.IOS, "MetalLibIOS", r.Path, i)
			writeGenerated(src.MetalLibs.IOSSimulator, "MetalLibIOSSimulator", r.Path, i)
			fmt.Fprintf(out, "}")
			if multiVariant {
				fmt.Fprintf(out, ",")
			}
			fmt.Fprintf(out, "\n")
		}
		if multiVariant {
			fmt.Fprintf(out, "}\n")
		}
		writeEmbedded := func(src []byte, prefix, path string, idx int) {
			base := fmt.Sprintf("z%s.%d.%s", filepath.Base(path), idx, strings.ToLower(prefix))
			if _, err := os.Stat(base); err != nil {
				return
			}
			field := strings.ReplaceAll(base, ".", "_")
			fmt.Fprintf(out, "//go:embed %s\n", base)
			fmt.Fprintf(out, "%s string\n", field)
		}
		for i, src := range r.Shaders {
			writeEmbedded(src.SPIRV, "SPIRV", r.Path, i)
			writeEmbedded(src.GLSL100ES, "GLSL100ES", r.Path, i)
			writeEmbedded(src.GLSL150, "GLSL150", r.Path, i)
			writeEmbedded(src.DXBC, "DXBC", r.Path, i)
			writeEmbedded(src.MetalLibs.MacOS, "MetalLibMacOS", r.Path, i)
			writeEmbedded(src.MetalLibs.IOS, "MetalLibIOS", r.Path, i)
			writeEmbedded(src.MetalLibs.IOSSimulator, "MetalLibIOSSimulator", r.Path, i)
		}
	}
	fmt.Fprintf(out, ")\n")
	writeInit := func(src []byte, prefix, field, path string, idx int, variants bool) {
		name := filepath.Base(path)
		name = strings.ReplaceAll(name, ".", "_")
		base := fmt.Sprintf("z%s.%d.%s", filepath.Base(path), idx, strings.ToLower(prefix))
		if _, err := os.Stat(base); err != nil {
			return
		}
		variable := strings.ReplaceAll(base, ".", "_")
		index := ""
		if variants {
			index = fmt.Sprintf("[%d]", idx)
		}
		fmt.Fprintf(out, "\t\tShader_%s%s.%s = %s\n", name, index, field, variable)
	}
	fmt.Fprintf(out, "func init() {\n")
	fmt.Fprintf(out, "\tconst (\n")
	fmt.Fprintf(out, "\t\topengles = %s\n", geeseExpr("linux", "freebsd", "openbsd", "windows", "js", "android", "darwin", "ios"))
	fmt.Fprintf(out, "\t\topengl = %s\n", geeseExpr("darwin"))
	fmt.Fprintf(out, "\t\td3d11 = %s\n", geeseExpr("windows"))
	fmt.Fprintf(out, "\t\tvulkan = %s\n", geeseExpr("linux", "android"))
	fmt.Fprintf(out, "\t)\n")
	for _, r := range shaderResults {
		variants := len(r.Shaders) > 1
		for i, src := range r.Shaders {
			fmt.Fprintf(out, "\tif vulkan {\n")
			writeInit(src.SPIRV, "SPIRV", "SPIRV", r.Path, i, variants)
			fmt.Fprintf(out, "\t}\n")
			fmt.Fprintf(out, "\tif opengles {\n")
			writeInit(src.GLSL100ES, "GLSL100ES", "GLSL100ES", r.Path, i, variants)
			fmt.Fprintf(out, "\t}\n")
			fmt.Fprintf(out, "\tif opengl {\n")
			writeInit(src.GLSL150, "GLSL150", "GLSL150", r.Path, i, variants)
			fmt.Fprintf(out, "\t}\n")
			fmt.Fprintf(out, "\tif d3d11 {\n")
			writeInit(src.DXBC, "DXBC", "DXBC", r.Path, i, variants)
			fmt.Fprintf(out, "\t}\n")
			fmt.Fprintf(out, "\tif runtime.GOOS == \"darwin\" {\n")
			writeInit(src.MetalLibs.MacOS, "MetalLibMacOS", "MetalLib", r.Path, i, variants)
			fmt.Fprintf(out, "\t}\n")
			fmt.Fprintf(out, "\tif runtime.GOOS == \"ios\" {\n")
			fmt.Fprintf(out, "if runtime.GOARCH == \"amd64\" {\n")
			writeInit(src.MetalLibs.IOSSimulator, "MetalLibIOSSimulator", "MetalLib", r.Path, i, variants)
			fmt.Fprintf(out, "\t\t} else {\n")
			writeInit(src.MetalLibs.IOS, "MetalLibIOS", "MetalLib", r.Path, i, variants)
			fmt.Fprintf(out, "\t\t}\n")
			fmt.Fprintf(out, "\t}\n")
		}
	}
	fmt.Fprintf(out, "}\n")

	return genErr
}

func geeseExpr(geese ...string) string {
	var checks []string
	for _, goos := range geese {
		checks = append(checks, fmt.Sprintf("runtime.GOOS == %q", goos))
	}
	return strings.Join(checks, " || ")
}

func (conv *Converter) Shader(shaderPath string) ([]shaderSources, error) {
	type Variant struct {
		FetchColorExpr string
		Header         string
	}
	variantArgs := [...]Variant{
		{
			FetchColorExpr: `_color.color`,
			Header:         `layout(push_constant) uniform Color { layout(offset=112) vec4 color; } _color;`,
		},
		{
			FetchColorExpr: `mix(_gradient.color1, _gradient.color2, clamp(vUV.x, 0.0, 1.0))`,
			Header:         `layout(push_constant) uniform Gradient { layout(offset=96) vec4 color1; vec4 color2; } _gradient;`,
		},
		{
			FetchColorExpr: `texture(tex, vUV)`,
			Header:         `layout(binding=0) uniform sampler2D tex;`,
		},
	}

	shaderTemplate, err := template.ParseFiles(shaderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %q: %w", shaderPath, err)
	}

	var variants []shaderSources
	for i, variantArg := range variantArgs {
		variantName := strconv.Itoa(i)
		var buf bytes.Buffer
		err := shaderTemplate.Execute(&buf, variantArg)
		if err != nil {
			return nil, fmt.Errorf("failed to execute template %q with %#v: %w", shaderPath, variantArg, err)
		}

		var sources shaderSources
		sources.Name = filepath.Base(shaderPath)

		src := buf.Bytes()
		sources.SPIRV, err = conv.glslvalidator.Convert(shaderPath, variantName, "vulkan", src)
		if err != nil {
			return nil, fmt.Errorf("failed to generate SPIR-V for %q: %w", shaderPath, err)
		}

		sources.SPIRV, err = spirvOpt(sources.SPIRV)
		if err != nil {
			return nil, fmt.Errorf("failed to optimize SPIR-V for %q: %w", shaderPath, err)
		}

		var reflect Metadata
		sources.GLSL100ES, reflect, err = conv.ShaderVariant(shaderPath, variantName, src, "es", "100")
		if err != nil {
			return nil, fmt.Errorf("failed to convert GLSL100ES:\n%w", err)
		}

		metal, _, err := conv.ShaderVariant(shaderPath, variantName, src, "msl", "10000")
		if err != nil {
			return nil, fmt.Errorf("failed to convert to Metal:\n%w", err)
		}

		metalIOS, _, err := conv.ShaderVariant(shaderPath, variantName, src, "mslios", "10000")
		if err != nil {
			return nil, fmt.Errorf("failed to convert to Metal:\n%w", err)
		}

		sources.MetalLibs, err = conv.msl.Compile(shaderPath, variantName, metal, metalIOS)
		if err != nil {
			if !errors.Is(err, exec.ErrNotFound) {
				return nil, fmt.Errorf("failed to build .metallib library:\n%w", err)
			}
		}

		hlsl, _, err := conv.ShaderVariant(shaderPath, variantName, src, "hlsl", "40")
		if err != nil {
			return nil, fmt.Errorf("failed to convert HLSL:\n%w", err)
		}

		sources.DXBC, err = conv.fxc.Compile(shaderPath, variantName, []byte(hlsl), "main", "4_0_level_9_1")
		if err != nil {
			// Attempt shader model 4.0. Only the gpu/headless
			// test shaders use features not supported by level
			// 9.1.
			sources.DXBC, err = conv.fxc.Compile(shaderPath, variantName, []byte(hlsl), "main", "4_0")
			if err != nil {
				if !errors.Is(err, exec.ErrNotFound) {
					return nil, fmt.Errorf("failed to compile HLSL: %w", err)
				}
			}
		}

		sources.GLSL150, _, err = conv.ShaderVariant(shaderPath, variantName, src, "glsl", "150")
		if err != nil {
			return nil, fmt.Errorf("failed to convert GLSL150:\n%w", err)
		}

		sources.Reflect = reflect

		variants = append(variants, sources)
	}

	// If the shader don't use the variant arguments, output only a single version.
	if bytes.Equal(variants[0].GLSL100ES, variants[1].GLSL100ES) {
		variants = variants[:1]
	}

	return variants, nil
}

func (conv *Converter) ShaderVariant(shaderPath, variant string, src []byte, lang, profile string) ([]byte, Metadata, error) {
	spirv, err := conv.glslvalidator.Convert(shaderPath, variant, lang, src)
	if err != nil {
		return nil, Metadata{}, fmt.Errorf("failed to generate SPIR-V for %q: %w", shaderPath, err)
	}

	dst, err := conv.spirv.Convert(shaderPath, variant, spirv, lang, profile)
	if err != nil {
		return nil, Metadata{}, fmt.Errorf("failed to convert shader %q: %w", shaderPath, err)
	}

	meta, err := conv.spirv.Metadata(shaderPath, variant, spirv)
	if err != nil {
		return nil, Metadata{}, fmt.Errorf("failed to extract metadata for shader %q: %w", shaderPath, err)
	}

	return dst, meta, nil
}

func (conv *Converter) ComputeShader(shaderPath string) ([]shaderSources, error) {
	sh, err := os.ReadFile(shaderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load shader %q: %w", shaderPath, err)
	}

	sources := shaderSources{
		Name: filepath.Base(shaderPath),
	}
	spirv, err := conv.glslvalidator.Convert(shaderPath, "", "glsl", sh)
	if err != nil {
		return nil, fmt.Errorf("failed to convert compute shader %q: %w", shaderPath, err)
	}

	sources.SPIRV, err = spirvOpt(spirv)
	if err != nil {
		return nil, fmt.Errorf("failed to optimize SPIR-V for %q: %w", shaderPath, err)
	}

	meta, err := conv.spirv.Metadata(shaderPath, "", spirv)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata for shader %q: %w", shaderPath, err)
	}
	sources.Reflect = meta

	metal, err := conv.spirv.Convert(shaderPath, "", spirv, "msl", "10000")
	if err != nil {
		return nil, fmt.Errorf("failed to convert GLSL130:\n%w", err)
	}

	metalIOS, err := conv.spirv.Convert(shaderPath, "", spirv, "mslios", "10000")
	if err != nil {
		return nil, fmt.Errorf("failed to convert GLSL130:\n%w", err)
	}

	sources.MetalLibs, err = conv.msl.Compile(shaderPath, "", metal, metalIOS)
	if err != nil {
		if !errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("failed to build .metallib library:\n%w", err)
		}
	}

	hlslSource, err := conv.spirv.Convert(shaderPath, "", spirv, "hlsl", "50")
	if err != nil {
		return nil, fmt.Errorf("failed to convert hlsl compute shader %q: %w", shaderPath, err)
	}

	sources.DXBC, err = conv.fxc.Compile(shaderPath, "0", []byte(hlslSource), "main", "5_0")
	if err != nil {
		if !errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("failed to compile hlsl compute shader %q: %w", shaderPath, err)
		}
	}

	return []shaderSources{sources}, nil
}

// Workers implements wait group with synchronous logging.
type Workers struct {
	running sync.WaitGroup
}

func (lg *Workers) Go(fn func()) {
	lg.running.Add(1)
	go func() {
		defer lg.running.Done()
		fn()
	}()
}

func (lg *Workers) Wait() {
	lg.running.Wait()
}

func unixLineEnding(s []byte) []byte {
	return bytes.ReplaceAll(s, []byte("\r\n"), []byte("\n"))
}
