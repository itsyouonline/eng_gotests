package main

import (
	"testing"

	log "github.com/Sirupsen/logrus"
)

func benchmarkWriteOneReadOne(msgCount int, b *testing.B) {
	// Set the log level to warn so the method doesn't spamm us with info logs
	log.SetLevel(log.WarnLevel)
	for n := 0; n < b.N; n++ {
		writeOneReadOne(msgCount)
	}
}

func benchmarkWriteListRead(msgCount int, b *testing.B) {
	// Set the log level to warn so the method doesn't spamm us with info logs
	log.SetLevel(log.WarnLevel)
	for n := 0; n < b.N; n++ {
		writeListRead(msgCount)
	}
}

func BenchmarkWriteOneReadOne100(b *testing.B)     { benchmarkWriteOneReadOne(100, b) }
func BenchmarkWriteOneReadOne1000(b *testing.B)    { benchmarkWriteOneReadOne(1000, b) }
func BenchmarkWriteOneReadOne10000(b *testing.B)   { benchmarkWriteOneReadOne(10000, b) }
func BenchmarkWriteOneReadOne100000(b *testing.B)  { benchmarkWriteOneReadOne(100000, b) }
func BenchmarkWriteOneReadOne1000000(b *testing.B) { benchmarkWriteOneReadOne(1000000, b) }

func BenchmarkWriteLIstRead100(b *testing.B)     { benchmarkWriteListRead(100, b) }
func BenchmarkWriteLIstRead1000(b *testing.B)    { benchmarkWriteListRead(1000, b) }
func BenchmarkWriteLIstRead10000(b *testing.B)   { benchmarkWriteListRead(10000, b) }
func BenchmarkWriteLIstRead100000(b *testing.B)  { benchmarkWriteListRead(100000, b) }
func BenchmarkWriteLIstRead1000000(b *testing.B) { benchmarkWriteListRead(1000000, b) }
