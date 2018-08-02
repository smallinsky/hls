package notify

import (
	"github.com/pkg/errors"
)

var (
	ErrNotImpl = errors.New("Function is not yet supported for OSX")
)

// Not implemented
func WatchDir(dir string) (chan string, error) {
	return nil, ErrNotImpl
}
