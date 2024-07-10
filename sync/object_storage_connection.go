package sync

import (
	"context"
	"fmt"
	"math/rand"
	"path"
	"path/filepath"
	"time"

	"github.com/the-singularity-labs/cornelius/log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type ObjectStorageConnection struct {
	ctx          context.Context
	minioClient  *minio.Client
	bucket       string
	prefix       string
	isRecursive  bool
	logger       log.Logger
	tmpDirectory string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewObjectStorageConnection(ctx context.Context, logger log.Logger, tmpDirectory, host, bucket, prefix, accessId, secretKey string, isSecure, isRecursive bool) (*ObjectStorageConnection, error) {
	minioClient, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessId, secretKey, ""), // TODO: add support for temp creds
		Secure: isSecure,
	})

	if err != nil {
		return nil, fmt.Errorf("unable to initialize object storage client: %w", err)
	}

	return &ObjectStorageConnection{
		minioClient:  minioClient,
		bucket:       bucket,
		prefix:       prefix,
		isRecursive:  isRecursive,
		logger:       logger,
		tmpDirectory: path.Join(tmpDirectory, randCharSeq(5)),
	}, nil
}

func (conn *ObjectStorageConnection) ListFiles() (ObjectStorageFiles, error) {
	opts := minio.ListObjectsOptions{
		Recursive: conn.isRecursive,
		Prefix:    conn.prefix,
	}

	results := ObjectStorageFiles{}
	for objectInfo := range conn.minioClient.ListObjects(context.Background(), conn.bucket, opts) {

		if objectInfo.Err != nil {
			return nil, fmt.Errorf("unable to iterate through objects: %w", objectInfo.Err)
		}

		logger := conn.logger.With("key", objectInfo.Key)
		if objectInfo.Size == 0 {
			logger.Warn("skipping file, file is empty and is likely just a folder")
			continue
		} else if objectInfo.Size > ArdriveCliFileSizeLimit {
			logger.Warn("skipping file, exceeds 2GB limit")
			continue
		}

		lastModified := objectInfo.LastModified
		if lastModified.IsZero() {
			logger.Debug("Detected last modified tme of epoch 0. This is likely a brand new file. Setting LastModified to now()")
			lastModified = time.Now()
		}

		results = append(results, ObjectStorageFile{
			Key:          objectInfo.Key,
			LastModified: lastModified,
			Mimetype:     objectInfo.ContentType,
		})
	}

	return results, nil
}

func (conn *ObjectStorageConnection) DownloadFile(objectStorageFile ObjectStorageFile) (LocalFile, error) {
	localFilePath := objectStorageFile.Key

	err := conn.minioClient.FGetObject(context.Background(), conn.bucket, objectStorageFile.Key, localFilePath, minio.GetObjectOptions{})
	if err != nil {
		return LocalFile{}, fmt.Errorf("unable to download file from object storage: %w", err)
	}

	return LocalFile{
		Dir:      filepath.Dir(localFilePath),
		Path:     localFilePath,
		Mimetype: objectStorageFile.Mimetype,
	}, nil
}

func randCharSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
