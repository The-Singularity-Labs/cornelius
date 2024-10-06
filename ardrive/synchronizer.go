package sync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/the-singularity-labs/cornelius/log"

	"golang.org/x/sync/errgroup"
)

type Synchronizer struct {
	ardrivecliPath string
	config         Config
	logger         log.Logger
}

func New(logger log.Logger, ardrivecliPath string, config Config) *Synchronizer {
	return &Synchronizer{
		logger:         logger,
		ardrivecliPath: ardrivecliPath,
		config:         config,
	}
}

func (s *Synchronizer) Start(ctx context.Context) error {
	var g errgroup.Group

	s.logger.Info("initializing pipelines", "count", len(s.config.Pipelines))
	for _, pipeline := range s.config.Pipelines {
		g.Go(func() error {
			return s.handlePipeline(ctx, pipeline)
		})
	}

	s.logger.Info("all pipelines initialized")

	return g.Wait()
}

func (s *Synchronizer) handlePipeline(ctx context.Context, pipeline Pipeline) error {
	logger := s.logger.With("pipeline", pipeline.Name)

	objConn, err := NewObjectStorageConnection(ctx, logger, s.config.TmpDirectory, pipeline.Bucket.Host, pipeline.Bucket.Name, pipeline.Bucket.Prefix, pipeline.Bucket.AccessId, pipeline.Bucket.SecretKey, pipeline.Bucket.IsSecure, pipeline.Bucket.IsRecursive)
	if err != nil {
		return fmt.Errorf("unable to initialize object storage connection %q: %w", pipeline.Name, err)
	}

	ardriveClient, err := NewArdriveClient(logger, s.ardrivecliPath, pipeline.DestinationDrive.WalletPath, pipeline.DestinationDrive.Password, pipeline.DestinationDrive.Id, pipeline.DestinationDrive.ParentFolderId, pipeline.DestinationDrive.IsPublic, pipeline.EnableManifest)
	if err != nil {
		return fmt.Errorf("unable to initialize ardrive client %q: %w", pipeline.Name, err)
	}

	exists, err := ardriveClient.DriveExists()
	if err != nil {
		return fmt.Errorf("unable to check if ardrive drive exists for pipeline %q: %w", pipeline.Name, err)
	} else if !exists {
		return fmt.Errorf("ardrive with ID %q does not exist for pipeline %q", pipeline.DestinationDrive.Id, pipeline.Name)
	}

	parentPath, err := ardriveClient.GetParentPath()
	if err != nil {
		return fmt.Errorf("unable to check if ardrive folder exists for pipeline %q: %w", pipeline.Name, err)
	} else if parentPath == "" {
		return fmt.Errorf("ardrive folder with id %q does not exist for pipeline %q", pipeline.DestinationDrive.ParentFolderId, pipeline.Name)
	}

	repeatOnSetFrequency := true
	sleepDuration := time.Duration(pipeline.Frequency)
	if sleepDuration == time.Duration(0) {
		logger.Info("no frequency set, proces will exit after first iteration.")
		repeatOnSetFrequency = false
	}

	logger.Info("starting sync")
	for {
		logger.Info("getting existing files")
		objectStorageFiles, err := objConn.ListFiles()
		if err != nil {
			return fmt.Errorf("unable to get files to sync: %w", err)
		}

		logger.Info("acquired object storage files", "count", len(objectStorageFiles))

		ardriveFiles, err := ardriveClient.ListFiles()
		if err != nil {
			return fmt.Errorf("unable to get drives to sync: %w", err)
		}

		logger.Info("acquired ardrive files", "count", len(ardriveFiles))

		deltaObjectStorageFiles, err := identifyNetNewFiles(objectStorageFiles, ardriveFiles, parentPath)
		if err != nil {
			return fmt.Errorf("unable to compare object storage files to ardrive files: %w", err)
		}
		logger.Info("idenitifed files to sync", "count", len(deltaObjectStorageFiles))

		for _, objectStorageFileToSync := range deltaObjectStorageFiles {
			err := func() error {
				logger := logger.With("object", objectStorageFileToSync.Key)
				logger.Debug("downloading file from object storage")
				localFile, err := objConn.DownloadFile(objectStorageFileToSync)
				if err != nil {
					return fmt.Errorf("unable to download object %q in order to reupload to arweave: %w", objectStorageFileToSync.Key, err)
				}

				defer func() {
					logger.Debug("removing staged file")
					removeLocalFile(localFile)
				}()

				logger.Debug("finished dowloading file from object storage", "path", localFile.Path)

				txData, err := ardriveClient.upsertFile(localFile) // TODO: compile response statistics intometrics
				if err != nil {
					return fmt.Errorf("unable to upsert %q to arweave: %w", localFile.Dir, err)
				}

				totalFees, err := txData.TotalFees()
				if err != nil {
					return fmt.Errorf("unable to upsert %q to arweave: %w", strings.Join(txData.EntityIds(), ", "), err)
				}

				logger.Info("file uploaded to arweave", "fees_paid", totalFees)

				return nil
			}()
			if err != nil {
				return err
			}
		}

		if !repeatOnSetFrequency {
			break
		}

		logger.Info("sleeping inside pipeline", "duration", sleepDuration)
		time.Sleep(sleepDuration)
	}

	return nil
}

func identifyNetNewFiles(objectStorageFiles ObjectStorageFiles, ardriveFiles ArdriveFiles, parentPath string) (ObjectStorageFiles, error) {
	filtered := ObjectStorageFiles{}

	objectStorageFileMap := map[string]ObjectStorageFile{}
	for _, objectStorageFile := range objectStorageFiles {
		objectStorageFileMap[objectStorageFile.Key] = objectStorageFile
	}

	ardriveFileMap := map[string]ArdriveFile{}
	for _, ardriveFile := range ardriveFiles {
		ardriveFileMap[ardriveFile.Path] = ardriveFile
	}

	for key, objectStorageFile := range objectStorageFileMap {
		potentialPath := filepath.Join(parentPath, key)

		if ardriveFile, exists := ardriveFileMap[potentialPath]; !exists || objectStorageFile.LastModified.After(ardriveFile.LastModified) {
			filtered = append(filtered, objectStorageFile)
		}
	}

	return filtered, nil
}

func removeLocalFile(localFile LocalFile) {
	os.Remove(localFile.Path)
}
