package sync

type Pipeline struct {
	Name             string           `yaml:"name"`
	Bucket           Bucket           `yaml:"bucket"`
	DestinationDrive DestinationDrive `yaml:"drive"`
	EnableManifest   bool             `yaml:"enable_manifest"`
	Frequency        Duration         `yaml:"frequency"`
}

type Bucket struct {
	Name        string `yaml:"name"`
	Host        string `yaml:"host"`
	Prefix      string `yaml:"prefix"`
	AccessId    string `yaml:"access_id"`
	SecretKey   string `yaml:"secret_key"`
	IsSecure    bool   `yaml:"is_secure"`
	IsRecursive bool   `yaml:"is_recursive"`
}

type DestinationDrive struct {
	Id             string `yaml:"id"`
	WalletPath     string `yaml:"wallet_path"`
	Password       string `yaml:"password"`
	ParentFolderId string `yaml:"parent_folder_id"`
	IsPublic       bool   `yaml:"is_public"`
}
