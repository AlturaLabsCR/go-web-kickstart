package config

import (
	"app/storage/kv"
	"app/storage/s3"
)

func InitStorage(store kv.Store[s3.Object]) s3.Storage {
	var storage s3.Storage
	var err error

	if Config.Storage.Type != "remote" {
		storage = s3.NewFS(s3.S3Params{
			Bucket:        Config.Storage.Bucket,
			Store:         store,
			MaxObjectSize: Config.Storage.MaxObjectSize,
			MaxBucketSize: Config.Storage.MaxBucketSize,
		})
	} else {
		storage, err = s3.New(s3.S3Params{
			Bucket:        Config.Storage.Bucket,
			Store:         store,
			MaxObjectSize: Config.Storage.MaxObjectSize,
			MaxBucketSize: Config.Storage.MaxBucketSize,
		})
		if err != nil {
			panic("error setting up s3 client")
		}
	}

	return storage
}
