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

type Config map[Asset]Collection

type config struct {
	Css Collection
	Js  Collection
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

type Collection map[string]Group

type Group struct {
	// Location inside the AppPath directory
	// specify this in case the root of static folder is not the default "/static"
	Root    string `json:",omitempty"`
	Sources []string
	Output  string
}

// Return absolute path for provided path, prepending AppPath and Root
func (o *Group) Path(path string) string {
	root := o.Root
	if root == "" {
		root = "/static"
	}
	return filepath.Join(beego.AppPath, root, path)
}

func (o *Group) SourcePaths() ([]string, error) {
	p := []string{}
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
func (o *Group) OutputPath() string {
	return o.Path(o.Output)
}

// find conf/pipeline.conf and load it
func loadConfig(cp ConfigPath) (*Config, error) {
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
		return nil, err
	}

	c := config{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	beego.Debug("Loaded pipeline data", c)
	return &Config{
		AssetCss: c.Css,
		AssetJs:  c.Js,
	}, nil
}
