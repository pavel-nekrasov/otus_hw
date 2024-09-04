package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type internalData struct {
	Key   Key
	Value interface{}
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	data := internalData{Key: key, Value: value}

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = data
		return true
	}

	c.storeItem(key, c.queue.PushFront(data))
	if c.queue.Len() > c.capacity {
		c.removeItem(c.queue.Back())
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		data, ok := item.Value.(internalData)
		return data.Value, ok
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) storeItem(key Key, item *ListItem) {
	c.items[key] = item
}

func (c *lruCache) removeItem(item *ListItem) {
	data, ok := item.Value.(internalData)

	if ok {
		delete(c.items, data.Key)
		c.queue.Remove(item)
	}
}
