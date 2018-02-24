package lrucache

import (
	"time"
	"testing"
	"reflect"
)

var (
	cache          *LRUCache
	cacheNoTimeout *LRUCache
)

func init() {
	cache, _ = InitLRUCache(3, time.Millisecond*1)
	cacheNoTimeout, _ = InitLRUCache(3, 0)
}

func TestClear(t *testing.T) {
	cache.Set("key1", "value1")
	cache.Clear()
	if cache.Count() != 0 {
		t.Errorf("Expected cache empty")
	}
}

func TestSetAndGet(t *testing.T) {
	cache.Clear()
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	if !reflect.DeepEqual(keysToStrArr(cache.Keys()), []string{"key3", "key2", "key1"}) {
		t.Errorf("Expected right cache sequence")
	}

	if value, found := cache.Get("key1"); !found || value != "value1" {
		t.Errorf("Expected get value")
		if !reflect.DeepEqual(keysToStrArr(cache.Keys()), []string{"key1", "key3", "key2"}) {
			t.Errorf("Expected right cache sequence")
		}
	}

	if _, found := cache.Get("key0"); found {
		t.Errorf("Expected cannot get value")
	}

	time.Sleep(time.Millisecond * 1)
	if _, found := cache.Get("key1"); found {
		t.Errorf("Expected cannot get value")
	}

	cacheNoTimeout.Set("key1", "value1")
	time.Sleep(time.Millisecond * 1)
	if _, found := cacheNoTimeout.Get("key1"); !found {
		t.Errorf("Expected get value")
	}
}

func TestPeek(t *testing.T) {
	cache.Clear()
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	if value, found := cache.Peek("key1"); !found || value != "value1" {
		t.Errorf("Expected peek value")
		if !reflect.DeepEqual(keysToStrArr(cache.Keys()), []string{"key1", "key3", "key2"}) {
			t.Errorf("Expected right cache sequence")
		}

	}
}

func TestRemove(t *testing.T) {
	cache.Clear()
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	cache.Remove("key3")
	if !reflect.DeepEqual(keysToStrArr(cache.Keys()), []string{"key2", "key1"}) {
		t.Errorf("Expected right cache sequence")
	}
}

func keysToStrArr(keys []interface{}) []string {
	var ret []string
	for _, value := range keys {
		switch v := value.(type) {
		case string:
			ret = append(ret, v)
		default:
			ret = append(ret, "")
		}

	}
	return ret
}

