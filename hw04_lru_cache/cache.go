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
	mu       sync.RWMutex
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

func newCacheItem(key Key, value interface{}) *cacheItem {
	return &cacheItem{
		key:   key,
		value: value,
	}
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	node, isInCache := l.items[key]
	var cachedValue interface{} = nil

	if isInCache {
		l.queue.MoveToFront(node)
		item := l.cachedItemFromNode(node)
		cachedValue = item.value
	}

	return cachedValue, isInCache
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	node, isInCache := l.items[key]
	if isInCache {
		item := l.cachedItemFromNode(node)
		item.value = value
		l.queue.MoveToFront(node)
	} else {
		if l.quequeIsFull() {
			l.popLastFromQueue()
		}
		item := newCacheItem(key, value)
		node = l.queue.PushFront(item)
		l.items[key] = node
	}
	return isInCache
}

func (l *lruCache) quequeIsFull() bool {
	return l.queue.Len() == l.capacity
}

func (l *lruCache) popLastFromQueue() {
	leastUsedNode := l.queue.Back()
	leastUsedItem := l.cachedItemFromNode(leastUsedNode)
	delete(l.items, leastUsedItem.key)
	l.queue.Remove(leastUsedNode)
}

func (l *lruCache) cachedItemFromNode(node *ListItem) *cacheItem {
	return node.Value.(*cacheItem)
}
