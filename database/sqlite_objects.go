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
		ObjectBucket: data.ObjectBucket,
		ObjectMime:   data.ObjectMime,
		ObjectSize:   data.ObjectSize,
	})
}

func (s *SqliteObjectStore) Get(ctx context.Context, key string) (s3.Object, error) {
	obj, err := s.Queries.SelectObject(ctx, key)
	if err != nil {
		return s3.Object{}, err
	}

	return s3.Object{
		ObjectKey:      obj.ObjectKey,
		ObjectBucket:   obj.ObjectBucket,
		ObjectMime:     obj.ObjectMime,
		ObjectSize:     obj.ObjectSize,
		ObjectCreated:  time.Unix(obj.ObjectCreated, 0),
		ObjectModified: time.Unix(obj.ObjectModified, 0),
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

	for _, ob := range dbObjects {
		obs[ob.ObjectKey] = s3.Object{
			ObjectKey:      ob.ObjectKey,
			ObjectBucket:   ob.ObjectBucket,
			ObjectMime:     ob.ObjectMime,
			ObjectSize:     ob.ObjectSize,
			ObjectCreated:  time.Unix(ob.ObjectCreated, 0),
			ObjectModified: time.Unix(ob.ObjectModified, 0),
		}
	}

	return obs, nil
}
