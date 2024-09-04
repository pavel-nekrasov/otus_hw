package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu          sync.Mutex
	capacity    int
	queue       List
	items       map[Key]*ListItem
	itemsContra map[*ListItem]Key
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity:    capacity,
		queue:       NewList(),
		items:       make(map[Key]*ListItem, capacity),
		itemsContra: make(map[*ListItem]Key, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = value
		return true
	}

	c.storeItem(key, c.queue.PushFront(value))
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
		return item.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.itemsContra = make(map[*ListItem]Key, c.capacity)
}

func (c *lruCache) storeItem(key Key, item *ListItem) {
	c.items[key] = item
	c.itemsContra[item] = key
}

func (c *lruCache) removeItem(item *ListItem) {
	key := c.itemsContra[item]
	delete(c.items, key)
	delete(c.itemsContra, item)
	c.queue.Remove(item)
}
