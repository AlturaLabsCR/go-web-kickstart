package config

import (
	"app/storage"
	"app/storage/fs"
	"app/storage/s3"
)

func InitStorage() (storage.ObjectStorage, error) {
	var empty storage.ObjectStorage

	// TODO: Initialize properly and allow using S3 as well

	_ = &fs.FSParams{
		// Root          string
		// Cache         cache.Cache
		// L2Cache cache.Cache // optional
		// MaxObjectSize int64
		// MaxBucketSize int64
	}

	_ = &s3.S3Params{
		// Client        *s3.Client
		// BucketName    string
		// Cache         cache.Cache
		// L2Cache       cache.Cache // optional
		// MaxObjectSize int64
		// MaxBucketSize int64
	}

	// fstore, err := fs.NewFS(fsparams)
	// if err != nil {
	// 	return empty, err
	// }

	// s3tore, err := s3.NewS3(s3params)
	// if err != nil {
	// 	return empty, err
	// }

	return empty, nil
}
