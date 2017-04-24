package redis

import (
	"strconv"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
)

func benchmarkStoreInHset(client RedisClient, b *testing.B) {
	// generate 200 bytes of data
	data := make([]byte, 200)
	// key is the current unix timestamp
	key := strconv.FormatInt(time.Now().Unix(), 10)
	// declare the error variable so we only assign inside the benchmark loop, not declare and assign
	var err error
	// maintain a spereate field counter so we dont have field collisions in case the
	// runtime decides to re run the test with the same key somehow
	var fieldCounter int
	for n := 0; n < b.N; n++ {
		err = client.StoreInHset(key, strconv.Itoa(fieldCounter), data)
		if err != nil {
			log.Errorf("Error while storing data in HSet with key %v in field %v: %v", key, fieldCounter, err)
			return
		}
		fieldCounter++
	}
}

// BenchmarkStoreInHsetGoRedisTcp benchmarks the StoreInHset function using a go-redis
// client with a tcp connection
func BenchmarkStoreInHsetGoRedisTcp(b *testing.B) {
	client := NewRedisClient("go-redis", "localhost:6379", Tcp)
	benchmarkStoreInHset(client, b)
}

// BenchmarkStoreInHsetGoRedisUnix benchmarks the StoreInHset function using a go-redis
// client with a unix socket
func BenchmarkStoreInHsetGoRedisUnix(b *testing.B) {
	client := NewRedisClient("go-redis", "/tmp/redis.sock", Unix)
	benchmarkStoreInHset(client, b)
}
