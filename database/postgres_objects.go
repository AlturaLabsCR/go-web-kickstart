package database

import (
	"context"

	"app/database/postgres/db"
	"app/storage/s3"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresObjectStore struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

func NewPostgresObjectStore(s *Postgres) *PostgresObjectStore {
	return &PostgresObjectStore{
		Pool:    s.Pool,
		Queries: s.Queries,
	}
}

func (p *PostgresObjectStore) Set(ctx context.Context, key string, data s3.Object) error {
	return p.Queries.UpsertObject(ctx, db.UpsertObjectParams{
		ObjectKey:    key,
		ObjectBucket: data.Bucket,
		ObjectMime:   data.Mime,
		ObjectMd5:    data.MD5,
		ObjectSize:   data.Size,
	})
}

func (p *PostgresObjectStore) Get(ctx context.Context, key string) (s3.Object, error) {
	obj, err := p.Queries.SelectObject(ctx, key)
	if err != nil {
		return s3.Object{}, err
	}

	return s3.Object{
		Key:      obj.ObjectKey,
		Bucket:   obj.ObjectBucket,
		Mime:     obj.ObjectMime,
		MD5:      obj.ObjectMd5,
		Size:     obj.ObjectSize,
		Created:  obj.ObjectCreated.Time,
		Modified: obj.ObjectModified.Time,
	}, nil
}

func (p *PostgresObjectStore) Delete(ctx context.Context, key string) error {
	return p.Queries.DeleteObject(ctx, key)
}

func (p *PostgresObjectStore) GetElems(ctx context.Context) (map[string]s3.Object, error) {
	dbObjects, err := p.Queries.GetObjects(ctx)
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
			Created:  obj.ObjectCreated.Time,
			Modified: obj.ObjectModified.Time,
		}
	}

	return obs, nil
}
