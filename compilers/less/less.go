package less

import (
	"github.com/astaxie/beego"
	"github.com/bogh/beego-pipeline"
	"io"
	"strings"
)

type LessCompiler struct {
	*pipeline.Executor
}

func NewLessCompiler() *LessCompiler {
	path := beego.AppConfig.DefaultString("pipeline.command.less", "/usr/local/bin/lessc")
	args := beego.AppConfig.DefaultString("pipeline.command.less.arguments", "")
	return &LessCompiler{pipeline.NewExecutor(path, args)}
}

func (l *LessCompiler) Match(asset pipeline.Asset, filepath string) bool {
	return asset == pipeline.AssetCss &&
		strings.HasSuffix(filepath, ".less")
}

func (l *LessCompiler) Compile(filepath string) (io.Reader, error) {
	// start command and pipe the data through it
	return l.Executor.Pipe(l.BuildCmd(filepath), nil)
}

func (l *LessCompiler) String() string {
	return "Less Compiler"
}

func init() {
	pipeline.RegisterCompiler(NewLessCompiler())
}
