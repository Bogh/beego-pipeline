package pipeline

import (
	"github.com/astaxie/beego"
	"io"
)

var (
	compressors map[string]Compressor = make(map[string]Compressor)
)

// Define compressor interface
type Compressor interface {
	// Should compress and concatenate the file in paths and save them in the
	// output
	Compress(done chan bool, r io.Reader) (io.ReadCloser, error)
}

// There can be only one compressor per type of asset
// Receives type of compressor (css, js) and the Compressor interface
func SetCompressor(t string, c Compressor) {
	compressors[t] = c
}

// Run compilers and compressors in this pipeline
// TODO: make this concurrent using context
func Execute() error {
	// compress all in the output file
	beego.Debug("Processing CSS")
	p := NewProcessor("css", config.Css)
	err := p.Process()
	if err != nil {
		return err
	}

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
