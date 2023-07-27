package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

type Config struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"accessKey"`
	SecretKey  string `yaml:"secretKey"`
	BucketName string `yaml:"bucketName"`
}

func NewMinioClient(cfg *Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, errors.Wrap(err, "minio.New")
	}

	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, cfg.BucketName)
		if errBucketExists != nil || !exists {
			return nil, errors.Wrap(errBucketExists, "minioClient.BucketExists")
		}
	}

	return minioClient, nil
}
