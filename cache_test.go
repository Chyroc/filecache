package filecache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Chyroc/filecache"
)

func TestNew(t *testing.T) {
	as := assert.New(t)

	c := filecache.NewDefault("").(*filecache.CacheImpl)
	as.Equal(1024, c.Mod)
	as.Equal(2097152, c.MaxSize)
}
