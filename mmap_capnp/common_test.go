package main

import (
	"testing"

	log "github.com/Sirupsen/logrus"
)

var block *TlogBlock

// benchmarkCreateBlock benchmark the time it takes to create blocks with a given
// data length
func benchmarkCreateBlock(dataLenght int, b *testing.B) {
	optDataLen = dataLenght
	for n := 0; n < b.N; n++ {
		newblock, _, err := createBlock(n)
		if err != nil {
			log.Fatal("Failed to create block")
		}
		// assign to global variable to avoid aggressive compilar optimization
		block = newblock
	}
	// reset optDataLen for further tests
	optDataLen = 0
}

// execute benchmarks with diferent data lengths
func BenchmarkCreateBlock0(b *testing.B)   { benchmarkCreateBlock(0, b) }
func BenchmarkCreateBlock1(b *testing.B)   { benchmarkCreateBlock(1, b) }
func BenchmarkCreateBlock2(b *testing.B)   { benchmarkCreateBlock(2, b) }
func BenchmarkCreateBlock4(b *testing.B)   { benchmarkCreateBlock(4, b) }
func BenchmarkCreateBlock8(b *testing.B)   { benchmarkCreateBlock(8, b) }
func BenchmarkCreateBlock16(b *testing.B)  { benchmarkCreateBlock(16, b) }
func BenchmarkCreateBlock32(b *testing.B)  { benchmarkCreateBlock(32, b) }
func BenchmarkCreateBlock64(b *testing.B)  { benchmarkCreateBlock(64, b) }
func BenchmarkCreateBlock128(b *testing.B) { benchmarkCreateBlock(128, b) }
