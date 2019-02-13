package filecache_test

import (
	"github.com/Chyroc/filecache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	as := assert.New(t)

	c := filecache.New("./test").(*filecache.CacheImpl)

	t.Run("", func(t *testing.T) {
		as.Equal(5242880, c.CurrentSize)
	})

	t.Run("", func(t *testing.T) {
		as.Nil(c.Set("k1", "v1", time.Second))
	})

}
