package sync

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/the-singularity-labs/cornelius/log"
)

const ArdriveCliFileSizeLimit = 2000000000

type ArdriveClient struct {
	logger         log.Logger
	executablePath string
	driveId        string
	isPublic       bool
	enableManifest bool
	parentFolderId string
	walletPath     string
	walletPassword string
}

func NewArdriveClient(logger log.Logger, executablePath, walletPath, walletPassword, driveId, parentFolderId string, isPublic, enableManifest bool) (*ArdriveClient, error) {
	return &ArdriveClient{
		logger:         logger,
		executablePath: executablePath,
		driveId:        driveId,
		parentFolderId: parentFolderId,
		isPublic:       isPublic,
		walletPath:     walletPath,
		walletPassword: walletPassword,
		enableManifest: enableManifest,
	}, nil
}

func (client *ArdriveClient) exec(args ...string) ([]byte, error) {
	args = append(args, []string{"-w", client.walletPath, "--unsafe-drive-password", client.walletPassword}...)
	client.logger.Info(client.executablePath, "args", args)
	resp, err := ExecCmd(client.executablePath, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to exec private ardrive cli command: %w", err)
	}

	return resp, nil

}

func (client *ArdriveClient) execPrivateOrPublic(args ...string) ([]byte, error) {
	if !client.isPublic {
		args = append(args, []string{"-w", client.walletPath, "--unsafe-drive-password", client.walletPassword}...)
	}

	resp, err := ExecCmd(client.executablePath, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to exec ardrive cli command: %w", err)
	}

	return resp, nil
}

func (client *ArdriveClient) ListDrives() (ArdriveDrives, error) {
	resp, err := client.exec("list-all-drives")
	if err != nil {
		return nil, fmt.Errorf("unable to list ardrive drives: %w", err)
	}

	results := ArdriveDrives{}
	err = json.Unmarshal(resp, &results)
	if err != nil {
		return nil, fmt.Errorf("unable to parse list-all-drives response: %w", err)
	}

	return results, nil
}

func (client *ArdriveClient) DriveExists() (bool, error) {
	drives, err := client.ListDrives()
	if err != nil {
		return false, fmt.Errorf("unable to get drives to check for id %q: %w", client.driveId, err)
	}

	for _, drive := range drives {
		if drive.DriveId == client.driveId {
			return true, nil
		}
	}

	return false, nil
}

func (client *ArdriveClient) getFolderPath(folderId string) (ArdriveFolderInfo, error) {
	resp, err := client.execPrivateOrPublic("folder-info", "--folder-id", folderId)
	if err != nil {
		return ArdriveFolderInfo{}, fmt.Errorf("unable to get parent folder: %w", err)
	}

	ardriveFolderInfo := ArdriveFolderInfo{}
	err = json.Unmarshal(resp, &ardriveFolderInfo)
	if err != nil {
		return ArdriveFolderInfo{}, fmt.Errorf("unable to parse file-info response: %w", err)
	}

	return ardriveFolderInfo, nil
}

func (client *ArdriveClient) recursiveGenerateFolderPath(folderId string) (string, error) {
	ardriveFolderInfo, err := client.getFolderPath(folderId)
	if err != nil {
		return "", err
	}

	if ardriveFolderInfo.ParentFolderId == "root folder" {
		absRootPath := "/" + ardriveFolderInfo.Name
		return absRootPath, nil
	}

	childPath, err := client.recursiveGenerateFolderPath(ardriveFolderInfo.ParentFolderId)
	if err != nil {
		return "", err
	}

	return path.Join(ardriveFolderInfo.Name, childPath), nil
}

func (client *ArdriveClient) GetParentPath() (string, error) {
	return client.recursiveGenerateFolderPath(client.parentFolderId)
}

func (client *ArdriveClient) ListFiles() (ArdriveFiles, error) {
	resp, err := client.execPrivateOrPublic("list-folder", "--parent-folder-id", client.parentFolderId, "--all")
	if err != nil {
		return nil, fmt.Errorf("unable to list ardrive files: %w", err)
	}

	results := []ArdriveFileInfo{}
	err = json.Unmarshal(resp, &results)
	if err != nil {
		return nil, fmt.Errorf("unable to parse list-folder response: %w", err)
	}

	foundFiles := ArdriveFiles{}
	for _, ardrivefileInfo := range results {
		foundFiles = append(foundFiles, ArdriveFile{
			Path:         ardrivefileInfo.Path,
			Mimetype:     ardrivefileInfo.DataContentType,
			LastModified: time.Unix(ardrivefileInfo.LastModifiedDate, 0),
		})
	}

	return foundFiles, nil
}

func (client *ArdriveClient) upsertFile(localFile LocalFile) (TxData, error) {
	args := []string{
		"upload-file",
		"--parent-folder-id",
		client.parentFolderId,
		"--local-path",
		localFile.Dir,
	}

	if localFile.Mimetype != "" {
		args = append(args, "--content-type", localFile.Mimetype)
	}

	resp, err := client.exec(args...)
	if err != nil {
		return TxData{}, fmt.Errorf("unable to upsert ardrive file: %w", err)
	}

	results := TxData{}
	err = json.Unmarshal(resp, &results)
	if err != nil {
		return TxData{}, fmt.Errorf("unable to parse upload-file response: %w", err)
	}

	if client.enableManifest && localFile.Filename() == "index.html" {
		err = client.createManifest(results.EntityId())
		if err != nil {
			return TxData{}, fmt.Errorf("unable to create corresponding manifest response: %w", err)
		}
	}

	return results, nil
}

func (client *ArdriveClient) createManifest(indexEntityId string) error {
	args := []string{
		"file-info",
		"--file-id",
		indexEntityId,
	}

	resp, err := client.exec(args...)
	if err != nil {
		return fmt.Errorf("unable to upsert ardrive file: %w", err)
	}

	indexEntity := ArdriveFileInfo{}
	err = json.Unmarshal(resp, &indexEntity)
	if err != nil {
		return fmt.Errorf("unable to parse file-info response: %w", err)
	}

	args = []string{
		"create-manifest",
		"--f",
		indexEntity.ParentFolderId,
	}

	client.logger.Info("creating manifest", "parent_id", indexEntity.ParentFolderId, "index_entity_id", indexEntityId)
	_, err = client.exec(args...)
	if err != nil {
		return fmt.Errorf("unable to create manifest file: %w", err)
	}

	return nil
}

func pathWithoutRootFolder(fullPath string) string {
	base := path.Base(fullPath)
	return path.Join(path.Dir(fullPath)[1:], base)
}
