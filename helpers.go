package pipeline

import (
	"github.com/astaxie/beego"
)

func asset(asset Asset, name string) (out string, err error) {
	return "", nil
}

func PipelineCss(name string) (string, error) {
	return asset(AssetCss, name)
}

func PipelineJs(name string) (string, error) {
	return asset(AssetJs, name)
}

func registerHelpers(config *Config) {
	beego.AddFuncMap("pipeline_css", PipelineCss)
	beego.AddFuncMap("pipeline_js", PipelineJs)
}
