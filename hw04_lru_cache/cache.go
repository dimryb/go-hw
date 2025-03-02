package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

// Для совместного хранения ключа в элементе очереди.
type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func getKeyFromValue(value interface{}) Key {
	if item, ok := value.(cacheItem); ok {
		return item.key
	}
	return ""
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item, ok := c.items[key]
	cacheValue := cacheItem{key, value}

	if ok {
		item.Value = cacheValue
		c.queue.MoveToFront(item)
		return true
	}

	if c.queue.Len() >= c.capacity {
		backItem := c.queue.Back()
		if backItem != nil {
			delete(c.items, getKeyFromValue(backItem.Value))
			c.queue.Remove(backItem)
		}
	}

	c.items[key] = c.queue.PushFront(cacheValue)
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)
	if i, ok := item.Value.(cacheItem); ok {
		return i.value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {}
