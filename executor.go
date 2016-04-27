package pipeline

import (
	"bytes"
	"github.com/astaxie/beego"
	"io"
	"os/exec"
)

type Executor struct {
	CmdPath string
	CmdArgs string
}

func NewExecutor(path, args string) *Executor {
	return &Executor{path, args}
}

func (e *Executor) BuildCmd(extraArgs ...string) *exec.Cmd {
	args := append(append([]string{}, e.CmdArgs), extraArgs...)
	return exec.Command(e.CmdPath, args...)
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
