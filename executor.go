package pipeline

import (
	"github.com/astaxie/beego"
	"io"
	"os/exec"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Pipe(cmd *exec.Cmd, r io.Reader) (out io.ReadCloser, err error) {
	// read from stdout and write to file
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if r != nil {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}

		go func(command string) {
			defer stdin.Close()
			_, err = io.Copy(stdin, r)
			if err != nil {
				beego.Error(command, err)
			}
		}(cmd.Path)
	}

	go func() {
		err = cmd.Run()
		if err != nil {
			beego.Error(err)
		}
	}()
	return stdout, nil
}
