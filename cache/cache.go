package cache

import "time"

// Cache is an interface that defines the methods for an in-memory cache system.
type Cache interface {
	// SetMaxMemory sets the maximum allowed memory size for the cache.
	// size: a string representing the maximum memory size, e.g., "10MB".
	// Returns true if the size was set successfully, false otherwise.
	SetMaxMemory(size string) bool

	// Set adds a key-value pair to the cache with an expiration duration.
	// key: the key under which the value will be stored.
	// value: the value to store in the cache.
	// expire: the duration after which the key-value pair will expire.
	// Returns true if the value was set successfully, false otherwise (e.g., if it exceeds max memory).
	Set(key string, value any, expire time.Duration) bool

	// Get retrieves a value from the cache by its key.
	// key: the key to look up in the cache.
	// Returns the value associated with the key and a boolean indicating whether the key was found.
	Get(key string) (any, bool)

	// Delete removes a key-value pair from the cache.
	// key: the key to remove from the cache.
	// Returns true if the key was deleted successfully, false otherwise.
	Delete(key string) bool

	// Exists checks if a key is present in the cache.
	// key: the key to check for existence in the cache.
	// Returns true if the key exists, false otherwise.
	Exists(key string) bool

	// Flush clears all entries in the cache.
	// Returns true if the cache was cleared successfully, false otherwise.
	Flush() bool

	// KeyNum returns the number of keys currently in the cache.
	// Returns the number of keys as an int64.
	KeyNum() int64
}
