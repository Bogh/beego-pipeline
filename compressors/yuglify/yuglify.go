package yuglify

import (
	"github.com/astaxie/beego"
	"github.com/bogh/beego-pipeline"
	// "os/exec"
)

// TODO: allow for customization of command
type YuglifyCompressor struct{}

func (y YuglifyCompressor) Compress(paths pipeline.Paths, output string) error {
	beego.Debug("Compressing ", paths)
	// concatenate paths
	return nil
}

func init() {
	pipeline.SetCompressor("css", YuglifyCompressor{})
	pipeline.SetCompressor("js", YuglifyCompressor{})
}
