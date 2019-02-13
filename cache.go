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

// MaxSize Mb
func New(filepath string, maxSize int) Cache {
	b := maxSize * 1024 * 1024
	mod := b / 1024 / 2 / 10 // 2Kb 一个 entry，10个 entry一个区域
	return &CacheImpl{
		filepath: filepath,
		MaxSize:  1024 * 2 * mod, // B
		Mod:      mod,
	}
}

// every entry: 2Kb: key(256B), ttl(ms,13位数字,7B), value(1024*2-256-7 = 1785B)
type CacheImpl struct {
	filepath string
	MaxSize  int
	Mod      int
}

func (r *CacheImpl) Get(key string) (string, error) {
	panic("implement me")
}

func (r *CacheImpl) Set(key, val string, ttl time.Duration) error {
	panic("implement me")
}

func (r *CacheImpl) TTL(key string) (time.Duration, error) {
	panic("implement me")
}
