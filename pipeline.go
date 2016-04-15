package pipeline

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"os"
)

const (
	AssetCss Asset = iota
	AssetJs
)

var (
	compressors = make([]Compressor, 0)
	compilers   = make([]Compiler, 0)

	ErrNoCompiler = errors.New("No compiler found")
)

type Asset int

type Compiler interface {
	Match(Asset, string) bool
	Compile(string) (io.Reader, error)
}

// Define compressor interface
type Compressor interface {
	Match(Asset) bool
	// Should compress and concatenate the file in paths and save them in the
	// output
	Compress(io.Reader) (io.Reader, error)
}

func RegisterCompressor(c Compressor) {
	compressors = append(compressors, c)
}

func RegisterCompiler(c Compiler) {
	compilers = append(compilers, c)
}

type Processor struct {
	Asset      Asset
	Collection Collection
}

func NewProcessor(asset Asset, collection Collection) *Processor {
	return &Processor{asset, collection}
}

func (p *Processor) Process() error {
	// compile then compress
	for _, group := range p.Collection {
		compiled, err := p.Compile(group)
		if err != nil {
			beego.Error(err)
			continue
		}

		r, err := p.Compress(compiled)
		if err != nil {
			beego.Error(err)
			continue
		}

		go p.WriteGroup(group.OutputPath(), r)
	}

	return nil
}

func (p *Processor) WriteGroup(path string, r io.Reader) {
	oFile, err := os.Create(path)
	if err != nil {
		beego.Error(err)
		return
	}
	defer oFile.Close()

	_, err = io.Copy(oFile, io.Reader(r))
	if err != nil {
		beego.Error(err)
		return
	}
	beego.Debug("Generated: ", path)
}

func (p *Processor) Compile(group Group) (io.Reader, error) {
	// find compiler for each file in the group
	paths, err := group.SourcePaths()
	if err != nil {
		beego.Debug("Error found in matching group: ", err)
		return nil, err
	}

	readers := make([]io.Reader, 0, len(paths))

	for _, path := range paths {
		compiler := p.GetCompiler(path)
		if compiler == nil {
			return nil, fmt.Errorf("Compiler not found for asset: %s (%s)", path, p.GetAsset())
		}

		rc, err := compiler.Compile(path)
		if err != nil {
			beego.Debug("Error compiling ", path, ":", err)
			return nil, err
		}

		readers = append(readers, rc)
		beego.Debug("Found compiler", compiler, "for path", path)
	}

	return io.MultiReader(readers...), nil
}

func (p *Processor) GetCompiler(path string) Compiler {
	for _, compiler := range compilers {
		if compiler.Match(p.Asset, path) {
			return compiler
		}
	}
	return DefaultCompiler
}

// Accepts an io.Writer and returns an io.Reader
func (p *Processor) Compress(in io.Reader) (io.Reader, error) {
	compressor := p.GetCompressor()
	if compressor == nil {
		return nil, fmt.Errorf("Compressor not found for type: %s", p.GetAsset())
	}

	out, err := compressor.Compress(in)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	return out, nil
}

func (p *Processor) GetCompressor() Compressor {
	for _, compressor := range compressors {
		if compressor.Match(p.Asset) {
			return compressor
		}
	}
	return nil
}

func (p *Processor) GetAsset() string {
	switch p.Asset {
	case AssetCss:
		return "css"
	case AssetJs:
		return "js"
	default:
		return "unknown"
	}
}

// Run compilers and compressors in this pipeline
// TODO: make this concurrent using context
func Execute(config Config) error {
	for asset, collection := range config {
		processor := NewProcessor(asset, collection)
		processor.Process()
	}

	// p := NewProcessor("css", config.Css)
	// err := p.Process()
	// if err != nil {
	// 	return err
	// }

	return nil
}

// TODO: handle any errors that can be handled different
func appStartHook() error {
	config, err := loadConfig(nil)
	if err != nil {
		beego.Error(err)
		return err
	}

	// execute pipeline
	err = Execute(config)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	beego.AddAPPStartHook(appStartHook)
}
