package pipeline

import (
	"github.com/astaxie/beego"
	"github.com/fsnotify/fsnotify"
	"path/filepath"
)

// A map of asset groups by name
type Collection map[string]Group

// Init other group fields
func NewCollection(c Collection) Collection {
	for i, group := range c {
		group.events = make(chan fsnotify.Event)
		c[i] = group
	}
	return c
}

type Collections map[Asset]Collection

func NewCollections(css Collection, js Collection) Collections {
	return Collections{
		AssetCss: NewCollection(css),
		AssetJs:  NewCollection(js),
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

	// Events channel
	events chan fsnotify.Event `json:"-"`
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
	return g.AbsPath(g.Output)
}

// Determine the Result path and return the value
// TODO: This method will calculate the version hash
func (g *Group) ResultPath() string {
	g.Result = g.RootedPath(g.Output)
	return g.Result
}
