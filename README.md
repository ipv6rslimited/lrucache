# lrucache

## Overview

lrucache is a thread-safe, high-performance Least Recently Used (LRU) cache implementation in Go
with the simplest of APIs.

## Use case

We are improving our reverse proxy and needed a more efficient caching mechanism for DNS records.

## Features

- **Thread-Safe**: Uses `sync.RWMutex` to handle concurrent read/write operations safely.
- **High Performance**: Optimized for frequent access and rapid insertion/eviction.
- **Customizable Capacity**: Configure the cache size to suit your application's needs.
- **Hit/Miss Tracking**: Monitor cache performance with hit/miss statistics.

## Usage Examples

### Create a cache with a capacity of 100 items

```
import "github.com/ipv6rslimited/lrucache"

cache := lrucache.NewLRUCache(100) 
```

### Adding Items to the Cache

```
cache.Put("key1", "value1")
```

### Getting Items from the Cache

```
value, found := cache.Get("key1")
if found {
  fmt.Println("Found:", value)
} else {
  fmt.Println("Not found")
}
```

### Get Hit/Miss Count

```
import "fmt"

hits, misses := cache.GetHitMissCount()
fmt.Printf("Hits: %d, Misses: %d\n", hits, misses)
```

### Test Cases

- **Frequent Access**: Ensures that frequently accessed items are not evicted.
- **Rapid Insertion/Eviction**: Verifies cache behavior under rapid insertions and evictions.
- **Hit/Miss Ratio**: Validates the accuracy of hit/miss tracking.
- **Large Data Set**: Tests cache performance with a large number of items.
- **Boundary Conditions**: Checks edge cases, such as zero capacity and large value sizes.
- **High Concurrency Stress**: Tests cache performance under high concurrency.
- **Extreme Edge Cases**: Ensures the cache handles extreme scenarios gracefully.

## Running Tests

You can run a pretty sweet suite of tests by typing:

```
go test
```

## License

Distributed under the COOL License.

Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
All Rights Reserved

