package radish

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	radish *Radish
}

func NewServer() *Server {
	return &Server{
		radish: NewLocal(),
	}
}

func (s *Server) Stop() {
	s.radish.Stop()
}

func (s *Server) Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Println("Listen on", addr)

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(c)
	}

}

func (s *Server) handleConn(c net.Conn) {
	log.Println("new connection", c.RemoteAddr().String())
	r := bufio.NewReader(c)
	for {
		cmd, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println(err.Error())
			}
			break
		}
		cmd = strings.TrimSpace(cmd)

		result, err := s.proceedCmd(cmd)
		if err != nil {
			if err := senderr(err, c); err != nil {
				log.Println(err.Error())
				break
			}
		} else if err := sendok(result, c); err != nil {
			log.Println(err.Error())
			break
		}

	}

}

func (s *Server) proceedCmd(cmd string) (string, error) {
	c := strings.SplitN(cmd, " ", 4)
	if len(c) <= 0 {
		return "", BadCommand
	}
	switch c[0] {
	case "GET":
		return s.handleGet(c)
	case "GETINDEX":
		return s.handleGetIndex(c)
	case "GETDICT":
		return s.handleGetDict(c)
	case "REMOVE":
		return s.handleRemove(c)
	case "KEYS":
		return s.handleKeys(c)
	case "SET":
		return s.handleSet(c)
	}
	return "", BadCommand
}

func (s *Server) handleGet(c []string) (string, error) {
	if len(c) <= 1 {
		return "", BadCommand
	}
	v, err := s.radish.Get(c[1])
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *Server) handleGetIndex(c []string) (string, error) {
	if len(c) <= 2 {
		return "", BadCommand
	}
	index, err := strconv.Atoi(c[2])
	if err != nil {
		return "", err
	}
	v, err := s.radish.GetIndex(c[1], index)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *Server) handleGetDict(c []string) (string, error) {
	if len(c) <= 2 {
		return "", BadCommand
	}
	v, err := s.radish.GetDict(c[1], c[2])
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *Server) handleRemove(c []string) (string, error) {
	if len(c) <= 1 {
		return "", BadCommand
	}
	return "", s.radish.Remove(c[1])
}

func (s *Server) handleKeys(c []string) (string, error) {
	keys, err := s.radish.Keys()
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(keys)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *Server) handleSet(c []string) (string, error) {
	if len(c) <= 3 {
		return "", BadCommand
	}

	ttl, err := strconv.Atoi(c[2])
	if err != nil {
		return "", err
	}

	var value interface{}
	if err := json.Unmarshal([]byte(c[3]), &value); err != nil {
		return "", err
	}

	s.radish.Set(c[1], value, time.Duration(ttl))
	return "", nil
}

func sendok(result string, w io.Writer) error {
	_, err := fmt.Fprintf(w, "OK %s\n", result)
	if err != nil {
		return err
	}
	return nil
}

func senderr(err error, w io.Writer) error {
	_, err = fmt.Fprintf(w, "ERROR %s\n", err.Error())
	if err != nil {
		return err
	}
	return nil
}
