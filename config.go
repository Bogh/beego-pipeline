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

var config *Config

// Hold a map of asset types each containing a collection
type Config struct {
	Collections
	Watcher *fsnotify.Watcher
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
	return group, nil
}

func (c *Config) GetAssetTpl(asset Asset) string {
	return map[Asset]string{
		AssetCss: `<link href="%s">`,
		AssetJs:  `<script src="%s"></script>`,
	}[asset]
}

// Add all files to watcher
func (c *Config) watches() error {
	if !isDev() {
		return nil
	}

	for _, collection := range c.Collections {
		for _, group := range collection {
			paths, err := group.SourcePaths()
			if err != nil {
				beego.Error("Error adding watches.", err)
				return err
			}

			for _, path := range paths {
				c.Watcher.Add(path)
			}
		}
	}

	go c.listen()
	return nil
}

// wait for file changes
func (c *Config) listen() {
	for {
		select {
		case e := <-c.Watcher.Events:
			// identify the groups that have been changed and forward the event
			c.forward(e)
		case err := <-c.Watcher.Errors:
			beego.Error("Watcher error:", err)
		}
	}
}

func (c *Config) forward(e fsnotify.Event) {
	for _, collection := range c.Collections {
		for _, group := range collection {
			paths, err := group.SourcePaths()
			if err != nil {
				beego.Error("Error getting group paths:", err)
				return
			}

			for _, p := range paths {
				if p == e.Name {
					// send event to the matched group
					group.triggerWatch(e)
					continue
				}
			}
		}
	}
}

func getConfigPath() (string, error) {
	fn := filepath.Join(beego.AppPath, "conf", "pipeline.json")
	if !utils.FileExists(fn) {
		beego.Debug("pipeline.json not found.")
		return "", errors.New("File does not exist")
	}
	return fn, nil
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
		Collections: newCollections(c.Css, c.Js),
		Watcher:     watcher,
	}

	config.watches()

	return config, nil
}
