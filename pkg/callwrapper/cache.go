package callwrapper

import (
	"time"

	"github.com/karlseguin/ccache"
)

// CacheConfig : using karlseguin/ccache lib as basis. Accepting cache size (max units) and cache TTL (lifetime)
type CacheConfig struct {
	CacheTTLSec int
	CacheSize   int
}

type cache struct {
	cc      *ccache.Cache
	timeout time.Duration
}

func newCache(cacheSize, cacheTimeoutSec int) *cache {
	return &cache{
		cc:      ccache.New(ccache.Configure().MaxSize(int64(cacheSize))),
		timeout: time.Duration(cacheTimeoutSec) * time.Second,
	}
}

func (c *cache) Set(key string, val interface{}) {
	c.cc.Set(key, val, c.timeout)
}

func (c *cache) Get(key string) (interface{}, bool) {
	item := c.cc.Get(key)
	if item == nil || item.Expired() {
		return nil, false
	}
	return item.Value(), true
}
