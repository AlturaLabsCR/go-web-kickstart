package database

import (
	"context"
	"database/sql"
	"time"

	"app/database/sqlite/db"
	"app/storage/s3"
)

type SqliteObjectStore struct {
	DB      *sql.DB
	Queries *db.Queries
}

func NewSqliteObjectStore(s *Sqlite) *SqliteObjectStore {
	return &SqliteObjectStore{
		DB:      s.DB,
		Queries: s.Queries,
	}
}

func (s *SqliteObjectStore) Set(ctx context.Context, key string, data s3.Object) error {
	return s.Queries.UpsertObject(ctx, db.UpsertObjectParams{
		ObjectKey:    key,
		ObjectBucket: data.Bucket,
		ObjectMime:   data.Mime,
		ObjectMd5:    data.MD5,
		ObjectSize:   data.Size,
	})
}

func (s *SqliteObjectStore) Get(ctx context.Context, key string) (s3.Object, error) {
	obj, err := s.Queries.SelectObject(ctx, key)
	if err != nil {
		return s3.Object{}, err
	}

	return s3.Object{
		Key:      obj.ObjectKey,
		Bucket:   obj.ObjectBucket,
		Mime:     obj.ObjectMime,
		MD5:      obj.ObjectMd5,
		Size:     obj.ObjectSize,
		Created:  time.Unix(obj.ObjectCreated, 0),
		Modified: time.Unix(obj.ObjectModified, 0),
	}, nil
}

func (s *SqliteObjectStore) Delete(ctx context.Context, key string) error {
	return s.Queries.DeleteObject(ctx, key)
}

func (s *SqliteObjectStore) GetElems(ctx context.Context) (map[string]s3.Object, error) {
	dbObjects, err := s.Queries.GetObjects(ctx)
	if err != nil {
		return map[string]s3.Object{}, err
	}

	obs := map[string]s3.Object{}

	for _, obj := range dbObjects {
		obs[obj.ObjectKey] = s3.Object{
			Key:      obj.ObjectKey,
			Bucket:   obj.ObjectBucket,
			Mime:     obj.ObjectMime,
			MD5:      obj.ObjectMd5,
			Size:     obj.ObjectSize,
			Created:  time.Unix(obj.ObjectCreated, 0),
			Modified: time.Unix(obj.ObjectModified, 0),
		}
	}

	return obs, nil
}
