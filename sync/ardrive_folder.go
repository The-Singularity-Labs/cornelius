package sync

type ArdriveFolderInfo struct {
	AppName        string `json:"appName"`
	AppVersion     string `json:"appVersion"`
	ArFS           string `json:"arFS"`
	ContentType    string `json:"contentType"`
	DriveId        string `json:"driveId"`
	EntityType     string `json:"entityType"`
	Name           string `json:"name"`
	TxId           string `json:"txId"`
	UnixTime       int64  `json:"unixTime"`
	ParentFolderId string `json:"parentFolderId"` // Note the potential string value
	EntityId       string `json:"entityId"`
	FolderId       string `json:"folderId"`
}
