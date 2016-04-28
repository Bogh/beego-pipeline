package pipeline

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"os"
)

type Processor struct {
	Asset      Asset
	Collection Collection
}

func NewProcessor(asset Asset, collection Collection) *Processor {
	return &Processor{
		asset,
		collection,
	}
}

// Watch files in this group for changes and recompile
func (p *Processor) Watch() error {
	beego.Debug("Start watching files for changes.")
	// start watching groups
	for _, group := range p.Collection {
		go func(g *Group) {
			for {
				select {
				case <-g.events:
					p.ProcessGroup(g)
				}
			}
		}(group)
	}
	return nil
}

func (p *Processor) Process() error {
	// Start watching
	p.Watch()

	// compile then compress
	for _, group := range p.Collection {
		if err := p.ProcessGroup(group); err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) ProcessGroup(group *Group) error {
	compiled, err := p.Compile(group, isDev())
	if err != nil {
		beego.Error(err)
		return err
	}

	if isDev() {
		return nil
	}

	versionHash := md5.New()

	r, err := p.Compress(io.TeeReader(compiled, versionHash))
	if err != nil {
		beego.Error(err)
		return err
	}

	group.version = string(fmt.Sprintf("%x", versionHash.Sum(nil)[:8]))

	p.Write(group.OutputPath(), r)
	return nil
}

// Write data from reader to the group file
func (p *Processor) Write(path string, r io.Reader) error {
	oFile, err := os.Create(path)
	if err != nil {
		beego.Error(err)
		return err
	}
	defer oFile.Close()

	_, err = io.Copy(oFile, r)
	if err != nil {
		beego.Error(err)
		return err
	}

	beego.Debug("Generated: ", path)
	return nil
}

// If `write` is true then save each file to it's own path
func (p *Processor) Compile(group *Group, write bool) (io.Reader, error) {
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

		if write && !compiler.RequireCompile() {
			continue
		}

		rc, err := compiler.Compile(path)
		if err != nil {
			beego.Debug("Error compiling ", path, ":", err)
			return nil, err
		}

		if write {
			// write to the file
			nPath := group.NormAsset(path, p.Asset)
			if path == nPath {
				return nil, &ErrOverridingPath{path}
			}

			p.Write(nPath, rc)
		}

		readers = append(readers, rc)
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
