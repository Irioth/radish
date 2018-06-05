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
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Remove(key string) error
	Keys() ([]string, error)

	GetIndex(key string, index int) (interface{}, error)
	GetDict(dictName string, key string) (interface{}, error)
}
