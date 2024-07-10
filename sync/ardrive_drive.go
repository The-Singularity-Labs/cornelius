package sync

type ArdriveDrive struct {
	AppName               string            `json:"appName"`
	AppVersion            string            `json:"appVersion"`
	ArFS                  string            `json:"arFS"`
	ContentType           string            `json:"contentType"`
	DriveId               string            `json:"driveId"`
	EntityType            string            `json:"entityType"`
	Name                  string            `json:"name"`
	TxId                  string            `json:"txId"`
	UnixTime              int64             `json:"unixTime"`
	CustomMetaDataGqlTags map[string]string `json:"customMetaDataGqlTags"`
	DrivePrivacy          string            `json:"drivePrivacy"`
	RootFolderId          string            `json:"rootFolderId"`
}

type ArdriveDrives []ArdriveDrive
