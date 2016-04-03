package yuglify

import (
	"github.com/astaxie/beego"
	"github.com/bogh/beego-pipeline"
)

type YuglifyCompressor struct{}

func (y YuglifyCompressor) Compress(paths pipeline.Paths, output string) error {
	return nil
}

func init() {
	pipeline.SetCompressor("css", YuglifyCompressor{})
	pipeline.SetCompressor("js", YuglifyCompressor{})
}
