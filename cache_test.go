package filecache_test

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Chyroc/filecache"

	chyrocFileCache "github.com/Chyroc/filecache"
	dannyBenFileCache "github.com/DannyBen/filecache"
	fabiorphpFileCache "github.com/fabiorphp/cachego"
	gadelkareemFileCache "github.com/gadelkareem/cachita"
	gookitFileCache "github.com/gookit/cache"
	huntsmanFileCache "github.com/huntsman-li/go-cache"
	miguelmotaFileCache "github.com/miguelmota/go-filecache"
	"github.com/stretchr/testify/assert"
)

func Example_Cache() {
	cache := filecache.New("./cache.data")
	defer os.Remove("./cache.data")

	_, err := cache.Get("not-exist")
	fmt.Println(err)

	fmt.Println(cache.Set("k", "v", time.Minute))

	v, err := cache.Get("k")
	fmt.Println(v, err)

	ttl, err := cache.TTL("k")
	fmt.Println(int(math.Ceil(ttl.Seconds())), err)

	time.Sleep(time.Second)

	ttl, err = cache.TTL("k")
	fmt.Println(int(math.Ceil(ttl.Seconds())), err)

	fmt.Println(cache.Del("k"))

	_, err = cache.Get("k")
	fmt.Println(err)

	// output:
	// not found
	// <nil>
	// v <nil>
	// 60 <nil>
	// 59 <nil>
	// <nil>
	// not found
}

func TestNew(t *testing.T) {
	as := assert.New(t)
	defer os.Remove("./test")

	os.Remove("./test")
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

	t.Run("large count get set del", func(t *testing.T) {
		as.Nil(os.Remove("./test"))
		c = filecache.New("./test").(*filecache.CacheImpl)

		// set
		for i := 0; i <= 63124; i++ {
			j := strconv.Itoa(i)
			as.Nil(c.Set(j, j, time.Minute), i)
		}

		// get exist
		for i := 0; i <= 61324; i++ {
			j := strconv.Itoa(i)
			v, err := c.Get(j)
			as.Nil(err)
			as.Equal(v, j)
		}

		// del
		for i := 0; i <= 61324; i++ {
			j := strconv.Itoa(i)
			as.Nil(c.Del(j))
		}

		// get not-exist
		for i := 0; i <= 61324; i++ {
			j := strconv.Itoa(i)
			_, err := c.Get(j)
			as.Equal(filecache.NotFound, err)
		}

		// set again: should success
		// set
		for i := 0; i <= 63124; i++ {
			j := strconv.Itoa(i)
			as.Nil(c.Set(j, j, time.Minute), i)
		}
	})

	t.Run("range", func(t *testing.T) {
		as.Nil(os.Remove("./test"))
		c = filecache.New("./test").(*filecache.CacheImpl)

		for i := 0; i < 1000; i++ {
			j := strconv.Itoa(i)
			as.Nil(c.Set(j, j, time.Minute), i)
		}

		kvs, err := c.Range()
		as.Nil(err)
		for _, v := range kvs {
			as.Equal(v.Key, v.Val)
		}
		as.Len(kvs, 1000)
	})
}

func BenchmarkFileCache(b *testing.B) {
	as := assert.New(b)

	b.Run("chyroc", func(b *testing.B) {
		file := "./test-file-chyroc"
		defer os.RemoveAll(file)
		os.RemoveAll(file)

		b.ResetTimer()

		c := chyrocFileCache.New(file)

		for i := 0; i < b.N; i++ {
			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				as.Nil(c.Set(j, j, time.Second))
			}

			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				v, err := c.Get(strconv.Itoa(i))
				as.Nil(err)
				as.Equal(j, v)
			}
		}
	})

	b.Run("huntsman", func(b *testing.B) {
		for _, v := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"} {
			defer os.RemoveAll(v)
			os.RemoveAll(v)
		}

		b.ResetTimer()

		c := huntsmanFileCache.NewFileCacher()

		for i := 0; i < b.N; i++ {
			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				as.Nil(c.Put(j, j, 100))
			}

			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				vv := c.Get(strconv.Itoa(i))
				v, ok := vv.(string)
				as.True(ok)
				as.Equal(j, v)
			}
		}

	})

	b.Run("fabiorphp", func(b *testing.B) {
		file := "./test-file-fabiorphp"
		os.RemoveAll(file)
		as.Nil(os.Mkdir(file, 0755))
		defer os.RemoveAll(file)

		b.ResetTimer()

		c := fabiorphpFileCache.NewFile(file)

		for i := 0; i < b.N; i++ {
			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				as.Nil(c.Save(j, j, time.Minute))
			}

			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				v, err := c.Fetch(strconv.Itoa(i))
				as.Nil(err)
				as.Equal(j, v)
			}
		}
	})

	b.Run("dannyBen", func(b *testing.B) {
		file := "./test-file-dannyBen"
		os.RemoveAll(file)
		as.Nil(os.Mkdir(file, 0755))
		defer os.RemoveAll(file)

		b.ResetTimer()

		c := dannyBenFileCache.Handler{file, 600}

		for i := 0; i < b.N; i++ {
			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				as.Nil(c.Set(j, []byte(j)))
			}

			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				as.Equal(j, string(c.Get(strconv.Itoa(i))))
			}
		}
	})

	// too slow, skip
	b.Run("gookit", func(b *testing.B) {
		file := "./test-file-gookit"
		os.RemoveAll(file)
		defer os.RemoveAll(file)

		b.Skip()
		b.ResetTimer()

		c := gookitFileCache.NewFileCache(file)

		for i := 0; i < b.N; i++ {
			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				as.Nil(c.Set(j, j, time.Minute))
			}

			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				vv := c.Get(strconv.Itoa(i))
				v, ok := vv.(string)
				as.True(ok)
				as.Equal(j, v)
			}
		}
	})

	b.Run("gadelkareem", func(b *testing.B) {
		file := "./test-file-gadelkareem"
		os.RemoveAll(file)
		defer os.RemoveAll(file)

		b.ResetTimer()

		c, err := gadelkareemFileCache.NewFileCache(file, time.Minute, time.Minute)
		as.Nil(err)

		for i := 0; i < b.N; i++ {
			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				as.Nil(c.Put(j, j, time.Minute))
			}

			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				var v string
				as.Nil(c.Get(strconv.Itoa(i), &v))
				as.Equal(j, v)
			}
		}
	})

	b.Run("miguelmota", func(b *testing.B) {
		file := "./test-file-miguelmota"
		os.RemoveAll(file)
		defer os.RemoveAll(file)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				as.Nil(miguelmotaFileCache.Set(j, j, time.Minute))
			}

			for i := 0; i <= 1000; i++ {
				j := strconv.Itoa(i)
				var v string
				ok, err := miguelmotaFileCache.Get(strconv.Itoa(i), &v)
				as.True(ok)
				as.Nil(err)
				as.Equal(j, v)
			}
		}
	})

}
