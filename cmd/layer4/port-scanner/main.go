package main

import (
	"log"
	"sync"
	"time"

	"github.com/countingtoten/network-layers/scanner"
)

func main() {
	timeout := 5 * time.Second
	s := scanner.NewUDP(timeout)

	conSema := make(chan struct{}, 100)
	wg := sync.WaitGroup{}

	for i := 1; i <= 65535; i++ {
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
