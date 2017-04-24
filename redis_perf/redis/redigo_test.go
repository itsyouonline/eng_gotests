package redis

import (
	"testing"
)

// BenchmarkStoreInHsetRedigoTcp benchmarks the StoreInHset function using a redigo
// client with a tcp connection
func BenchmarkStoreInHsetRedigoTcp(b *testing.B) {
	client := NewRedisClient("redigo", "localhost:6379", Tcp)
	benchmarkStoreInHset(client, b)
}

// BenchmarkStoreInHsetRedigoUnix benchmarks the StoreInHset function using a redigo
// client with a unix socket
func BenchmarkStoreInHsetRedigoUnix(b *testing.B) {
	client := NewRedisClient("redigo", "/tmp/redis.sock", Unix)
	benchmarkStoreInHset(client, b)
}
