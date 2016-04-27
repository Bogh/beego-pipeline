package yuglify

import (
	"github.com/astaxie/beego"
	"github.com/bogh/beego-pipeline"
	"io"
)

// TODO: allow for customization of command
type YuglifyCompressor struct {
	*pipeline.Executor
	asset pipeline.Asset
}

func NewYuglifyCompressor(asset pipeline.Asset) *YuglifyCompressor {
	path := beego.AppConfig.DefaultString(
		"pipeline.command.yuglify",
		"/usr/local/bin/yuglify",
	)
	args := beego.AppConfig.DefaultString("pipeline.command.yuglify.arguments", "")
	return &YuglifyCompressor{pipeline.NewExecutor(path, args), asset}
}

func (y *YuglifyCompressor) Match(pipeline.Asset) bool {
	return true
}

func (y *YuglifyCompressor) Compress(r io.Reader) (io.Reader, error) {
	// start command and pipe the data through it
	stdout, err := y.Pipe(y.BuildCmd("--terminal", "--type", string(y.asset)), r)
	if err != nil {
		return nil, err
	}
	return stdout, nil
}

func init() {
	pipeline.RegisterCompressor(NewYuglifyCompressor(pipeline.AssetCss))
	pipeline.RegisterCompressor(NewYuglifyCompressor(pipeline.AssetJs))
}
