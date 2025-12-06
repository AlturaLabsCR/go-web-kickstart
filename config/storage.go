package config

import (
	"context"

	"app/storage"
	"app/storage/kv"
	"app/storage/s3"
)

func InitStorage(store kv.Store[s3.Object]) storage.ObjectStorage {
	var storage storage.ObjectStorage
	var err error

	switch Config.Storage.Type {
	case "local":
		storage, err = s3.NewFS(&s3.StorageParams{
			BucketName:     Config.Storage.Local.Root,
			Store:          store,
			Cache:          kv.NewMemoryStore[s3.Object](),
			MaxObjectSize:  Config.Storage.MaxObjectSize,
			MaxBucketSize:  Config.Storage.MaxBucketSize,
			PublicEndpoint: Config.Storage.Local.PublicEndpointURL,
		})
		if err != nil {
			panic("error setting up filesystem storage client")
		}
	case "remote":
		storage, err = s3.NewS3(&s3.StorageParams{
			BucketName:     Config.Storage.Remote.Bucket,
			Store:          store,
			Cache:          kv.NewMemoryStore[s3.Object](),
			MaxObjectSize:  Config.Storage.MaxObjectSize,
			MaxBucketSize:  Config.Storage.MaxBucketSize,
			PublicEndpoint: Config.Storage.Remote.PublicEndpointURL,
		})
		if err != nil {
			panic("error setting up s3 storage client")
		}
	default:
		panic("error unsupported storage type")
	}

	if err := storage.LoadCache(context.Background()); err != nil {
		panic("error loading cache")
	}

	return storage
}
