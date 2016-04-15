package pipeline

import (
	"github.com/astaxie/beego"
	"io"
	"os"
)

var (
	DefaultCompiler = new(NopCompiler)
)

type NopCompiler struct{}

func (r *NopCompiler) Match(asset Asset, filepath string) bool {
	return true
}

func (r *NopCompiler) Compile(filepath string) (io.Reader, error) {
	// read the file and return an io
	beego.Debug("Nop compiler: ", filepath)
	f, err := os.Open(filepath)
	return &AutoCloseReader{f}, err
}

func (n *NopCompiler) String() string {
	return "Nop Compiler"
}
