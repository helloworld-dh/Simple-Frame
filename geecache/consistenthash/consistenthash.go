package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // 用户自定义hash函数
	replicas int            // 节点数
	keys     []int          // hash环
	hashMap  map[int]string // hash环对应的是真节点
}

func NewMap(replicas int, hash Hash) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加真实节点，每一个真实节点创建replicas个虚拟节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
}

// Get 根据key获得真实节点索引值
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// 在hash环中查找key对应的虚拟idx
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	virNode := m.keys[idx%len(m.keys)]
	if realNode, ok := m.hashMap[virNode]; ok {
		return realNode
	}
	return ""
}

// 删除节点
