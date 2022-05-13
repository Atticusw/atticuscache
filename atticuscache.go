package atticuscache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 声明一个 与 Get 相同的 func
type GetterFunc func(key string) ([]byte, error)

// Get 接口型函数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{name: name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get value for a key from cache
func (g Group) Get(key string) (ByteView, error) {
	// 判断 key
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required...")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[Cache] hit...")
		return v, nil
	}
	return g.load(key)
}

func (g Group) load(key string) (ByteView, error) {
	return g.getLocally(key)

}

func (g Group) getLocally(key string) (ByteView, error) {
	// 调用用户回调函数，获取源数据
	bytes, err := g.getter.Get(key)

	if err != nil {
		return ByteView{}, err
	}
	// 将数据copy
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
