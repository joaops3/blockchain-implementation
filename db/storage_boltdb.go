package db

import (
	"errors"

	"github.com/boltDb/bolt"
)

type BoltStorage struct {
	db *bolt.DB
}

func NewBoltStorage(path string) (*BoltStorage, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &BoltStorage{db: db}, nil
}

func (b *BoltStorage) View(fn func(ReadBucket) error) error {
	return b.db.View(func(tx *bolt.Tx) error {
		return fn(&boltReadBucket{tx: tx})
	})
}

func (b *BoltStorage) Update(fn func(WriteBucket) error) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return fn(&boltWriteBucket{tx: tx})
	})
}

func (b *BoltStorage) Close() error {
	return b.db.Close()
}

type boltReadBucket struct {
	tx *bolt.Tx
}

func (r *boltReadBucket) Get(bucketName string, key []byte) ([]byte, error) {
	b := r.tx.Bucket([]byte(bucketName))
	if b == nil {
		return nil, errors.New("bucket not found")
	}
	val := b.Get(key)
	if val == nil {
		return nil, errors.New("key not found")
	}
	return val, nil
}


func (r *boltReadBucket) ForEach(bucketName string, fn func(k, v []byte) error) error {
	b := r.tx.Bucket([]byte(bucketName))
	if b == nil {
		return errors.New("bucket not found")
	}
	return b.ForEach(fn)
}


type boltWriteBucket struct {
	tx *bolt.Tx
}

func (w *boltWriteBucket) ForEach(bucketName string, fn func(k, v []byte) error) error {
	b := w.tx.Bucket([]byte(bucketName))
	if b == nil {
		return errors.New("bucket not found")
	}
	return b.ForEach(fn)
}

func (w *boltWriteBucket) Get(bucketName string, key []byte) ([]byte, error) {
	b := w.tx.Bucket([]byte(bucketName))
	if b == nil {
		return nil, errors.New("bucket not found")
	}
	val := b.Get(key)
	if val == nil {
		return nil, errors.New("key not found")
	}
	return val, nil
}

func (w *boltWriteBucket) Put(bucketName string, key []byte, value []byte) error {
	b := w.tx.Bucket([]byte(bucketName))
	if b == nil {
		return errors.New("bucket not found")
	}
	return b.Put(key, value)
}

func (w *boltWriteBucket) CreateBucketIfNotExists(bucketName string) error {
	_, err := w.tx.CreateBucketIfNotExists([]byte(bucketName))
	return err
}

func (w *boltWriteBucket) Delete(bucketName string, key []byte) error {
	b := w.tx.Bucket([]byte(bucketName))
	if b == nil {
		return errors.New("bucket not found")
	}
	return b.Delete(key)
}