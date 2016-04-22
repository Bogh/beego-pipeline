package pipeline

import (
	"fmt"
	"github.com/astaxie/beego"
	"html/template"
)

const (
	linkTpl = `<link href="%s">`
)

func asset(asset Asset, name string) (out template.HTML, err error) {
	// Find the asset
	group, err := config.GetAssetGroup(asset, name)
	if err != nil {
		// If asset not found log a warning and return empty string
		if err == ErrAssetNotFound {
			beego.Warn("%s: Asset group not found", name)
			return template.HTML(""), nil

		}
		return template.HTML(""), err
	}

	return template.HTML(fmt.Sprintf(linkTpl, group.ResultPath())), nil
}

func PipelineCss(name string) (template.HTML, error) {
	return asset(AssetCss, name)
}

func PipelineJs(name string) (template.HTML, error) {
	return asset(AssetJs, name)
}

func registerHelpers() {
	beego.AddFuncMap("pipeline_css", PipelineCss)
	beego.AddFuncMap("pipeline_js", PipelineJs)
}
