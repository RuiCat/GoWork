package gen

import (
	"ui/util/gocc/internal/util/gen/golang"
)

func Gen(outDir string) {
	golang.GenRune(outDir)
	golang.GenLitConv(outDir)
}
