package main

import (
	"encoding/hex"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type scanner struct {
	openPorts []string
	m         sync.Mutex
	wg        sync.WaitGroup
}

func (s *scanner) Scan(port string) {
	log.Printf("Scanning %s\n", port)
	addr := net.JoinHostPort("127.0.0.1", port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	b := []byte{02, 00, 00, 00, 45, 00, 00, 1c, 51, 3a, 00, 00, 2d, 11, 00, 00, 7f, 00, 00, 01, 7f, 00, 00, 01, c1, d5, 00, 8a, 00, 08, 3f, 7c}
	if err != nil {
		log.Println(err)
	}

	conn.Write(b)

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	buf := make([]byte, 1000)
	n, err := conn.Read(buf)
	if err == nil {
		s.m.Lock()
		defer s.m.Unlock()
		s.openPorts = append(s.openPorts, port)
	} else {
		log.Println(err)
	}

	log.Println(n)
}

func main() {
	s := &scanner{}

	conSema := make(chan struct{}, 10)
	wg := sync.WaitGroup{}

	for i := 137; i <= 138; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()

			conSema <- struct{}{}

			p := strconv.Itoa(port)
			s.Scan(p)

			<-conSema
		}(i)
	}

	wg.Wait()
	log.Printf("Open Ports:\n%v\n", s.openPorts)
}
