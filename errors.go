package pipeline

import (
	"errors"
)

var (
	ErrNoCompiler = errors.New("No compiler found")

	ErrAssetNotFound = errors.New("No asset found.")
)
