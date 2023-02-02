package lru

import "container/list"

type Key string

type Cache struct {
	MaxEntries int                              // 最大容量
	ll         *list.List                       //双向链表
	cache      map[Key]*list.Element            //通过key找value，再将ll中的节点移到首节点
	OnEvicted  func(key Key, value interface{}) //某条记录被移除时的回调函数
}

type entry struct {
	key   Key
	value interface{}
}

func New(maxEntries int, onEvicted func(key Key, value interface{})) *Cache {
	return &Cache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[Key]*list.Element),
		OnEvicted:  onEvicted,
	}
}

// Get 查找某个值
func (c *Cache) Get(key Key) (interface{}, bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		et := ele.Value.(*entry)
		return et.value, ok
	}
	return nil, false
}

// Remove 删除最近最少使用的节点
func (c *Cache) Remove() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		et := ele.Value.(*entry)
		delete(c.cache, et.key)
		if c.OnEvicted != nil {
			c.OnEvicted(et.key, et.value)
		}
	}
}

// Add 添加某个节点
func (c *Cache) Add(key Key, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[Key]*list.Element)
		c.ll = list.New()
	}
	old, ok := c.cache[key]
	if ok {
		oldEt := old.Value.(*entry)
		oldEt.value = value
		c.ll.MoveToFront(old)
		return
	}
	newEt := c.ll.PushFront(&entry{key, value})
	c.cache[key] = newEt
	if c.MaxEntries > 0 && c.MaxEntries < c.ll.Len() {
		c.Remove()
	}
}
