package main

import (
	"net"
)

func handleConnection(c net.Conn) {
	c.Write([]byte("HTTP/1.1 200 OK\nContent-Type: text/html\n\n<HTML>Hello World</HTML>\n"))
	c.Close()
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}
