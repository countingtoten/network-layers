package main

import (
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Status string

const (
	Open    Status = "open"
	Unknown Status = "unknown"
	Closed  Status = "closed"
)

type scanner struct {
	portStatuses map[int]Status
	m            sync.Mutex
	wg           sync.WaitGroup
}

// Scan a port by attempting to read and write from it three times.
// If you can write to a port and it writes back an empty response it
// is definitely open.
// If writing to the port produces an error it is definitely closed.
// If you can write to a port but its response times out, it is unknown
// whether or not it is open.
func (s *scanner) Scan(port int) {
	p := strconv.Itoa(port)
	addr := net.JoinHostPort("", p)
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

	var (
		retry    int    = 0
		maxTries int    = 3
		status   Status = ""
	)
	for retry < maxTries && status != Open {
		retry++

		_, err = conn.Write([]byte{})
		if err != nil {
			status = Closed
		}

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		buf := make([]byte, 1000)
		_, err := conn.Read(buf)
		if err != nil {
			opErr, ok := err.(*net.OpError)

			if ok && opErr.Temporary() {
				status = Unknown
			} else {
				status = Closed
			}
		} else {
			status = Open
			break
		}
	}

	s.m.Lock()
	defer s.m.Unlock()
	s.portStatuses[port] = status
}

func (s *scanner) String() string {
	s.m.Lock()
	defer s.m.Unlock()

	var keys []int
	for k := range s.portStatuses {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	str := strings.Builder{}
	for _, k := range keys {
		status := s.portStatuses[k]

		if status == Open || status == Unknown {
			s := fmt.Sprintf("%d: %s\n", k, status)
			str.WriteString(s)
		}
	}

	return str.String()
}

func main() {
	s := &scanner{
		portStatuses: map[int]Status{},
	}

	conSema := make(chan struct{}, 10)
	wg := sync.WaitGroup{}

	for i := 1; i <= 200; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()

			conSema <- struct{}{}

			s.Scan(port)

			<-conSema
		}(i)
	}

	wg.Wait()
	log.Printf("Port Status:\n%s\n", s.String())
}
