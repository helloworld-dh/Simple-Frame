package geecache

import (
	"gin/geecache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int
}

func (c *cache) add(key lru.Key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key lru.Key) (ByteView, bool) {
	if c.lru == nil {
		return ByteView{b: []byte{}}, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.lru.Get(key)
	if !ok {
		return ByteView{b: []byte{}}, false
	}
	return v.(ByteView), true
}
