package pipeline

import (
	"io"
)

var compilers = make([]Compiler, 0)

type Compiler interface {
	Match(Asset, string) bool
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
