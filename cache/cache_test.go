package cache

import (
	"testing"
	"time"
)

// TestNewMemoryCache tests the creation of a new MemoryCache.
func TestNewMemoryCache(t *testing.T) {
	cache := NewMemoryCache()
	if cache == nil {
		t.Error("Expected new memory cache to be created")
	}
	if len(cache.cacheMap) != 0 {
		t.Error("Expected new cache map to be empty")
	}
}

// TestSetMaxMemory tests setting the maximum memory size.
func TestSetMaxMemory(t *testing.T) {
	cache := NewMemoryCache()
	result := cache.SetMaxMemory("10MB")
	if !result {
		t.Error("Expected setting max memory to succeed")
	}
	if cache.maxMemoryString != "10MB" {
		t.Errorf("Expected max memory string to be '10MB', got '%s'", cache.maxMemoryString)
	}
}

// TestSetAndGet tests setting and getting a cache value.
func TestSetAndGet(t *testing.T) {
	cache := NewMemoryCache()
	cache.SetMaxMemory("10MB")
	cache.Set("key1", "value1", time.Minute)

	value, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key 'key1' in cache")
	}
	if value.(*Value).value != "value1" {
		t.Errorf("Expected value 'value1', got '%v'", value)
	}
}

// TestGetNonExistentKey tests getting a non-existent key.
func TestGetNonExistentKey(t *testing.T) {
	cache := NewMemoryCache()
	value, found := cache.Get("key2")
	if found {
		t.Error("Expected not to find key 'key2' in cache")
	}
	if value != nil {
		t.Errorf("Expected nil value, got '%v'", value)
	}
}

// TestDelete tests deleting a cache value.
func TestDelete(t *testing.T) {
	cache := NewMemoryCache()
	cache.SetMaxMemory("10MB")
	cache.Set("key1", "value1", time.Minute)

	deleted := cache.Delete("key1")
	if !deleted {
		t.Error("Expected key 'key1' to be deleted")
	}
	value, found := cache.Get("key1")
	if found {
		t.Error("Expected not to find key 'key1' in cache")
	}
	if value != nil {
		t.Errorf("Expected nil value, got '%v'", value)
	}
}

// TestExists tests checking for the existence of a key.
func TestExists(t *testing.T) {
	cache := NewMemoryCache()
	cache.SetMaxMemory("10MB")
	cache.Set("key1", "value1", time.Minute)

	exists := cache.Exists("key1")
	if !exists {
		t.Error("Expected key 'key1' to exist in cache")
	}
	cache.Delete("key1")
	exists = cache.Exists("key1")
	if exists {
		t.Error("Expected key 'key1' to not exist in cache")
	}
}

// TestFlush tests flushing the cache.
func TestFlush(t *testing.T) {
	cache := NewMemoryCache()
	cache.SetMaxMemory("10MB")
	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)

	cache.Flush()
	if cache.KeyNum() != 0 {
		t.Error("Expected cache to be empty after flush")
	}
}

// TestKeyNum tests getting the number of keys in the cache.
func TestKeyNum(t *testing.T) {
	cache := NewMemoryCache()
	cache.SetMaxMemory("10MB")
	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)

	if cache.KeyNum() != 2 {
		t.Errorf("Expected 2 keys in cache, got %d", cache.KeyNum())
	}
}

// TestClearExpireCache tests the automatic removal of expired cache entries.
func TestClearExpireCache(t *testing.T) {
	cache := NewMemoryCache()
	cache.SetMaxMemory("10MB")
	cache.Set("key1", "value1", time.Second)

	time.Sleep(2 * time.Second) // Wait for the cache entry to expire

	value, found := cache.Get("key1")
	if found {
		t.Error("Expected not to find expired key 'key1' in cache")
	}
	if value != nil {
		t.Errorf("Expected nil value for expired key, got '%v'", value)
	}
}
