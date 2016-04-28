package pipeline

import (
	"errors"
	"fmt"
)

var (
	ErrNoCompiler = errors.New("No compiler found")

	ErrAssetNotFound = errors.New("No asset found.")
)

type ErrOverridingPath struct {
	Path string
}

func (e *ErrOverridingPath) Error() string {
	return fmt.Sprintf("Overriding original path: %s", e.Path)
}
