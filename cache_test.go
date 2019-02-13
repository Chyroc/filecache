package filecache_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Chyroc/filecache"
)

func TestNew(t *testing.T) {
	as := assert.New(t)

	as.Nil(os.Remove("./test"))
	c := filecache.New("./test").(*filecache.CacheImpl)

	t.Run("", func(t *testing.T) {
		as.Equal(5242880, c.CurrentSize)
	})

	t.Run("not found", func(t *testing.T) {
		v, err := c.Get("k1")
		as.Equal(filecache.NotFound, err)
		as.Equal("", v)
	})

	t.Run("exist get set", func(t *testing.T) {
		as.Nil(c.Set("k", "v", time.Second))

		v, err := c.Get("k")
		as.Nil(err)
		as.Equal("v", v)
	})

	t.Run("expired", func(t *testing.T) {
		as.Nil(c.Set("k", "v", time.Second))

		time.Sleep(time.Second)

		v, err := c.Get("k")
		as.Equal(filecache.NotFound, err)
		as.Equal("", v)
	})

	t.Run("ttl", func(t *testing.T) {
		as.Nil(c.Set("k", "v", time.Second))

		ttl, err := c.TTL("k")
		as.Nil(err)
		as.True(ttl <= time.Second && ttl >= time.Second-10*time.Millisecond)
	})
}

func BenchmarkCacheImpl_Get(b *testing.B) {
	as := assert.New(b)

	as.Nil(os.Remove("./test"))
	c := filecache.New("./test")

	for i := 0; i < b.N; i++ {
		j := strconv.Itoa(i)
		as.Nil(c.Set(j, j, time.Second), i)
	}
}
