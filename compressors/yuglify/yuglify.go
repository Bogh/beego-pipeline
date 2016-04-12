package yuglify

import (
	"github.com/bogh/beego-pipeline"
	"io"
	"os/exec"
)

// TODO: allow for customization of command
type YuglifyCompressor struct {
	*pipeline.Executor
	asset pipeline.Asset
}

func NewYuglifyCompressor(asset pipeline.Asset) *YuglifyCompressor {
	return &YuglifyCompressor{pipeline.NewExecutor(), asset}
}

func (y *YuglifyCompressor) Match(pipeline.Asset) bool {
	return true
}

func (y *YuglifyCompressor) Compress(r io.Reader) (*pipeline.AutoCloseReader, error) {
	// start command and pipe the data through it
	cmdArgs := []string{
		"--terminal",
		"--type", "css",
	}
	stdout, err := y.Executor.Pipe(
		exec.Command("/usr/local/bin/yuglify", cmdArgs...),
		r,
	)
	if err != nil {
		return nil, err
	}
	return &pipeline.AutoCloseReader{stdout}, nil
}

func init() {
	pipeline.RegisterCompressor(NewYuglifyCompressor(pipeline.AssetCss))
	pipeline.RegisterCompressor(NewYuglifyCompressor(pipeline.AssetJs))
}
