package s3

import (
	wfv1 "github.com/argoproj/argo/api/workflow/v1alpha1"
	"github.com/argoproj/argo/errors"
	minio "github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

// S3ArtifactDriver is a driver for AWS S3
type S3ArtifactDriver struct {
	Endpoint  string
	Secure    bool
	AccessKey string
	SecretKey string
}

// newMinioClient instantiates a new minio client object.
func (s3Driver *S3ArtifactDriver) newMinioClient() (*minio.Client, error) {
	minioClient, err := minio.New(s3Driver.Endpoint, s3Driver.AccessKey, s3Driver.SecretKey, s3Driver.Secure)
	if err != nil {
		return nil, errors.InternalWrapError(err)
	}
	return minioClient, nil
}

// Load downloads artifacts from S3 compliant storage using Minio Go SDK
func (s3Driver *S3ArtifactDriver) Load(inputArtifact *wfv1.Artifact, path string) error {
	minioClient, err := s3Driver.newMinioClient()
	if err != nil {
		return err
	}
	// Download the file to a local file path
	log.Infof("Loading from s3 (endpoint: %s, bucket: %s, key: %s) to %s",
		inputArtifact.S3.Endpoint, inputArtifact.S3.Bucket, inputArtifact.S3.Key, path)
	err = minioClient.FGetObject(inputArtifact.S3.Bucket, inputArtifact.S3.Key, path)
	if err != nil {
		return errors.InternalWrapError(err)
	}
	return nil
}

func (s3Driver *S3ArtifactDriver) Save(path string, outputArtifact *wfv1.Artifact) error {
	minioClient, err := s3Driver.newMinioClient()
	if err != nil {
		return err
	}
	log.Infof("Saving from %s to s3 (endpoint: %s, bucket: %s, key: %s)",
		path, outputArtifact.S3.Endpoint, outputArtifact.S3.Bucket, outputArtifact.S3.Key)

	_, err = minioClient.FPutObject(outputArtifact.S3.Bucket, outputArtifact.S3.Key, path, "application/gzip")
	if err != nil {
		return errors.InternalWrapError(err)
	}
	return nil
}
