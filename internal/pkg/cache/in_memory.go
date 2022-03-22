package cache

import (
	"sync"
	"time"
)

// Clock  - interface for easier mock time.Now in tests
type Clock interface {
	Now() time.Time
}

// InMemoryCache - implementation of Cache interface
// local in-memory storage, replacement for Redis in tests
// Mutex is used to protect map (sync.Map can be used too)
type InMemoryCache struct {
	dataMap map[int]inMemoryValue
	lock    *sync.Mutex
	clock   Clock
}

// inMemoryValue - internal struct to check expiration on values in cache
type inMemoryValue struct {
	SetTime    int64
	Expiration int64
}

// InitInMemoryCache - create new instance of InMemoryCache
// clock - instance of Clock to get time.Now() (and mocks in tests)
func InitInMemoryCache(clock Clock) *InMemoryCache {
	return &InMemoryCache{
		dataMap: make(map[int]inMemoryValue, 0),
		lock:    &sync.Mutex{},
		clock:   clock,
	}
}

// Add - add rand value with expiration (in seconds) to cache
func (c *InMemoryCache) Add(key int, expiration int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.dataMap[key] = inMemoryValue{
		SetTime:    c.clock.Now().Unix(),
		Expiration: expiration,
	}
	return nil
}

// Get - check existence of int key in cache
func (c *InMemoryCache) Get(key int) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.dataMap[key]
	if ok && c.clock.Now().Unix()-value.SetTime > value.Expiration {
		return false, nil
	}
	return ok, nil
}

// Delete - delete key from cache
func (c *InMemoryCache) Delete(key int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.dataMap, key)
}
