package cache

import (
	"container/list"
	"sync"
)

type Cached[K comparable, V any] struct {
	key  K
	data V
}

type Cache[K comparable, V any] struct {
	maxEntries int
	cache      map[K]*list.Element
	list       *list.List
	lock       sync.RWMutex
}

func New[K comparable, V any](maxEntries int) *Cache[K, V] {
	return &Cache[K, V]{
		maxEntries: maxEntries,
		list:       list.New(),
		cache:      make(map[K]*list.Element),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if val, ok := c.cache[key]; ok {
		c.list.MoveToFront(val)
		val.Value.(*Cached[K, V]).data = value
		return
	}
	val := c.list.PushFront(&Cached[K, V]{
		key:  key,
		data: value,
	})
	c.cache[key] = val
	if c.list.Len() > c.maxEntries {
		c.removeOldest()
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if val, ok := c.cache[key]; ok {
		c.list.MoveToFront(val)
		return val.Value.(*Cached[K, V]).data, true
	}
	var v V
	return v, false
}

func (c *Cache[K, V]) Del(key K) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if val, ok := c.cache[key]; ok {
		c.remove(val)
	}
}

func (c *Cache[K, V]) removeOldest() {
	el := c.list.Back()
	if el != nil {
		c.remove(el)
	}
}

func (c *Cache[K, V]) remove(el *list.Element) {
	c.list.Remove(el)
	val := el.Value.(*Cached[K, V])
	delete(c.cache, val.key)
}
