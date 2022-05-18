package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash Hash
	// 虚拟节点倍数
	replicas int
	// hash 环
	keys []int
	// 虚拟节点与真实节点的映射表，
	// 键是虚拟节点的hash 值
	// 值是真实节点的名称
	hashMap map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加真实节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 虚拟副本
		for i := 0; i < m.replicas; i++ {
			// 计算虚拟节点的hash值
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 存这个节点的hash值到 哈希环上
			m.keys = append(m.keys, hash)
			// 对应真实的节点
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get 选择节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
