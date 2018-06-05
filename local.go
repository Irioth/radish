package radish

import (
	"time"
)

type Radish struct {
	cache  *cache
	stop   chan bool
	ticker *time.Ticker
}

var _ Client = &Radish{}

func NewLocal() *Radish {
	r := &Radish{cache: newcache(), stop: make(chan bool, 1)}
	r.startCleanup()
	return r
}

func (r *Radish) startCleanup() {
	r.ticker = time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-r.ticker.C:
				r.cache.nextEpoch()
			case <-r.stop:
				break
			}
		}
	}()
}

func (r *Radish) Stop() {
	r.ticker.Stop()
	r.stop <- true
}

func (r *Radish) Get(key string) (interface{}, error) {
	if v, ok := r.cache.get(key); ok {
		return v, nil
	}
	return nil, NotFound
}

func (r *Radish) Set(key string, value interface{}, ttl time.Duration) error {
	r.cache.set(key, value, ttl)
	return nil
}

func (r *Radish) Remove(key string) error {
	r.cache.remove(key)
	return nil
}

func (r *Radish) Keys() ([]string, error) {
	return r.cache.keys(), nil
}

func (r *Radish) GetIndex(key string, index int) (interface{}, error) {
	v, ok := r.cache.get(key)
	if !ok {
		return nil, NotFound
	}
	x, ok := v.([]interface{})
	if !ok {
		return nil, NotList
	}
	if index < 0 || index >= len(x) {
		return nil, OutOfBound
	}
	return x[index], nil
}

func (r *Radish) GetDict(dictName string, key string) (interface{}, error) {
	v, ok := r.cache.get(dictName)
	if !ok {
		return nil, NotFound
	}
	x, ok := v.(map[string]interface{})
	if !ok {
		return nil, NotDict
	}
	return x[key], nil
}
