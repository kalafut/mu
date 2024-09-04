package mu

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewExpLRU(t *testing.T) {
	cache := NewCache[string, int]().TTL(time.Minute)
	assert.NotNil(t, cache, "NewExpLRU should not return nil")
	assert.Equal(t, time.Minute, cache.ttl, "Expiration time should be 1 minute")
}

func TestExpLRU_Add(t *testing.T) {
	cache := NewCache[string, int]()
	cache.Add("key", 42)

	val, ok := cache.Get("key")
	assert.True(t, ok, "Should be able to retrieve added value")
	assert.Equal(t, 42, val, "Retrieved value should be 42")
}

func TestExpLRU_Get_Expired(t *testing.T) {
	cache := NewCache[string, int]().TTL(time.Millisecond)
	cache.Add("key", 42)

	time.Sleep(2 * time.Millisecond)

	_, ok := cache.Get("key")
	assert.False(t, ok, "Should not retrieve expired value")
}

func TestExpLRU_Get_NotFound(t *testing.T) {
	cache := NewCache[string, int]()

	_, ok := cache.Get("non-existent")
	assert.False(t, ok, "Should not retrieve non-existent value")
}

func TestExpLRU_Eviction(t *testing.T) {
	cache := NewCache[int, int]().Cap(10).TTL(100 * time.Millisecond)

	// Fill cache
	for i := 0; i < 10; i++ {
		cache.Add(i, i)
	}
	assert.Equal(t, 10, len(cache.cache))

	// Check that adding one more results in 0.5 + 1 cache size
	cache.Add(40, 40)
	assert.Equal(t, 6, len(cache.cache))

	// Fill cache
	for i := 0; i < 4; i++ {
		cache.Add(i+100, i)
	}
	assert.Equal(t, 10, len(cache.cache))

	// Wait, add 1, and check the all of the expired values are gone
	time.Sleep(200 * time.Millisecond)
	cache.Add(500, 500)
	assert.Equal(t, 1, len(cache.cache))
}

func TestExpLRU_Concurrent(t *testing.T) {
	cache := NewCache[int, int]().Cap(100)
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				cache.Add(i*100+j, j)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	count := 0
	for i := 0; i < 1000; i++ {
		if _, ok := cache.Get(i); ok {
			count++
		}
	}

	assert.Equal(t, 100, count, "Expected 100 items in cache")
}

func TestKeepAlive(t *testing.T) {
	// Test with KeepAlive set to false (default behavior)
	t.Run("KeepAlive false", func(t *testing.T) {
		cache := NewCache[string, int]().
			TTL(100 * time.Millisecond)

		cache.Add("key", 42)

		// First get should succeed
		val, ok := cache.Get("key")
		assert.True(t, ok, "Expected to get the value")
		assert.Equal(t, 42, val, "Expected to get 42")

		// Wait for half the TTL
		time.Sleep(50 * time.Millisecond)

		// Second get should still succeed
		val, ok = cache.Get("key")
		assert.True(t, ok, "Expected to get the value")
		assert.Equal(t, 42, val, "Expected to get 42")

		// Wait for the TTL to expire
		time.Sleep(60 * time.Millisecond)

		// Third get should fail
		_, ok = cache.Get("key")
		assert.False(t, ok, "Expected key to be expired")
	})

	// Test with KeepAlive set to true
	t.Run("KeepAlive true", func(t *testing.T) {
		cache := NewCache[string, int]().
			TTL(100 * time.Millisecond).
			KeepAlive(true)

		cache.Add("key", 42)

		// First get should succeed
		val, ok := cache.Get("key")
		assert.True(t, ok, "Expected to get the value")
		assert.Equal(t, 42, val, "Expected to get 42")

		// Wait for half the TTL
		time.Sleep(50 * time.Millisecond)

		// Second get should succeed and reset the expiration
		val, ok = cache.Get("key")
		assert.True(t, ok, "Expected to get the value")
		assert.Equal(t, 42, val, "Expected to get 42")

		// Wait for the original TTL to expire
		time.Sleep(60 * time.Millisecond)

		// Third get should still succeed because the expiration was reset
		val, ok = cache.Get("key")
		assert.True(t, ok, "Expected to get the value")
		assert.Equal(t, 42, val, "Expected to get 42")

		// Wait for the TTL to expire again
		time.Sleep(110 * time.Millisecond)

		// Fourth get should fail
		_, ok = cache.Get("key")
		assert.False(t, ok, "Expected key to be expired")
	})
}

func TestExpLRU_Remove(t *testing.T) {
	cache := NewCache[string, int]()

	// Add some items
	cache.Add("key1", 1)
	cache.Add("key2", 2)
	cache.Add("key3", 3)

	// Remove an existing key
	cache.Remove("key2")

	// Check that the removed key is no longer in the cache
	_, ok := cache.Get("key2")
	assert.False(t, ok, "Key 'key2' should have been removed")

	// Check that other keys are still present
	val1, ok := cache.Get("key1")
	assert.True(t, ok, "Key 'key1' should still be in the cache")
	assert.Equal(t, 1, val1, "Value for 'key1' should be 1")

	val3, ok := cache.Get("key3")
	assert.True(t, ok, "Key 'key3' should still be in the cache")
	assert.Equal(t, 3, val3, "Value for 'key3' should be 3")

	// Try to remove a non-existent key (should not cause any errors)
	cache.Remove("non-existent")
}

func TestExpLRU_Clear(t *testing.T) {
	cache := NewCache[string, int]()

	// Add some items
	cache.Add("key1", 1)
	cache.Add("key2", 2)
	cache.Add("key3", 3)

	// Clear the cache
	cache.Clear()

	// Check that all keys have been removed
	_, ok1 := cache.Get("key1")
	assert.False(t, ok1, "Key 'key1' should have been removed")

	_, ok2 := cache.Get("key2")
	assert.False(t, ok2, "Key 'key2' should have been removed")

	_, ok3 := cache.Get("key3")
	assert.False(t, ok3, "Key 'key3' should have been removed")

	// Add a new item after clearing to ensure the cache is still functional
	cache.Add("new-key", 42)
	val, ok := cache.Get("new-key")
	assert.True(t, ok, "Should be able to add and retrieve a new item after clearing")
	assert.Equal(t, 42, val, "New value should be retrievable after clearing")
}

func TestExpLRU_RemoveAndAdd(t *testing.T) {
	cache := NewCache[string, int]()

	// Add an item
	cache.Add("key", 1)

	// Remove the item
	cache.Remove("key")

	// Add the same key with a different value
	cache.Add("key", 2)

	// Check that the new value is present
	val, ok := cache.Get("key")
	assert.True(t, ok, "Key should be present after removing and re-adding")
	assert.Equal(t, 2, val, "New value should be retrieved after removing and re-adding")
}
