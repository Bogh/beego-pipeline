package pipeline

import (
	"github.com/astaxie/beego"
)

// Run compilers and compressors in this pipeline
// TODO: make this concurrent using context
func Execute() error {
	for asset, collection := range *config {
		processor := NewProcessor(asset, collection)
		processor.Process()
	}

	return nil
}

// TODO: handle any errors that can be handled different
func registerPipeline() error {
	_, err := loadConfig()
	if err != nil {
		beego.Error(err)
		return err
	}

	// execute pipeline
	go Execute()
	go registerHelpers()
	return nil
}

func init() {
	beego.AddAPPStartHook(registerPipeline)
}
