package pipeline

import (
	"io"
)

var compressors = make([]Compressor, 0)

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

func (p *Processor) GetCompressor() Compressor {
	for _, compressor := range compressors {
		if compressor.Match(p.Asset) {
			return compressor
		}
	}
	return nil
}
