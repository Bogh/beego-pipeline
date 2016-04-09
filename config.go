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
	fp FilePath
)

type Config map[Asset]Outputs

type config struct {
	Css Outputs
	Js  Outputs
}

type ConfigPath interface {
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
func loadConfig(cp ConfigPath) (Config, error) {
	if cp == nil {
		cp = fp
	}
	path, err := cp.Path()
	if err != nil {
		return nil, err
	}
	beego.Debug("Found pipeline config file: ", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	c = config{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	beego.Debug("Loaded pipeline data", *config)
	return nil, Config{
		AssetCss: c.Css,
		AssetJs: c.Js
	}
}
