package pipeline

import (
	"github.com/astaxie/beego"
)

var (
	DefaultPipeline *Pipeline
)

func NewPipeline(config *Config) (*Pipeline, error) {
	return &Pipeline{config}, nil
}

type Pipeline struct {
	Config *Config
}

// TODO: handle any errors that can be handled different
func appStartHook() error {
	config, err := loadConfig(nil)
	if err != nil {
		beego.Error(err)
	}

	DefaultPipeline, err = NewPipeline(config)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	beego.AddAPPStartHook(appStartHook)
}
