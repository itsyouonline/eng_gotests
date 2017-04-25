package redis

import (
	"strconv"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
)

// store up to 10 statements before executing them
const pipelength = 10

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

func benchmarkStoreInHsetWithPipe(client RedisClient, b *testing.B) {
	// generate 200 bytes of data
	data := make([]byte, 200)
	// key is the current unix timestamp
	key := strconv.FormatInt(time.Now().Unix(), 10)
	// declare the error variable so we only assign inside the benchmark loop, not declare and assign
	var err error
	// maintain a spereate field counter so we dont have field collisions in case the
	// runtime decides to re run the test with the same key somehow
	var fieldCounter, pipeCounter int
	pipe := client.StartPipe()
	for n := 0; n < b.N; n++ {
		err = pipe.StoreInHset(key, strconv.Itoa(fieldCounter), data)
		if err != nil {
			log.Errorf("Error while storing data in HSet with key %v in field %v: %v", key, fieldCounter, err)
			return
		}
		fieldCounter++
		pipeCounter++
		if pipeCounter%pipelength == 0 && n != 0 {
			_, err := pipe.Execute()
			if err != nil {
				log.Error("Error while executing pipe: ", err)
				return
			}
			// reset the counter now that we executed the pipe
			pipeCounter = 0
		}
	}
	// If the pipe is not empty, execute the remaining statements.
	if pipeCounter%pipelength == 0 && pipelength == 0 {
		_, err := pipe.Execute()
		if err != nil {
			log.Error("Error while executing remainder in pipe: ", err)
			return
		}
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

// BenchmarkStoreInHsetGoRedisTcpWithPipe benchmarks the StoreInHset function using a go-redis
// client with a tcp connection and using a pipe
func BenchmarkStoreInHsetGoRedisTcpWithPipe(b *testing.B) {
	client := NewRedisClient("go-redis", "localhost:6379", Tcp)
	benchmarkStoreInHsetWithPipe(client, b)
}

// BenchmarkStoreInHsetGoRedisUnixWithPipe benchmarks the StoreInHset function using a go-redis
// client with a unix socket and using a pipe
func BenchmarkStoreInHsetGoRedisUnixWithPipe(b *testing.B) {
	client := NewRedisClient("go-redis", "/tmp/redis.sock", Unix)
	benchmarkStoreInHsetWithPipe(client, b)
}
