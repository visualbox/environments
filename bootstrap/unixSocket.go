package main

import (
	"log"
	"net"
)

func unixSocketServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		go Output(string(data))
	}
}

// InitUnixSocket ...
func InitUnixSocket() {
	l, err := net.Listen("unix", "/tmp/out")
	if err != nil {
		log.Fatal(err)
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go unixSocketServer(fd)
	}
}
