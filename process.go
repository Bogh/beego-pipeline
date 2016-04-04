package pipeline

import (
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"os"
	"os/exec"
)

type Processor struct {
	// type of asset
	t       string
	outputs Outputs
}

func NewProcessor(t string, outputs Outputs) *Processor {
	return &Processor{t, outputs}
}

func (p *Processor) Process() error {
	// compile then compress

	// get list of files for each output and
	// normalize paths and check if they exist, otherwise issue an ignore
	for _, output := range p.outputs {
		err := p.Compress(output)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

// Accepts an io.Writer and returns an io.Reader
func (p *Processor) Compress(output Output) error {
	// normalized paths
	sources, _ := output.Paths()
	beego.Debug("Found paths: ", sources, " for output ", output.Output)

	files := make([]io.Reader, len(sources))
	for i, path := range sources {
		f, _ := os.Open(path)
		files[i] = io.Reader(f)
		defer f.Close()
	}
	r := io.MultiReader(files...)

	// output File
	oFile, _ := os.OpenFile(output.NOutput(), os.O_WRONLY|os.O_CREATE, 0644)
	defer oFile.Close()

	// start command and pipe the data through it
	cmd := exec.Command("yuglify", "--terminal", "--type css")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer stdin.Close()

	// write data to command
	go func() {
		io.Copy(stdin, r)
	}()

	// read from stdout and write to file
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	defer stdout.Close()

	go func() {
		io.Copy(os.Stdout, stdout)
	}()

	// stderr, err := cmd.StderrPipe()
	// if err != nil {
	// 	panic(err)
	// }

	// go func() {
	//        ioutil.ReadAll(cmd.Stderr)
	// }()

	err = cmd.Run()
	if err != nil {
		beego.Error(fmt.Sprintf("%T", err))
		panic(err)
	}
	return nil
}
