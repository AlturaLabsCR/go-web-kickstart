package config

import (
	"context"
	"fmt"

	"app/cache"
	"app/database"
	"app/storage"
	"app/storage/fs"
	"app/storage/s3"

	sdkconfig "github.com/aws/aws-sdk-go-v2/config"
	sdk "github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	fsType = "fs"
	s3Type = "s3"
)

func InitStorage(ctx context.Context, querier database.Querier) (storage.ObjectStorage, error) {
	var empty storage.ObjectStorage

	switch Config.Storage.Type {
	case fsType:
		return initFS(querier)
	case s3Type:
		return initS3(ctx, querier)
	}

	return empty, fmt.Errorf("'%s' is not a valid storage type", Config.Storage.Type)
}

func initFS(querier database.Querier) (*fs.FS, error) {
	params := &fs.FSParams{
		Root:          Config.Storage.FS.Root,
		Cache:         cache.NewMemoryStore(),
		L2Cache:       querier,
		MaxObjectSize: Config.Storage.MaxObjectSize,
		MaxBucketSize: Config.Storage.MaxBucketSize,
	}

	fsStore, err := fs.NewFS(params)
	if err != nil {
		return nil, err
	}

	return fsStore, nil
}

func initS3(ctx context.Context, querier database.Querier) (*s3.S3Bucket, error) {
	cfg, err := sdkconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := sdk.NewFromConfig(cfg)

	params := &s3.S3Params{
		Client:        client,
		BucketName:    Config.Storage.S3.Bucket,
		Cache:         cache.NewMemoryStore(),
		L2Cache:       querier,
		MaxObjectSize: Config.Storage.MaxObjectSize,
		MaxBucketSize: Config.Storage.MaxBucketSize,
	}

	s3Store, err := s3.NewS3(params)
	if err != nil {
		return nil, err
	}

	return s3Store, nil
}
