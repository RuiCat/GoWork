// Code generated by gioui.org/cpu/cmd/compile DO NOT EDIT.

//go:build !(linux && (arm64 || arm || amd64))
// +build !linux !arm64,!arm,!amd64

package piet

import "ui/util/x/cpu"

var Kernel4ProgramInfo *cpu.ProgramInfo

type Kernel4DescriptorSetLayout struct{}

const Kernel4Hash = ""

func (l *Kernel4DescriptorSetLayout) Binding0() *cpu.BufferDescriptor {
	panic("unsupported")
}

func (l *Kernel4DescriptorSetLayout) Binding1() *cpu.BufferDescriptor {
	panic("unsupported")
}

func (l *Kernel4DescriptorSetLayout) Binding2() *cpu.ImageDescriptor {
	panic("unsupported")
}

func (l *Kernel4DescriptorSetLayout) Binding3() *cpu.ImageDescriptor {
	panic("unsupported")
}
