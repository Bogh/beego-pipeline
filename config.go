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
	Config ConfigContainer
)

type ConfigContainer map[string]Output

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

func loadConfig() error {
	AppDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	confFile := filepath.Join(AppDir, "conf", "pipeline.json")
	beego.Debug("Found pipeline config file: ", confFile)

	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		return err
	}
	beego.Debug("Loaded pipeline data", Config)
	return nil
}

// find conf/pipeline.conf and load it
// TODO: handle any errors that can be handled different
func appStartHook() error {
	err := loadConfig()
	if err != nil {
		beego.Error(err)
	}

	return nil
}

func init() {
	beego.AddAPPStartHook(appStartHook)
}
