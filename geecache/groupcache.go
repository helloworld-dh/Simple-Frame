package geecache

import (
	"fmt"
	"gin/geecache/lru"
	"gin/geecache/singleflight"
	"log"
	"sync"
)

// Getter 加载数据
type Getter interface {
	Get(Key lru.Key) ([]byte, error)
}

type GetterFunc func(key lru.Key) ([]byte, error)

func (g GetterFunc) Get(key lru.Key) ([]byte, error) {
	return g(key)
}

// Group 每个缓存的命名空间
type Group struct {
	name      string              //缓存名称
	getter    Getter              //缓存未命中，调用回调函数获取缓存
	mainCache cache               //缓存
	peers     PeerPicker          // 将实现了PeerPicker接口的HTTPPool加入到Group中
	loader    *singleflight.Group //对于相同的key，只请求数据一次
}

var (
	groups = make(map[string]*Group)
	mu     sync.RWMutex
)

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// Get 获取数据
func (g *Group) Get(key lru.Key) (ByteView, error) {
	if key == "" {
		return ByteView{nil}, fmt.Errorf("key is required")
	}
	bv, ok := g.mainCache.get(key)
	if !ok {
		return g.load(key)
	}
	return bv, nil
}

// load 加载数据
func (g *Group) load(key lru.Key) (ByteView, error) {

	viewi, err := g.loader.Do(string(key), func() (interface{}, error) {
		// 存在分布式节点
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(string(key)); ok {
				value, err := g.getFromPeer(peer, string(key))
				if err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}

		}
		return g.getLocally(key)
	})
	if err != nil {
		return ByteView{}, err
	}
	return viewi.(ByteView), err
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{bytes}, nil
}

// getLocally 从本地加载数据
func (g *Group) getLocally(key lru.Key) (ByteView, error) {
	localCache, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	localVal := ByteView{b: localCache}
	g.populateCache(key, localVal)
	return localVal, nil
}

// populateCache 获得数据，并加入缓存
func (g *Group) populateCache(key lru.Key, value ByteView) {
	g.mainCache.add(key, value)
}

func NewGroup(name string, getter Getter, cacheBytes int) *Group {
	if getter == nil {
		panic("Getter nothing")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := groups[name]; ok {
		panic("duplicate registration of group" + name)
	}
	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g, ok := groups[name]
	if !ok {
		return nil
	}
	return g
}
