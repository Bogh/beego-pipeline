package pipeline

import (
	"io"
)

// Reader that closes it self when reaches EOF
type AutoCloseReader struct {
	io.ReadCloser
}

func (a *AutoCloseReader) Read(p []byte) (n int, err error) {
	n, err = a.ReadCloser.Read(p)
	if err == io.EOF {
		defer a.Close()
	}
	return
}
