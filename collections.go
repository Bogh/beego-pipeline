package pipeline

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"strings"
	"time"
)

const (
	AssetCss Asset = "css"
	AssetJs        = "js"
)

// Used for constants to define the asset type set
type Asset string

// A map of asset groups by name
type Collection map[string]*Group

// Init other group fields
func newCollection(c Collection) Collection {
	for i, group := range c {
		group.events = make(chan fsnotify.Event)
		c[i] = group
	}
	return c
}

type Collections map[Asset]Collection

func newCollections(css Collection, js Collection) Collections {
	return Collections{
		AssetCss: newCollection(css),
		AssetJs:  newCollection(js),
	}
}

// Keep configuration for an asset output
type Group struct {
	// Location inside the AppPath directory
	// specify this in case the root of static folder is not the default "/static"
	Root        string `json:",omitempty"`
	Sources     []string
	sourcePaths []string
	Output      string

	// Resulted file, default is the Output
	Result string `json:"-"`

	// Version string represents buy a hash generated from the file contents
	version string

	// Events channel
	events     chan fsnotify.Event `json:"-"`
	eventTimer *time.Timer
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
	if len(g.sourcePaths) == 0 {
		p := make([]string, 0)
		for _, pattern := range g.Sources {
			matches, err := filepath.Glob(g.AbsPath(pattern))
			if err != nil {
				return p, err
			}
			p = append(p, matches...)
		}
		g.sourcePaths = p
	}
	return g.sourcePaths, nil
}

// Normalized Output
func (g *Group) OutputPath() string {
	return g.AbsPath(g.VersionedPath())
}

func (g *Group) VersionedPath() string {
	path := g.Output
	if beego.BConfig.RunMode != "dev" && g.version != "" {
		dir := filepath.Dir(path)
		ext := filepath.Ext(path)
		base := strings.Replace(filepath.Base(path), ext, "", -1)

		path = filepath.Join(dir, fmt.Sprintf("%s.%s%s", base, g.version, ext))
	}
	return path
}

// Determine the Result path and return the value
func (g *Group) ResultPath() string {
	g.Result = g.RootedPath(g.Output)
	return g.Result
}

// sends the event to the events channel after waiting 500ms or cancel the time
// if a new arrives
// basically a debounce function, in case too many events are sent
func (g *Group) triggerWatch(event fsnotify.Event) {
	if g.eventTimer != nil {
		g.eventTimer.Stop()
	}

	g.eventTimer = time.AfterFunc(500*time.Millisecond, func() {
		g.events <- event
	})
}
