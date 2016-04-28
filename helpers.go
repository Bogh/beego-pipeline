package pipeline

import (
	"fmt"
	"github.com/astaxie/beego"
	"html/template"
)

func asset(asset Asset, name string) (template.HTML, error) {
	var html string

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

	tpl := config.GetAssetTpl(asset)

	paths, err := group.ResultPaths(asset)
	if err != nil {
		return template.HTML(""), err
	}
	for _, path := range paths {
		html += fmt.Sprintf(tpl, path) + "\n"
	}
	return template.HTML(html), nil
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
