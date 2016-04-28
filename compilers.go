package pipeline

import (
	"io"
	"os"
)

var (
	compilers       = make([]Compiler, 0)
	DefaultCompiler = new(NopCompiler)
)

type Compiler interface {
	Match(Asset, string) bool

	// Return true if the asset must be compiled in order to work
	RequireCompile() bool

	Compile(string) (io.Reader, error)
}

func RegisterCompiler(c Compiler) {
	compilers = append(compilers, c)
}

func (p *Processor) GetCompiler(path string) Compiler {
	for _, compiler := range compilers {
		if compiler.Match(p.Asset, path) {
			return compiler
		}
	}
	return DefaultCompiler
}

type NopCompiler struct{}

func (n *NopCompiler) RequireCompile() bool {
	return false
}

func (n *NopCompiler) Match(asset Asset, filepath string) bool {
	return true
}

func (n *NopCompiler) Compile(filepath string) (io.Reader, error) {
	// read the file and return an io
	f, err := os.Open(filepath)
	return &AutoCloseReader{f}, err
}

func (n *NopCompiler) String() string {
	return "Nop Compiler"
}
