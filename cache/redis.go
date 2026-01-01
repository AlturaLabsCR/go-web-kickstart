package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	r *redis.Client
}

var _ Cache = (*RedisStore)(nil)

func NewRedisStore(opt *redis.Options) *RedisStore {
	return &RedisStore{r: redis.NewClient(opt)}
}

func (s *RedisStore) Set(ctx context.Context, key, value string) error {
	return s.r.Set(ctx, key, value, 0).Err()
}

func (s *RedisStore) Get(ctx context.Context, key string) (string, error) {
	val, err := s.r.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrNotFound
		}
		return "", err
	}
	return val, nil
}

func (s *RedisStore) Del(ctx context.Context, key string) error {
	return s.r.Del(ctx, key).Err()
}

func (s *RedisStore) GetAll(ctx context.Context) (map[string]string, error) {
	values := make(map[string]string)

	var cursor uint64
	for {
		keys, nextCursor, err := s.r.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			return nil, err
		}

		if len(keys) > 0 {
			vals, err := s.r.MGet(ctx, keys...).Result()
			if err != nil {
				return nil, err
			}

			for i, v := range vals {
				if v == nil {
					continue
				}
				values[keys[i]] = v.(string)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return values, nil
}
