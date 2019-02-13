package filecache

import "time"

type Cache interface {
	Get(key string) (string, error)
	Set(key, val string, ttl time.Duration) error
	TTL(key string) (time.Duration, error)
}

func NewDefault(filepath string) Cache {
	return New(filepath, 20)
}

// maxSize Mb
func New(filepath string, maxSize int) Cache {
	return &cacheImpl{
		filepath: filepath,
		maxSize:  maxSize * 1024 * 1024, // B
	}
}

type cacheImpl struct {
	filepath string
	maxSize  int
}

func (r *cacheImpl) Get(key string) (string, error) {
	panic("implement me")
}

func (r *cacheImpl) Set(key, val string, ttl time.Duration) error {
	panic("implement me")
}

func (r *cacheImpl) TTL(key string) (time.Duration, error) {
	panic("implement me")
}
