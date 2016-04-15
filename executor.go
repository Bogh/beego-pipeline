package pipeline

import (
	"bytes"
	"github.com/astaxie/beego"
	"io"
	"os/exec"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Pipe(cmd *exec.Cmd, r io.Reader) (io.Reader, error) {
	if r != nil {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}

		go func(command string) {
			defer stdin.Close()
			var buf bytes.Buffer
			tr := io.TeeReader(r, &buf)
			_, err = io.Copy(stdin, tr)
			if err != nil {
				beego.Error(command, err)
				beego.Debug("Read data:", buf.String())
			}
		}(cmd.Path)
	}

	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}