package file

import (
	"time"
)

type Info struct {
	Name    string
	Size    int64
	ModTime time.Time
}
