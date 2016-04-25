package pipeline

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"path/filepath"
)

const (
	AssetCss Asset = "css"
	AssetJs        = "js"
)

var (
	config  *Config
	watcher *fsnotify.Watcher
)

// Used for constants to define the asset type set
type Asset string

// Hold a map of asset types each containing a collection
type Config struct {
	Collections map[Asset]Collection
	Watcher     *fsnotify.Watcher
}

// struct to hold different types of assets

func (c *Config) GetAssetGroup(asset Asset, name string) (*Group, error) {
	collection, ok := c.Collections[asset]
	if !ok {
		return nil, ErrAssetNotFound
	}

	group, ok := collection[name]
	if !ok {
		return nil, ErrAssetNotFound
	}
	return &group, nil
}

func (c *Config) GetAssetTpl(asset Asset) string {
	return map[Asset]string{
		AssetCss: `<link href="%s">`,
		AssetJs:  `<script src="%s"></script>`,
	}[asset]
}

func getConfigPath() (string, error) {
	fn := filepath.Join(beego.AppPath, "conf", "pipeline.json")
	if !utils.FileExists(fn) {
		beego.Debug("pipeline.json not found.")
		return "", errors.New("File does not exist")
	}
	return fn, nil
}

// A map of asset groups by name
type Collection map[string]Group

// Keep configuration for an asset output
type Group struct {
	// Location inside the AppPath directory
	// specify this in case the root of static folder is not the default "/static"
	Root    string `json:",omitempty"`
	Sources []string
	Output  string

	// Resulted file, default is the Output
	Result string `json:"-"`
}

// Return absolute path for provided path, prepending AppPath and Root
func (g *Group) AbsPath(path string) string {
	return filepath.Join(beego.AppPath, g.RootedPath(path))
}

func (g *Group) RootedPath(paths ...string) string {
	if g.Root == "" {
		g.Root = "/static"
	}

	return filepath.Join(append([]string{g.Root}, paths...)...)
}

func (g *Group) SourcePaths() ([]string, error) {
	p := []string{}
	for _, pattern := range g.Sources {
		matches, err := filepath.Glob(g.AbsPath(pattern))
		if err != nil {
			return p, err
		}
		p = append(p, matches...)
	}
	return p, nil
}

// Normalized Output
func (g *Group) OutputPath() string {
	return g.AbsPath(g.Output)
}

// Determine the Result path and return the value
// TODO: This method will calculate the version hash
func (g *Group) ResultPath() string {
	g.Result = g.RootedPath(g.Output)
	return g.Result
}

// find conf/pipeline.conf and load it
func loadConfig() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	beego.Debug("Found pipeline config file: ", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := struct {
		Css Collection
		Js  Collection
	}{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	beego.Debug("Loaded pipeline data", c)
	config = &Config{
		Collections: map[Asset]Collection{
			AssetCss: c.Css,
			AssetJs:  c.Js,
		},
		Watcher: watcher,
	}

	return config, nil
}
