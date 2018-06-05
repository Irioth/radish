package radish

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type RadishClient struct {
	conn   net.Conn
	reader *bufio.Reader
}

var _ Client = &RadishClient{}

func Open(addr string) (*RadishClient, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &RadishClient{c, bufio.NewReader(c)}, nil
}

func (r *RadishClient) Close() error {
	return r.conn.Close()
}

func (r *RadishClient) Get(key string) (interface{}, error) {
	if hasSpaces(key) {
		return nil, InvalidKey
	}
	fmt.Fprintf(r.conn, "GET %s\n", key)
	return r.readAnswer()
}

func (r *RadishClient) Set(key string, value interface{}, ttl time.Duration) error {
	if hasSpaces(key) {
		return InvalidKey
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	fmt.Fprintf(r.conn, "SET %s %d %s\n", key, ttl, string(data))
	_, err = r.readAnswer()
	return err
}

func (r *RadishClient) Remove(key string) error {
	if hasSpaces(key) {
		return InvalidKey
	}
	fmt.Fprintf(r.conn, "REMOVE %s\n", key)
	_, err := r.readAnswer()
	return err
}

func (r *RadishClient) Keys() ([]string, error) {
	fmt.Fprintf(r.conn, "KEYS\n")
	v, err := r.readAnswer()
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, val := range v.([]interface{}) {
		result = append(result, val.(string))
	}
	return result, nil
}

func (r *RadishClient) GetIndex(key string, index int) (interface{}, error) {
	if hasSpaces(key) {
		return nil, InvalidKey
	}
	fmt.Fprintf(r.conn, "GETINDEX %s %d\n", key, index)
	return r.readAnswer()
}
func (r *RadishClient) GetDict(dictName string, key string) (interface{}, error) {
	if hasSpaces(dictName) || hasSpaces(key) {
		return nil, InvalidKey
	}
	fmt.Fprintf(r.conn, "GETDICT %s %s\n", dictName, key)
	return r.readAnswer()
}

func (r *RadishClient) readAnswer() (interface{}, error) {
	result, err := r.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	result = result[:len(result)-1]

	s := strings.SplitN(result, " ", 2)
	switch s[0] {
	case "OK":
		if len(s[1]) == 0 {
			return nil, nil
		}
		var v interface{}
		if err := json.Unmarshal([]byte(s[1]), &v); err != nil {
			return nil, err
		}
		return v, nil
	case "ERROR":
		if err, ok := errorsmap[s[1]]; ok {
			return nil, err
		}
		return nil, errors.New(s[1])
	}
	panic("ups. invalid server")
}

// TODO implement escaping in proto
func hasSpaces(s string) bool {
	return strings.IndexByte(s, ' ') >= 0
}
