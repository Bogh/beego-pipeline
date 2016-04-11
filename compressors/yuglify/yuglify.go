package yuglify

import (
	"github.com/astaxie/beego"
	"github.com/bogh/beego-pipeline"
	"io"
	"os/exec"
)

// TODO: allow for customization of command
type YuglifyCompressor struct {
	asset pipeline.Asset
}

func (y *YuglifyCompressor) Compress(r io.Reader) (*pipeline.AutoCloseReader, error) {
	// start command and pipe the data through it
	cmdArgs := []string{
		"--terminal",
		"--type", "css",
	}
	cmd := exec.Command("/usr/local/bin/yuglify", cmdArgs...)
	// write data to command
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	// read from stdout and write to file
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	defer stdin.Close()
	n, _ := io.Copy(stdin, r)
	beego.Debug("Bytes sent to stdin:", n)

	go func() {
		err = cmd.Run()
		if err != nil {
			beego.Error(err)
		}
	}()

	return &pipeline.AutoCloseReader{stdout}, nil
}

func (y *YuglifyCompressor) Match(pipeline.Asset) bool {
	return true
}

func init() {
	pipeline.SetCompressor(&YuglifyCompressor{pipeline.AssetCss})
	pipeline.SetCompressor(&YuglifyCompressor{pipeline.AssetJs})
}
