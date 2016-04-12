package pipeline

import (
	"github.com/astaxie/beego"
	"os"
)

var (
	DefaultCompiler = new(NopCompiler)
)

type NopCompiler struct{}

func (r *NopCompiler) Match(asset Asset, filepath string) bool {
	return true
}

func (r *NopCompiler) Compile(filepath string) (*AutoCloseReader, error) {
	// read the file and return an io
	beego.Debug("Nop compiler: ", filepath)
	f, err := os.Open(filepath)
	return &AutoCloseReader{f}, err
}

func (n *NopCompiler) String() string {
	return "Nop Compiler"
}
