package mu

import (
	"sync"
	"time"
)

const defaultCap = 100
const defaultTTL = 5 * time.Minute

type Cache[K comparable, V any] struct {
	lock      sync.Mutex
	cache     map[K]*entry[K, V]
	ttl       time.Duration
	cap       int
	keepAlive bool
}

type entry[K comparable, V any] struct {
	value      V
	expiration time.Time
}

func NewExpLRU[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		cache: make(map[K]*entry[K, V], defaultCap),
		ttl:   defaultTTL,
		cap:   defaultCap,
	}
}

func (c *Cache[K, V]) Cap(cap int) *Cache[K, V] {
	c.cap = cap
	return c
}

func (c *Cache[K, V]) TTL(ttl time.Duration) *Cache[K, V] {
	c.ttl = ttl
	return c
}

func (c *Cache[K, V]) KeepAlive(ka bool) *Cache[K, V] {
	c.keepAlive = ka
	return c
}

func (c *Cache[K, V]) ensureCapacity() {
	if len(c.cache) < c.cap {
		return
	}

	var keysToPurge []K
	now := time.Now()
	for k, v := range c.cache {
		if v.expiration.Before(now) {
			keysToPurge = append(keysToPurge, k)
		}
	}
	for _, key := range keysToPurge {
		delete(c.cache, key)
	}

	if len(c.cache) < c.cap {
		return
	}

	numKeysToDelete := c.cap / 2
	keysToPurge = []K{}

	for k := range c.cache {
		keysToPurge = append(keysToPurge, k)
		if len(keysToPurge) >= numKeysToDelete {
			break
		}
	}

	for _, key := range keysToPurge {
		delete(c.cache, key)
	}
}

func (c *Cache[K, V]) Add(key K, val V) {
	c.lock.Lock()
	defer c.lock.Unlock()

	ent := entry[K, V]{
		value:      val,
		expiration: time.Now().Add(c.ttl),
	}
	c.ensureCapacity()
	c.cache[key] = &ent
}

func (c *Cache[K, V]) Get(key K) (val V, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	ent, ok := c.cache[key]
	if !ok {
		return val, false
	}

	if ent.expiration.Before(time.Now()) {
		delete(c.cache, key)
		return val, false
	}
	if c.keepAlive {
		ent.expiration = time.Now().Add(c.ttl)
	}

	return ent.value, true
}
