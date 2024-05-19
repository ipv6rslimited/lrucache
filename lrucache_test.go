/*
**
** lrucache_test
** Tests for lrucache
**
** Distributed under the COOL License.
**
** Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
** All Rights Reserved
**
*/

package lrucache

import (
  "fmt"
  "sync"
  "testing"
)

func TestLRUCache(t *testing.T) {
  fmt.Println("TestLRUCacheFrequentAccess")
  TestLRUCacheFrequentAccess(t)

  fmt.Println("TestLRUCacheRapidInsertionEviction")
  TestLRUCacheRapidInsertionEviction(t)

  fmt.Println("TestLRUCacheHitMissRatio")
  TestLRUCacheHitMissRatio(t)

  fmt.Println("TestLRUCacheLargeDataSet")
  TestLRUCacheLargeDataSet(t)

  fmt.Println("TestLRUCacheBoundaryConditions")
  TestLRUCacheBoundaryConditions(t)

  fmt.Println("TestLRUCacheHighConcurrencyStress")
  TestLRUCacheHighConcurrencyStress(t)

  fmt.Println("TestLRUCacheExtremeEdgeCases")
  TestLRUCacheExtremeEdgeCases(t)
}

func TestLRUCacheHighConcurrencyStress(t *testing.T) {
  cache := NewLRUCache(100)

  var wg sync.WaitGroup
  const numOperations = 10000
  const numGoroutines = 100

  writer := func(key string, value string) {
    defer wg.Done()
    for i := 0; i < numOperations; i++ {
      cache.Put(key, value)
    }
  }

  reader := func(key string) {
    defer wg.Done()
    for i := 0; i < numOperations; i++ {
      cache.Get(key)
    }
  }

  keys := make([]string, numGoroutines)
  for i := 0; i < numGoroutines; i++ {
    keys[i] = fmt.Sprintf("key%d", i)
  }

  for _, key := range keys {
    wg.Add(1)
    go writer(key, "value"+key)
  }

  for _, key := range keys {
    wg.Add(1)
    go reader(key)
  }

  wg.Wait()

  if len(cache.cache) > cache.capacity {
    t.Errorf("Cache exceeded capacity")
  }
}

func TestLRUCacheExtremeEdgeCases(t *testing.T) {
  cache := NewLRUCache(1000000)
  for i := 0; i < 1000000; i++ {
    key := fmt.Sprintf("key%d", i)
    value := fmt.Sprintf("value%d", i)
    cache.Put(key, value)
  }
  if len(cache.cache) != 1000000 {
    t.Errorf("Expected cache to have 1000000 elements, but got %d", len(cache.cache))
  }

  cache = NewLRUCache(1)
  largeValue := make([]byte, 1024*1024)
  cache.Put("A", string(largeValue))
  if len(cache.cache) != 1 {
    t.Errorf("Expected cache to have 1 element, but got %d", len(cache.cache))
  }

  cache.Put("B", string(largeValue))
  if len(cache.cache) != 1 {
    t.Errorf("Expected cache to have 1 element, but got %d", len(cache.cache))
  }
  if _, ok := cache.Get("A"); ok {
    t.Errorf("Expected key A to be evicted")
  }
  if _, ok := cache.Get("B"); !ok {
    t.Errorf("Expected key B to be in cache")
  }
}

func TestLRUCacheFrequentAccess(t *testing.T) {
  cache := NewLRUCache(3)

  cache.Put("A", "valueA")
  cache.Put("B", "valueB")
  cache.Put("C", "valueC")

  cache.Get("A")
  cache.Get("A")
  cache.Get("A")

  cache.Put("D", "valueD")

  if _, ok := cache.Get("B"); ok {
    t.Errorf("Expected key B to be evicted")
  }
  if _, ok := cache.Get("C"); !ok {
    t.Errorf("Expected key C to remain in cache")
  }
  if _, ok := cache.Get("A"); !ok {
    t.Errorf("Expected key A to remain in cache")
  }
}

func TestLRUCacheRapidInsertionEviction(t *testing.T) {
  cache := NewLRUCache(3)

  for i := 0; i < 100; i++ {
    key := fmt.Sprintf("key%d", i)
    value := fmt.Sprintf("value%d", i)
    cache.Put(key, value)

    if len(cache.cache) > cache.capacity {
      t.Errorf("Cache exceeded capacity at iteration %d", i)
    }
  }

  for i := 0; i < 97; i++ {
    key := fmt.Sprintf("key%d", i)
    if _, ok := cache.Get(key); ok {
      t.Errorf("Expected key %s to be evicted", key)
    }
  }
}

func TestLRUCacheHitMissRatio(t *testing.T) {
  cache := NewLRUCache(3)

  cache.Put("A", "valueA")
  cache.Put("B", "valueB")
  cache.Put("C", "valueC")

  cache.Get("A")
  cache.Get("B")
  cache.Get("D")

  cache.Put("E", "valueE")

  cache.Get("A")
  cache.Get("C")
  cache.Get("E")

  hits, misses := cache.GetHitMissCount()

  expectedHits := 4
  expectedMisses := 2

  if hits != expectedHits {
    t.Errorf("Expected %d hits, got %d", expectedHits, hits)
  }
  if misses != expectedMisses {
    t.Errorf("Expected %d misses, got %d", expectedMisses, misses)
  }
}

func TestLRUCacheLargeDataSet(t *testing.T) {
  cache := NewLRUCache(100)

  for i := 0; i < 1000; i++ {
    key := fmt.Sprintf("key%d", i)
    value := fmt.Sprintf("value%d", i)
    cache.Put(key, value)

    if len(cache.cache) > cache.capacity {
      t.Errorf("Cache exceeded capacity at iteration %d", i)
    }
  }

  for i := 900; i < 1000; i++ {
    key := fmt.Sprintf("key%d", i)
    if _, ok := cache.Get(key); !ok {
      t.Errorf("Expected key %s to be in cache", key)
    }
  }

  for i := 0; i < 900; i++ {
    key := fmt.Sprintf("key%d", i)
    if _, ok := cache.Get(key); ok {
      t.Errorf("Expected key %s to be evicted", key)
    }
  }
}

func TestLRUCacheBoundaryConditions(t *testing.T) {
  cache := NewLRUCache(0)
  cache.Put("A", "valueA")
  if len(cache.cache) != 0 {
    t.Errorf("Expected cache to have 0 elements, but got %d", len(cache.cache))
  }
  if _, ok := cache.Get("A"); ok {
    t.Errorf("Expected not to find any keys in a zero capacity cache")
  }

  cache = NewLRUCache(1)
  cache.Put("A", "valueA")
  if _, ok := cache.Get("A"); !ok {
    t.Errorf("Expected to find key A")
  }

  cache.Put("B", "valueB")
  if _, ok := cache.Get("A"); ok {
    t.Errorf("Expected key A to be evicted")
  }
  if _, ok := cache.Get("B"); !ok {
    t.Errorf("Expected to find key B")
  }
}
