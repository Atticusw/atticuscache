package lru

import "container/list"

// Cache 最近最少使用 lru，就是说最近被使用的表示后面也会被经常使用
// 当一个数据被访问，那么就将他移动到队尾表示。当内存不够的时候删除队头
// 使用双向链表存储key，使用map存储值
type Cache struct {
	// 最大的内存
	maxBytes int64
	// 当前已使用的内存
	nbytes int64
	// 实际的数据值存放在 双链表中
	ll *list.List
	// 数据存放在 map 中方便查找
	cache map[string]*list.Element
	// 某条记录被移除时的回调函数
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value 允许值是实现了 Value 接口的任意类型， len() 方法是返回值所占用的内存大小
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOdlest removes the oldest item
func (c *Cache) RemoveOdlest() {
	// 返回双链表最后的元素
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		// 当前的字节大小，等于 key的长度加上value的长度
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 如果存在移动到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 修改当前的内存值，
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		// 将新值替换给 老值
		kv.value = value
	} else {
		// 将新数据插入到队尾
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOdlest()
	}
}

// 返回缓存的长度
func (c *Cache) Len() int {
	return c.ll.Len()
}
