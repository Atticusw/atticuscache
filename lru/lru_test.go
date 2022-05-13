package lru

import (
	"reflect"
	"testing"
)

type String string

// 实现了 Value 的方法
func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	cache := New(int64(0), nil)
	cache.Add("key1", String("1234"))
	if value, ok := cache.Get("key1"); !ok || string(value.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed...")
	}

	if _, ok := cache.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed...")
	}

}

func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1 + k2 + k3)
	cache := New(int64(cap), nil)
	cache.Add(k1, String(v1))
	cache.Add(k2, String(v2))
	cache.Add(k3, String(v3))
	if _, ok := cache.Get("key1"); ok || cache.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	cache := New(int64(10), callback)
	cache.Add("key1", String("123456"))
	cache.Add("k2", String("k2"))
	cache.Add("k3", String("k3"))
	cache.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}

func TestAdd(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key", String("1"))
	lru.Add("key", String("111"))

	if lru.nbytes != int64(len("key")+len("111")) {
		t.Fatal("expected 6 but got", lru.nbytes)
	}
}
