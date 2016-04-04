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

type Outputs map[string]Output

type Paths []string

type Output struct {
	// Location inside the AppPath directory
	// specify this in case the root of static folder is not the default "/static"
	Root    string `json:",omitempty"`
	Sources Paths
	Output  string
}

// Return absolute path for provided path, prepending AppPath and Root
func (o Output) Path(path string) string {
	root := o.Root
	if root == "" {
		root = "/static"
	}
	return filepath.Join(beego.AppPath, root, path)
}

func (o Output) Paths() (Paths, error) {
	p := Paths{}
	for _, pattern := range o.Sources {
		matches, err := filepath.Glob(o.Path(pattern))
		if err != nil {
			return p, err
		}
		p = append(p, matches...)
	}
	return p, nil
}

// Normalized Output
func (o Output) NOutput() string {
	return o.Path(o.Output)
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
