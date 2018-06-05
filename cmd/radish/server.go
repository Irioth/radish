package main

import (
	"flag"
	"log"
	"net"

	"github.com/Irioth/radish"
)

var portPtr = flag.String("port", "1234", "server port")

func main() {
	flag.Parse()
	s := radish.NewServer()
	defer s.Stop()
	// TODO listen for signals for correct shutdown
	err := s.Listen(net.JoinHostPort("", *portPtr))
	log.Println(err)
}
