package less

import (
	"github.com/astaxie/beego"
	"github.com/bogh/beego-pipeline"
	"os/exec"
	"strings"
)

type LessCompiler struct {
}

func (l *LessCompiler) Match(asset pipeline.Asset, filepath string) bool {
	return asset == pipeline.AssetCss &&
		strings.HasSuffix(filepath, ".less")
}

func (l *LessCompiler) Compile(filepath string) (*pipeline.AutoCloseReader, error) {
	// start command and pipe the data through it
	cmd := exec.Command("/usr/local/bin/lessc", filepath)

	// // write data to command
	// stdin, err := cmd.StdinPipe()
	// if err != nil {
	// 	return nil, err
	// }

	// read from stdout and write to file
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	// defer stdin.Close()
	// n, _ := io.Copy(stdin, r)
	// beego.Debug("Bytes sent to lessc stdin:", n)

	go func() {
		err = cmd.Run()
		if err != nil {
			beego.Error(err)
		}
	}()

	return &pipeline.AutoCloseReader{stdout}, nil
}

func (l *LessCompiler) String() string {
	return "Less Compiler"
}

func init() {
	pipeline.RegisterCompiler(&LessCompiler{})
}
