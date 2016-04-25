package pipeline

import (
	"github.com/astaxie/beego"
)

// TODO: handle any errors that can be handled different
func registerPipeline() error {
	_, err := loadConfig()
	if err != nil {
		beego.Error(err)
		return err
	}

	// execute pipeline
	for asset, collection := range config.Collections {
		processor := NewProcessor(asset, collection)
		processor.Process()
	}

	return nil
	registerHelpers()
	return nil
}

func init() {
	beego.AddAPPStartHook(registerPipeline)
}
