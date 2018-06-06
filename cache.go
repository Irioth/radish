package radish

import (
	"sync"
	"sync/atomic"
	"time"
)

// Warning: Never do that in production,
// always use proven open source libraries for complex concurrent data structures and algorithms
// (for example https://github.com/patrickmn/go-cache or https://github.com/karlseguin/ccache),
// this is here exclusively for fun and educational purposes.

// simple cuncurrent epoch based cache with only ttl expiration

const buckets = 3 // minimum 3

type cache struct {
	bucket [buckets]sync.Map
	epoch  int32
}

type wrapped struct {
	value   interface{}
	created time.Time
	ttl     time.Duration
	removed bool
}

func (w wrapped) unwrap() (interface{}, bool) {
	if w.expired() {
		return nil, false
	}
	return w.value, true
}

func (w wrapped) expired() bool {
	if w.removed {
		return true
	}
	if w.ttl == NoExpiration {
		return false
	}
	return time.Since(w.created) >= w.ttl
}

func newcache() *cache {
	return &cache{}
}

// get life value from cache
func (c *cache) get(key string) (interface{}, bool) {
	cur := atomic.LoadInt32(&c.epoch)
	v, ok := c.bucket[cur].Load(key)
	if !ok {
		v, ok = c.bucket[last(cur)].Load(key)
		if !ok {
			return nil, false
		}
	}

	return v.(wrapped).unwrap()
}

// internal method put value with defined creation time
func (c *cache) setw(key string, value wrapped) {
	cur := atomic.LoadInt32(&c.epoch)
	c.bucket[cur].Store(key, value)
}

// put value into cache
func (c *cache) set(key string, value interface{}, ttl time.Duration) {
	c.setw(key, wrapped{value, time.Now(), ttl, false})
}

// put tomb value into cache
func (c *cache) remove(key string) {
	c.setw(key, wrapped{removed: true})
}

// returns collection of life keys in cache
func (c *cache) keys() []string {
	result := map[string]bool{}
	cur := atomic.LoadInt32(&c.epoch)
	// collect keys with expiration status from last epoch firstly
	c.bucket[last(cur)].Range(func(key, value interface{}) bool {
		result[key.(string)] = value.(wrapped).expired()
		return true
	})
	// then collect them from current epoch
	c.bucket[cur].Range(func(key, value interface{}) bool {
		result[key.(string)] = value.(wrapped).expired()
		return true
	})

	// take only live ones
	list := []string{}
	for key, expired := range result {
		if !expired {
			list = append(list, key)
		}
	}

	return list
}

// we expect that calls to nextEpoch occurs rare, so no concurrent writes to last epoch exist at this moment
func (c *cache) nextEpoch() {
	// steps order important!
	cur := atomic.LoadInt32(&c.epoch)
	// copy live values from last epoch
	c.bucket[last(cur)].Range(func(key, value interface{}) bool {
		if !value.(wrapped).expired() {
			c.bucket[cur].LoadOrStore(key, value)
		}
		return true
	})
	// cleanup next epoch
	c.bucket[next(cur)] = sync.Map{}
	atomic.StoreInt32(&c.epoch, next(cur))
}

// return next epoch
func next(epoch int32) int32 {
	return (epoch + 1) % buckets
}

// return last epoch
func last(epoch int32) int32 {
	return (epoch - 1 + buckets) % buckets
}
