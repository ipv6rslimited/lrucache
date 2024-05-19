/*
**
** lrucache
** A thread-safe, high-performance Least Recently Used (LRU) cache implementation.
**
** Distributed under the COOL License.
**
** Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
** All Rights Reserved
**
*/

package lrucache

import (
  "container/list"
  "sync"
)

type CacheItem struct {
  Key   string
  Value interface{}
}

type LRUCache struct {
  capacity int
  cache    map[string]*list.Element
  list     *list.List
  hits     int
  misses   int
  mutex    sync.RWMutex
}

func NewLRUCache(capacity int) *LRUCache {
  return &LRUCache{
    capacity: capacity,
    cache:    make(map[string]*list.Element),
    list:     list.New(),
  }
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
  c.mutex.RLock()
  element, found := c.cache[key]
  c.mutex.RUnlock()

  if found {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.list.MoveToFront(element)
    c.hits++
    return element.Value.(*CacheItem).Value, true
  }

  c.misses++
  return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) {
  c.mutex.Lock()
  defer c.mutex.Unlock()

  if c.capacity == 0 {
    return
  }

  if element, found := c.cache[key]; found {
    c.list.MoveToFront(element)
    element.Value.(*CacheItem).Value = value
    return
  }

  if c.list.Len() >= c.capacity {
    back := c.list.Back()
    if back != nil {
      evictedItem := back.Value.(*CacheItem)
      delete(c.cache, evictedItem.Key)
      c.list.Remove(back)
    }
  }

  item := &CacheItem{Key: key, Value: value}
  element := c.list.PushFront(item)
  c.cache[key] = element
}

func (c *LRUCache) GetHitMissCount() (int, int) {
  c.mutex.RLock()
  defer c.mutex.RUnlock()
  return c.hits, c.misses
}

