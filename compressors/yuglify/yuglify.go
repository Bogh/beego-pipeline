package yuglify

import (
	"github.com/astaxie/beego"
	"github.com/bogh/beego-pipeline"
	"io"
	"os/exec"
)

// TODO: allow for customization of command
type YuglifyCompressor struct {
	t string
}

func (y YuglifyCompressor) Compress(done chan bool, r io.Reader) (io.ReadCloser, error) {
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
		done <- true
	}()

	return stdout, nil
}

func init() {
	pipeline.SetCompressor("css", YuglifyCompressor{"css"})
	pipeline.SetCompressor("js", YuglifyCompressor{"js"})
}
