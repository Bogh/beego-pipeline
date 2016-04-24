package pipeline

import (
	"fmt"
	"github.com/astaxie/beego"
	"golang.org/x/net/context"
	"io"
	"os"
)

type Processor struct {
	Asset      Asset
	Collection Collection

	ctx context.Context
}

func NewProcessor(asset Asset, collection Collection) *Processor {
	return &Processor{
		asset,
		collection,
		context.Background(), // parent context
	}
}

// Watch files in this group for changes and recompile
func (p *Processor) Watch() error {
	return nil
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

		p.WriteGroup(group.OutputPath(), r)
	}

	return nil
}

// Write data from reader to the group file
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
