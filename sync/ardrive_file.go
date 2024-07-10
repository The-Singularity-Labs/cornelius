package sync

import (
	"time"
)

type ArdriveFile struct {
	Path         string
	Mimetype     string
	LastModified time.Time
}

type ArdriveFiles []ArdriveFile

type ArdriveFileInfo struct {
	AppName          string `json:"appName"`
	AppVersion       string `json:"appVersion"`
	ArFS             string `json:"arFS"`
	ContentType      string `json:"contentType"`
	DriveId          string `json:"driveId"`
	EntityType       string `json:"entityType"`
	Name             string `json:"name"`
	TxId             string `json:"txId"`
	UnixTime         int64  `json:"unixTime"`
	Size             int64  `json:"size"`
	LastModifiedDate int64  `json:"lastModifiedDate"`
	DataTxId         string `json:"dataTxId"`
	DataContentType  string `json:"dataContentType"`
	ParentFolderId   string `json:"parentFolderId"`
	EntityId         string `json:"entityId"`
	Path             string `json:"path"`
	TxIdPath         string `json:"txIdPath"`
	EntityIdPath     string `json:"entityIdPath"`
}
