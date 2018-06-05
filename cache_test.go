package radish

import (
	"sort"
	"testing"
	"time"
)

func TestGetNonExisted(t *testing.T) {
	c := newcache()
	v, ok := c.get("key")
	if ok {
		t.Fatal("found unexisted key")
	}
	if v != nil {
		t.Fatal("found unexisted value")
	}
}

func TestSetGet(t *testing.T) {
	c := newcache()
	c.set("key", "ups", NoExpiration)
	v, ok := c.get("key")
	if !ok {
		t.Fatal("key must be found")
	}
	if v != "ups" {
		t.Fatal("wrong value in cache")
	}
}

func TestSetGetNil(t *testing.T) {
	c := newcache()
	c.set("key", nil, NoExpiration)
	v, ok := c.get("key")
	if !ok {
		t.Fatal("key must be found")
	}
	if v != nil {
		t.Fatal("wrong value in cache")
	}
}

func TestSetGetDict(t *testing.T) {
	c := newcache()
	c.set("key", map[string]interface{}{"a": "b", "b": "c"}, NoExpiration)
	v, ok := c.get("key")
	if !ok {
		t.Fatal("key must be found")
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		t.Fatal("wrong value in cache")
	}
	if m["a"] != "b" || m["b"] != "c" {
		t.Fatal("wrong value in cache")
	}
}

func TestSetGetList(t *testing.T) {
	c := newcache()
	c.set("key", []interface{}{"a", "b"}, NoExpiration)
	v, ok := c.get("key")
	if !ok {
		t.Fatal("key must be found")
	}
	m, ok := v.([]interface{})
	if !ok {
		t.Fatal("wrong value in cache")
	}
	if len(m) != 2 || m[0] != "a" || m[1] != "b" {
		t.Fatal("wrong value in cache")
	}
}

func TestMultiKeys(t *testing.T) {
	c := newcache()
	c.set("key1", "ups1", NoExpiration)
	c.set("key2", "ups2", NoExpiration)
	c.set("key3", "ups3", NoExpiration)
	v, ok := c.get("key2")
	if !ok {
		t.Fatal("key must be found")
	}
	if v != "ups2" {
		t.Fatal("wrong value in cache")
	}
}

func TestRemove(t *testing.T) {
	c := newcache()
	c.set("key1", "ups1", NoExpiration)
	c.set("key2", "ups2", NoExpiration)
	c.set("key3", "ups3", NoExpiration)
	c.remove("key2")
	v, ok := c.get("key2")
	if ok {
		t.Fatal("key must not be found")
	}
	if v != nil {
		t.Fatal("found removed value")
	}
}

func TestKeys(t *testing.T) {
	c := newcache()
	c.set("key1", "ups1", NoExpiration)
	c.set("key2", "ups2", NoExpiration)
	c.set("key3", "ups3", NoExpiration)
	k := c.keys()
	sort.Strings(k)
	if len(k) != 3 || k[0] != "key1" || k[1] != "key2" || k[2] != "key3" {
		t.Fatal("invalid key set")
	}
}

// TODO use https://github.com/benbjohnson/clock
func TestTTL(t *testing.T) {
	c := newcache()
	c.set("key", "ups", 2*time.Millisecond)
	time.Sleep(time.Millisecond)
	v, ok := c.get("key")
	if !ok {
		t.Fatal("key must be found")
	}
	if v != "ups" {
		t.Fatal("wrong value in cache")
	}
	time.Sleep(2 * time.Millisecond)
	v, ok = c.get("key")
	if ok {
		t.Fatal("found expired value")
	}
}

func TestManyEpoches(t *testing.T) {
	c := newcache()
	c.set("key", "ups", NoExpiration)
	for i := 0; i < 2*buckets; i++ {
		c.nextEpoch()
	}
	v, ok := c.get("key")
	if !ok {
		t.Fatal("key must be found")
	}
	if v != "ups" {
		t.Fatal("wrong value in cache")
	}
}

func TestGetNextEpoch(t *testing.T) {
	c := newcache()
	c.set("key", "ups", NoExpiration)
	c.nextEpoch()
	v, ok := c.get("key")
	if !ok {
		t.Fatal("key must be found")
	}
	if v != "ups" {
		t.Fatal("wrong value in cache")
	}
}

func TestRemoveNextEpoch(t *testing.T) {
	c := newcache()
	c.set("key", "ups", NoExpiration)
	c.nextEpoch()
	v, ok := c.get("key")
	if !ok {
		t.Fatal("key must be found")
	}
	if v != "ups" {
		t.Fatal("wrong value in cache")
	}

	c.remove("key")
	if _, ok := c.get("key"); ok {
		t.Fatal("key must not be found")
	}
}

func TestKeysEpoches(t *testing.T) {
	c := newcache()
	c.set("key1", "ups1", NoExpiration)
	c.nextEpoch()
	c.set("key2", "ups2", NoExpiration)
	c.nextEpoch()
	c.set("key3", "ups3", NoExpiration)
	k := c.keys()
	sort.Strings(k)
	if len(k) != 3 || k[0] != "key1" || k[1] != "key2" || k[2] != "key3" {
		t.Fatal("invalid key set")
	}
}
