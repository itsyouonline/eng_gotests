package redis

import (
	"testing"
)

// BenchmarkStoreInHsetRedigoTcp benchmarks the StoreInHset function using a radix.v2
// client with a tcp connection
func BenchmarkStoreInHsetRadixTcp(b *testing.B) {
	client := NewRedisClient("radix", "localhost:6379", Tcp)
	benchmarkStoreInHset(client, b)
}

// BenchmarkStoreInHsetRedigoUnix benchmarks the StoreInHset function using a radix.v2
// client with a unix socket
func BenchmarkStoreInHsetRadixUnix(b *testing.B) {
	client := NewRedisClient("radix", "/tmp/redis.sock", Unix)
	benchmarkStoreInHset(client, b)
}

// BenchmarkStoreInHsetRedigoTcpWithPipe benchmarks the StoreInHset function using a radix.v2
// client with a tcp connection and using a pipe
func BenchmarkStoreInHsetRadixTcpWithPipe(b *testing.B) {
	client := NewRedisClient("radix", "localhost:6379", Tcp)
	benchmarkStoreInHsetWithPipe(client, b)
}

// BenchmarkStoreInHsetRedigoUnixWithPipe benchmarks the StoreInHset function using a radix.v2
// client with a unix socket and using a pipe
func BenchmarkStoreInHsetRadixUnixWithPipe(b *testing.B) {
	client := NewRedisClient("radix", "/tmp/redis.sock", Unix)
	benchmarkStoreInHsetWithPipe(client, b)
}
