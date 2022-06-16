package hw04lrucache

import "sync"

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
	mutex    sync.Mutex
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex.Lock()

	item, exists := cache.items[key]

	if !exists {
		cache.queue.PushFront(cacheItem{key: key, value: value})
		cache.items[key] = cache.queue.Front()

		deleteUnusedIfOverflow(cache)
	} else {
		updatingCacheItem := item.Value.(cacheItem)
		updatingCacheItem.value = value

		item.Value = updatingCacheItem

		cache.queue.MoveToFront(item)
	}

	cache.mutex.Unlock()

	return exists
}

func deleteUnusedIfOverflow(cache *lruCache) {
	if len(cache.items) <= cache.capacity {
		return
	}

	lastItem := cache.queue.Back()
	cache.queue.Remove(lastItem)

	delete(cache.items, lastItem.Value.(cacheItem).key)
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item, exists := cache.items[key]
	var value interface{}

	if exists {
		cache.queue.MoveToFront(item)
		value = item.Value.(cacheItem).value
	}

	return value, exists
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}

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
