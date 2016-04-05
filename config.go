package pipeline

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"io/ioutil"
	"path/filepath"
)

var (
	fp     FilePath
	config *Config
)

type Config struct {
	Css Outputs
	Js  Outputs
}

type ConfigPather interface {
	// returns file path to the conf file
	Path() (string, error)
}

type FilePath string

func (fp FilePath) Path() (string, error) {
	fn := filepath.Join(beego.AppPath, "conf", "pipeline.json")
	if !utils.FileExists(fn) {
		beego.Debug("pipeline.json not found.")
		return "", errors.New("File does not exist")
	}
	return fn, nil
}

// find conf/pipeline.conf and load it
func loadConfig(cp ConfigPather) error {
	if cp == nil {
		cp = fp
	}
	path, err := cp.Path()
	if err != nil {
		return err
	}
	beego.Debug("Found pipeline config file: ", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	beego.Debug("Loaded pipeline data", *config)
	return nil
}
