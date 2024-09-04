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
