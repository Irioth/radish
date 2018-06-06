package radish

import (
	"errors"
	"time"
)

var (
	NotFound   = newError("Key not found")
	NotList    = newError("Value not list")
	NotDict    = newError("Value not dictionary")
	OutOfBound = newError("Index out of bound")
	InvalidKey = newError("invalid key (key must not contains spaces)")
	BadCommand = newError("bad command")

	NoExpiration = time.Duration(0)
)

var errorsmap = make(map[string]error)

func newError(text string) error {
	err := errors.New(text)
	errorsmap[text] = err
	return err
}

type Client interface {
	// Get returns value by key or NotFound if key not exist or expired
	Get(key string) (interface{}, error)

	// Set updates key/value with provided ttl
	Set(key string, value interface{}, ttl time.Duration) error

	// Remove delete key from storage
	Remove(key string) error

	// Keys returns list of live keys in storage
	Keys() ([]string, error)

	// GetIndex returns element of stored list by index
	GetIndex(key string, index int) (interface{}, error)

	// GetDict returns element of stored dictionary by key
	GetDict(dictName string, key string) (interface{}, error)
}
