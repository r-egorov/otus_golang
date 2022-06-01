package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
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

func (l *lruCache) Set(key Key, value interface{}) bool {
	wasInCache := false
	if node, isInCache := l.items[key]; isInCache {
		item := node.Value.(*cacheItem)
		item.value = value
		l.queue.MoveToFront(node)
		wasInCache = true
	} else {
		item := newCacheItem(key, value)
		if l.queue.Len() == l.capacity {
			leastUsedNode := l.queue.Back()
			leastUsedItem := leastUsedNode.Value.(*cacheItem)
			delete(l.items, leastUsedItem.key)
			l.queue.Remove(leastUsedNode)
		}
		node = l.queue.PushFront(item)
		l.items[key] = node
	}
	return wasInCache
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	node, isInCache := l.items[key]
	var res interface{} = nil
	if isInCache {
		l.queue.MoveToFront(node)
		item := node.Value.(*cacheItem)
		res = item.value
	}
	return res, isInCache
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}
