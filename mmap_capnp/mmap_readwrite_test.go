package main

import (
	"testing"

	log "github.com/Sirupsen/logrus"
)

func BenchmarkWriteOneReadOne(b *testing.B) {
	// Set the log level to warn so the method doesn't spamm us with info logs
	log.SetLevel(log.WarnLevel)
	for n := 0; n < b.N; n++ {
		// benchmark writing 1K messages
		writeOneReadOne(1000000)
	}
}

func BenchmarkWriteListRead(b *testing.B) {
	// Set the log level to warn so the method doesn't spamm us with info logs
	log.SetLevel(log.WarnLevel)
	for n := 0; n < b.N; n++ {
		writeListRead(1000000)
	}
}
