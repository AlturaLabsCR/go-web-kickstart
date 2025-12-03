package config

import (
	"context"

	"app/storage/kv"
	"app/storage/s3"
)

func InitStorage(store kv.Store[s3.Object]) s3.Storage {
	var storage s3.Storage
	var err error

	switch Config.Storage.Type {
	case "local":
		storage = s3.NewFS(s3.S3Params{
			Bucket:         Config.Storage.Local.Root,
			Store:          store,
			MaxObjectSize:  Config.Storage.MaxObjectSize,
			MaxBucketSize:  Config.Storage.MaxBucketSize,
			PublicEndpoint: Config.Storage.Local.PublicEndpointURL,
		})
	case "remote":
		storage, err = s3.New(s3.S3Params{
			Bucket:         Config.Storage.Remote.Bucket,
			Store:          store,
			MaxObjectSize:  Config.Storage.MaxObjectSize,
			MaxBucketSize:  Config.Storage.MaxBucketSize,
			PublicEndpoint: Config.Storage.Remote.PublicEndpointURL,
		})
		if err != nil {
			panic("error setting up s3 client")
		}
	default:
		panic("error unsupported storage type")
	}

	if err := storage.LoadCache(context.Background()); err != nil {
		panic("error loading cache")
	}

	return storage
}
