package radish

import (
	"os"
	"sort"
	"testing"
	"time"
)

const testaddr = ":1234"

func TestMain(m *testing.M) {
	s := NewServer()
	go s.Listen(testaddr)
	time.Sleep(10 * time.Millisecond)
	ret := m.Run()
	s.Stop()
	os.Exit(ret)
}

func TestBBSimple(t *testing.T) {
	c, _ := Open(testaddr)
	defer c.Close()

	c.Set("keybbsimple", "ups", NoExpiration)
	if value, err := c.Get("keybbsimple"); err != nil || value != "ups" {
		t.Fatal("incorrect value in cache")
	}
}

func TestBBErrors(t *testing.T) {
	c, _ := Open(testaddr)
	defer c.Close()

	if _, err := c.Get("keybbserrors"); err != NotFound {
		t.Fatal("incorrect error")
	}
}

func TestBBDict(t *testing.T) {
	c, _ := Open(testaddr)
	defer c.Close()

	c.Set("keybbdict", map[string]interface{}{"key": "valuedict"}, NoExpiration)
	if value, err := c.GetDict("keybbdict", "key"); err != nil || value != "valuedict" {
		t.Fatal("incorrect value in cache")
	}
}

func TestBBIndex(t *testing.T) {
	c, _ := Open(testaddr)
	defer c.Close()

	c.Set("keybbindex", []interface{}{"valuelist1", "valuelist2"}, NoExpiration)
	if value, err := c.GetIndex("keybbindex", 1); err != nil || value != "valuelist2" {
		t.Fatal("incorrect value in cache")
	}
}

func TestBBRemove(t *testing.T) {
	c, _ := Open(testaddr)
	defer c.Close()

	c.Set("keybbremove", "val", NoExpiration)
	c.Remove("keybbremove")
	if _, err := c.GetIndex("keybbremove", 1); err != NotFound {
		t.Fatal("value must be removed")
	}
}

func TestBBKeys(t *testing.T) {
	c, _ := Open(testaddr)
	defer c.Close()

	c.Set("keybbkeys", "val", NoExpiration)

	k, _ := c.Keys()
	sort.Strings(k)
	if k[sort.SearchStrings(k, "keybbkeys")] != "keybbkeys" {
		t.Fatal("keys must contains at least keybbkeys")
	}
}
