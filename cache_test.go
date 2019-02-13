package filecache_test

import (
	"os"
	"strconv"
	"strings"
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

	t.Run("invalid length", func(t *testing.T) {
		var err error
		long := strings.Repeat("x", 9999)

		as.Equal(filecache.KeyTooShort, c.Set("", "v", time.Second))
		as.Equal(filecache.ValueTooShort, c.Set("k", "", time.Second))
		as.Equal(filecache.KeyTooLong, c.Set(long, "v", time.Second))
		as.Equal(filecache.ValueTooLong, c.Set("k", long, time.Second))

		_, err = c.Get(long)
		as.Equal(filecache.KeyTooLong, err)

		_, err = c.Get("")
		as.Equal(filecache.KeyTooShort, err)
	})

	t.Run("too large", func(t *testing.T) {
		as.Nil(os.Remove("./test"))
		c = filecache.New("./test").(*filecache.CacheImpl)

		for i := 0; i <= 63124; i++ {
			j := strconv.Itoa(i)
			as.Nil(c.Set(j, j, time.Second), i)
		}

		as.Equal(filecache.FileSizeTooLarge, c.Set("63125", "63125", time.Second))

		as.Nil(os.Remove("./test"))
		c = filecache.New("./test").(*filecache.CacheImpl)
		as.Nil(c.Set("63125", "63125", time.Second))
	})

	t.Run("large count get set", func(t *testing.T) {
		as.Nil(os.Remove("./test"))
		c = filecache.New("./test").(*filecache.CacheImpl)

		for i := 0; i <= 63124; i++ {
			j := strconv.Itoa(i)
			as.Nil(c.Set(j, j, time.Minute), i)
		}

		for i := 0; i <= 61324; i++ {
			j := strconv.Itoa(i)
			v, err := c.Get(j)
			as.Nil(err)
			as.Equal(v, j)
		}
	})
}
