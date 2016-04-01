package pipeline

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	AppDir string
)

type Config struct {
	Css map[string]Output
	Js  map[string]Output
}

type Output struct {
	// Location inside the AppDirectory
	// specify this in case the root of static folder is not the default "/static"
	Root    string `json:",omitempty"`
	Sources []string
	Output  string
}

// Return absolute path for provided path, prepending AppDir and Root
func (o Output) Path(path string) string {
	return filepath.Join(AppDir, o.Root, path)
}

// find conf/pipeline.conf and load it
func loadConfig() (*Config, error) {
	AppDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	f := filepath.Join(AppDir, "conf", "pipeline.json")
	beego.Debug("Found pipeline config file: ", f)

	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	beego.Debug("Loaded pipeline data", *config)
	return config, nil
}
