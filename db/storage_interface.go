package db

type Storage interface {
	View(fn func(ReadBucket) error) error
	Update(fn func(WriteBucket) error) error
	Close() error
}

type ReadBucket interface {
	Get(bucketName string, key []byte) ([]byte, error)
	ForEach(bucketName string, fn func(k, v []byte) error) error
}

type WriteBucket interface {
	ReadBucket
	Put(bucketName string, key []byte, value []byte) error
	Delete(bucketName string, key []byte) error
	CreateBucketIfNotExists(bucketName string) error
}