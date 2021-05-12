package scanner_test

import (
	"testing"
	"time"

	"../scanner"
)

func BenchmarkScan(t *testing.B) {
	s := scanner.NewUDP(5 * time.Nanosecond)

	for i := 0; i < t.N; i++ {
		s.Scan(i)
	}
}
