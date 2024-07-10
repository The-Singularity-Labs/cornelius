package sync

import (
	"time"
)

type ObjectStorageFile struct {
	Key          string
	Mimetype     string
	LastModified time.Time
}

type ObjectStorageFiles []ObjectStorageFile
