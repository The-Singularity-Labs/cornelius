package sync

import (
	"path/filepath"
)

type LocalFile struct {
	Dir      string
	Path     string
	Mimetype string
}

func (lf LocalFile) Filename() string {
	return filepath.Base(lf.Path)
}
