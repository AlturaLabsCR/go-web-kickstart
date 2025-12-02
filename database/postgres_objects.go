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

// type Object struct {
// 	ObjectBucket   string
// 	ObjectKey      string
// 	ObjectMime     string
// 	ObjectSize     int64
// 	ObjectCreated  time.Time
// 	ObjectModified time.Time
// }

func (p *PostgresObjectStore) Set(ctx context.Context, key string, data s3.Object) error {
	return p.Queries.UpsertObject(ctx, db.UpsertObjectParams{
		ObjectKey:    key,
		ObjectBucket: data.ObjectBucket,
		ObjectMime:   data.ObjectMime,
		ObjectSize:   data.ObjectSize,
	})
}

func (p *PostgresObjectStore) Get(ctx context.Context, key string) (s3.Object, error) {
	obj, err := p.Queries.SelectObject(ctx, key)
	if err != nil {
		return s3.Object{}, err
	}

	return s3.Object{
		ObjectKey:      obj.ObjectKey,
		ObjectBucket:   obj.ObjectBucket,
		ObjectMime:     obj.ObjectMime,
		ObjectSize:     obj.ObjectSize,
		ObjectCreated:  obj.ObjectCreated.Time,
		ObjectModified: obj.ObjectModified.Time,
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

	for _, ob := range dbObjects {
		obs[ob.ObjectKey] = s3.Object{
			ObjectKey:      ob.ObjectKey,
			ObjectBucket:   ob.ObjectBucket,
			ObjectMime:     ob.ObjectMime,
			ObjectSize:     ob.ObjectSize,
			ObjectCreated:  ob.ObjectCreated.Time,
			ObjectModified: ob.ObjectModified.Time,
		}
	}

	return obs, nil
}
