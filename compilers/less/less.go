package less

import (
	"github.com/bogh/beego-pipeline"
	"io"
	"os/exec"
	"strings"
)

type LessCompiler struct {
	*pipeline.Executor
}

func NewLessCompiler() *LessCompiler {
	return &LessCompiler{pipeline.NewExecutor()}
}

func (l *LessCompiler) Match(asset pipeline.Asset, filepath string) bool {
	return asset == pipeline.AssetCss &&
		strings.HasSuffix(filepath, ".less")
}

func (l *LessCompiler) Compile(filepath string) (io.Reader, error) {
	// start command and pipe the data through it
	cmd := exec.Command("/usr/local/bin/lessc", filepath)
	return l.Executor.Pipe(cmd, nil)
}

func (l *LessCompiler) String() string {
	return "Less Compiler"
}

func init() {
	pipeline.RegisterCompiler(&LessCompiler{})
}
