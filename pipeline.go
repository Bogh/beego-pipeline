package pipeline

import (
	"github.com/astaxie/beego"
)

var (
	DefaultPipeline *Pipeline
)

func NewPipeline(config *Config) (*Pipeline, error) {
	return &Pipeline{
		Config:      config,
		compressors: make(map[string]Compresser, 2),
	}, nil
}

type Pipeline struct {
	Config *Config

	compressors map[string]Compresser
}

// There can be only one compressor per type of asset
// Receives type of compressor (css, js) and the Compresser interface
func (p Pipeline) SetCompressor(t string, c Compresser) {
	p.compressors[t] = c
}

func (p Pipeline) Paths(o Output) Paths {
	return nil
}

// Run compilers and compressors in this pipeline
func (p Pipeline) Execute() {
	// compress all in the output file
	beego.Debug("Executing pipeline")

	beego.Debug("Processing CSS")
	p.Process(p.Config.Css)
}

// Process a set of outputs
func (p Pipeline) Process(os Outputs) {
	// get list of files and
	// normalize paths and check if they exist, otherwise issue an ignore
	for _, output := range os {
		// normalized paths
		np := p.Paths(output)
		beego.Debug("Found paths: ", np, " for output ", output.Output)
	}
}

// Define compressor interface
type Compresser interface {
	// Should compress and concatenate the file in paths and save them in the
	// output
	Compress(paths Paths, output string) error
}

// func (p *Pipeline) RegisterCompiler() error {
// 	return nil
// }

func SetCompressor(t string, c Compresser) {
	DefaultPipeline.SetCompressor(t, c)
}

func Execute() error {
	return DefaultPipeline.Execute()
}

// TODO: handle any errors that can be handled different
func appStartHook() error {
	config, err := loadConfig(nil)
	if err != nil {
		beego.Error(err)
	}

	DefaultPipeline, err = NewPipeline(config)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	beego.AddAPPStartHook(appStartHook)
}
