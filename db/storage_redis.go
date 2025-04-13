package db

import (
	"context"
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(addr string, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStorage{
		client: client,
		ctx:    context.Background(),
	}, nil
}

func (r *RedisStorage) View(fn func(ReadBucket) error) error {
	return fn(&redisBucket{client: r.client, ctx: r.ctx})
}

func (r *RedisStorage) Update(fn func(WriteBucket) error) error {
	return fn(&redisBucket{client: r.client, ctx: r.ctx})
}

func (r *RedisStorage) Close() error {
	return r.client.Close()
}


type redisBucket struct {
	client *redis.Client
	ctx    context.Context
}

func (rb *redisBucket) makeKey(bucket string, key []byte) string {
	return bucket + ":" + string(key)
}

func (rb *redisBucket) Get(bucketName string, key []byte) ([]byte, error) {
	val, err := rb.client.Get(rb.ctx, rb.makeKey(bucketName, key)).Bytes()
	if err == redis.Nil {
		return nil, errors.New("key not found")
	}
	return val, err
}

func (rb *redisBucket) Put(bucketName string, key []byte, value []byte) error {
	return rb.client.Set(rb.ctx, rb.makeKey(bucketName, key), value, 0).Err()
}

func (rb *redisBucket) Delete(bucketName string, key []byte) error {
	return rb.client.Del(rb.ctx, rb.makeKey(bucketName, key)).Err()
}

func (rb *redisBucket) CreateBucketIfNotExists(bucketName string) error {
	return nil
}

func (rb *redisBucket) ForEach(bucketName string, fn func(k, v []byte) error) error {
	var cursor uint64
	prefix := bucketName + ":"

	for {
		keys, newCursor, err := rb.client.Scan(rb.ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}
		cursor = newCursor

		for _, fullKey := range keys {
			val, err := rb.client.Get(rb.ctx, fullKey).Bytes()
			if err != nil && err != redis.Nil {
				return err
			}
			key := []byte(strings.TrimPrefix(fullKey, prefix))
			if err := fn(key, val); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}


