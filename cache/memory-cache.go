package cache

import (
	"github.com/evansnowly/cache-system/util"
	"log"
	"sync"
	"time"
)

// MemoryCache is a struct that represents an in-memory cache with a fixed maximum size.
type MemoryCache struct {
	maxMemorySize     int64             // Maximum allowed memory size for the cache
	maxMemoryString   string            // Human-readable string representation of the maximum memory size
	currentMemorySize int64             // Current memory usage of the cache
	cacheMap          map[string]*Value // The underlying map to store cache values
	locker            sync.RWMutex      // Read-write mutex for safe concurrent access
	clearTime         time.Duration     // Interval for clearing expired cache entries
}

// Value represents a cached value with an expiration time and size.
type Value struct {
	value  any       // The actual value stored in the cache
	expire time.Time // Expiration time of the cached value
	size   int64     // Size of the cached value in bytes
}

// NewMemoryCache initializes a new MemoryCache with default settings.
func NewMemoryCache() *MemoryCache {
	m := &MemoryCache{
		cacheMap:  make(map[string]*Value, 100), // Initialize the cache map with a capacity of 100
		clearTime: time.Second * 1,              // Set the default clear interval to 1 second
	}
	go m.clearExpireCache(m.clearTime) // Start a goroutine to clear expired cache entries periodically
	return m
}

// deleteKey removes a key from the cache and adjusts the current memory size.
func (m *MemoryCache) deleteKey(key string) bool {
	value, ok := m.cacheMap[key]
	if value != nil && ok {
		delete(m.cacheMap, key)
		m.currentMemorySize -= value.size
		return true
	}
	return false
}

// addKey adds a new key-value pair to the cache and updates the current memory size.
func (m *MemoryCache) addKey(key string, value *Value) {
	m.cacheMap[key] = value
	m.currentMemorySize += value.size
}

// SetMaxMemory sets the maximum allowed memory size for the cache.
func (m *MemoryCache) SetMaxMemory(size string) bool {
	var err error
	m.maxMemorySize, m.maxMemoryString, err = util.ParseSize(size)
	if err != nil {
		return false
	}
	return true
}

// Set adds a key-value pair to the cache with an expiration duration.
func (m *MemoryCache) Set(key string, value any, expire time.Duration) bool {
	m.locker.Lock()
	defer m.locker.Unlock()
	v := &Value{
		value:  value,
		expire: time.Now().Add(expire),
		size:   int64(util.SizeOfVariable(value)),
	}
	if v.size+m.currentMemorySize > m.maxMemorySize {
		log.Printf("超出内存限制, 当前对象使用内存 %d bytes, 最大内存 %d bytes", v.size, m.maxMemorySize)
		return false
	}
	m.deleteKey(key)
	m.addKey(key, v)
	return true
}

// Get retrieves a value from the cache by its key.
func (m *MemoryCache) Get(key string) (any, bool) {
	m.locker.Lock()
	defer m.locker.Unlock()
	value, ok := m.cacheMap[key]
	if !ok {
		return nil, false
	}
	if value.expire.Before(time.Now()) {
		m.deleteKey(key)
		return nil, false
	}
	return value, true
}

// Delete removes a key-value pair from the cache.
func (m *MemoryCache) Delete(key string) bool {
	m.locker.Lock()
	defer m.locker.Unlock()
	return m.deleteKey(key)
}

// Exists checks if a key is present in the cache.
func (m *MemoryCache) Exists(key string) bool {
	m.locker.Lock()
	defer m.locker.Unlock()
	_, ok := m.cacheMap[key]
	return ok
}

// Flush clears all entries in the cache.
func (m *MemoryCache) Flush() bool {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.cacheMap = make(map[string]*Value, 100)
	m.currentMemorySize = 0
	return true
}

// KeyNum returns the number of keys currently in the cache.
func (m *MemoryCache) KeyNum() int64 {
	m.locker.Lock()
	defer m.locker.Unlock()
	return int64(len(m.cacheMap))
}

// clearExpireCache periodically removes expired entries from the cache.
func (m *MemoryCache) clearExpireCache(dr time.Duration) {
	ticker := time.NewTicker(dr)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for key, value := range m.cacheMap {
				if value.expire.Before(time.Now()) {
					m.locker.Lock()
					m.deleteKey(key)
					m.locker.Unlock()
				}
			}

		}
	}
}
