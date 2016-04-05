package pipeline

import (
	"github.com/astaxie/beego"
	"io"
	"path/filepath"
)

const (
	AssetCss Asset = iota
	AssetJs
)

var (
	compressors = make([]Compressor, 1)
	compilers   = make([]Compiler, 1)
)

type Asset int

type Typer interface {
	// Return the type of asset that can be handled
	Type() Asset
}

type Matcher interface {
	// Return true if it can handle the file
	Match(string) bool
}

type Compiler interface {
	Typer
	Matcher
	Compile() (chan bool, chan error)
}

// Define compressor interface
type Compressor interface {
	Typer
	Matcher
	// Should compress and concatenate the file in paths and save them in the
	// output
	Compress(done chan bool, r io.Reader) (io.ReadCloser, error)
}

type Outputs map[string]Output

type Paths []string

type Output struct {
	// Location inside the AppPath directory
	// specify this in case the root of static folder is not the default "/static"
	Root    string `json:",omitempty"`
	Sources Paths
	Output  string
}

// Return absolute path for provided path, prepending AppPath and Root
func (o *Output) Path(path string) string {
	root := o.Root
	if root == "" {
		root = "/static"
	}
	return filepath.Join(beego.AppPath, root, path)
}

func (o *Output) Paths() (Paths, error) {
	p := Paths{}
	for _, pattern := range o.Sources {
		matches, err := filepath.Glob(o.Path(pattern))
		if err != nil {
			return p, err
		}
		p = append(p, matches...)
	}
	return p, nil
}

// Normalized Output
func (o *Output) NOutput() string {
	return o.Path(o.Output)
}

// There can be only one compressor per type of asset
// Receives type of compressor (css, js) and the Compressor interface
func SetCompressor(c Compressor) {
	compressors = append(compressors, c)
}

func AddCompiler(c Compiler) {
	compilers = append(compilers, c)
}

// Run compilers and compressors in this pipeline
// TODO: make this concurrent using context
func Execute() error {
	// compress all in the output file
	beego.Debug("Processing CSS")
	// p := NewProcessor("css", config.Css)
	// err := p.Process()
	// if err != nil {
	// 	return err
	// }

	return nil
}

// TODO: handle any errors that can be handled different
func appStartHook() error {
	err := loadConfig(nil)
	if err != nil {
		beego.Error(err)
		return err
	}

	// execute pipeline
	err = Execute()
	if err != nil {
		return err
	}

	// pipeline, err = NewPipeline(config)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func init() {
	beego.AddAPPStartHook(appStartHook)
}
