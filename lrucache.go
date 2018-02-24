package lrucache

import (
	"sync"
	"time"
	"container/list"
)

type LRUCache struct {
	mutex        sync.RWMutex
	items        map[interface{}]*list.Element
	linkedList   *list.List
	size         int
	cacheTimeout time.Duration
}

type entry struct {
	key        interface{}
	value      interface{}
	updateTime time.Time
}

func InitLRUCache(cacheSize int, timeout time.Duration) (cache *LRUCache, ok bool) {
	if cacheSize <= 0 {
		return nil, false
	}

	if timeout < 0 {
		return nil, false
	}

	cache = &LRUCache{
		items:        make(map[interface{}]*list.Element),
		size:         cacheSize,
		linkedList:   list.New(),
		cacheTimeout: timeout,
	}
	return cache, true
}

func (c *LRUCache) Set(key interface{}, value interface{}) {
	e := &entry{
		key:        key,
		value:      value,
		updateTime: time.Now(),
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.items[key]; ok {
		v.Value.(*entry).key = e.key
		v.Value.(*entry).value = e.value
		v.Value.(*entry).updateTime = e.updateTime
		c.linkedList.MoveToFront(v)
	} else {
		c.items[key] = c.linkedList.PushFront(e)
	}

	if c.linkedList.Len() > c.size {
		c.removeOldest()
		c.cleanTimeout()
	}
}

func (c *LRUCache) Get(key interface{}) (value interface{}, found bool) {
	c.mutex.RLock()
	if v, ok := c.items[key]; ok {
		if !c.expired(key) {
			c.mutex.RUnlock()
			c.mutex.Lock()
			c.linkedList.MoveToFront(v)
			c.mutex.Unlock()
			return v.Value.(*entry).value, true
		}
	}
	c.mutex.RUnlock()
	return nil, false
}

// Peek returns the value without updating
func (c *LRUCache) Peek(key interface{}) (value interface{}, found bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if v, ok := c.items[key]; ok {
		if !c.expired(key) {
			return v.Value.(*entry).value, true
		}
	}
	return nil, false
}

func (c *LRUCache) Remove(key interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.items[key]; ok {
		c.removeElement(v)
	}
}

// Keys returns keys order by latest used desc
func (c *LRUCache) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(c.items))

	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for e := c.linkedList.Front(); e != nil; e = e.Next() {
		keys[i] = e.Value.(*entry).key
		i++
	}
	return keys
}

func (c *LRUCache) Count() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.linkedList.Len()
}

func (c *LRUCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, v := range c.items {
		delete(c.items, v.Value.(*entry).key)
	}
	c.linkedList.Init()
}

func (c *LRUCache) cleanTimeout() {
	for k, v := range c.items {
		if c.expired(k) {
			c.removeElement(v)
		}
	}
}

func (c *LRUCache) expired(key interface{}) bool {
	if c.cacheTimeout == 0 {
		return false
	}

	if value, ok := c.items[key]; ok {
		if time.Now().Sub(value.Value.(*entry).updateTime) <= c.cacheTimeout {
			return false
		}
	}
	return true
}

func (c *LRUCache) removeElement(value *list.Element) {
	if _, ok := c.items[c.linkedList.Remove(value).(*entry).key]; ok {
		delete(c.items, c.linkedList.Remove(value).(*entry).key)
	}
}

func (c *LRUCache) removeOldest() {
	e := c.linkedList.Back()
	if e != nil {
		c.removeElement(e)
	}
}

