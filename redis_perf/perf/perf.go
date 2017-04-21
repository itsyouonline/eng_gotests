package perf

import (
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"

	"docs.greenitglobe.com/despiegk/gotests/redis_perf/redis"
)

// storeDataHSetRandom stores a defined amount of objects with a given size in a redis hset.
// the key is generated from the current unix timestamp.
func StoreDataHSetRandom(amount, size int, client redis.RedisClient) error {
	log.Infof("Storing %v objects in a redis HSet, %v bytes per object.", amount, size)
	// make the byteslice to hold our data
	data := make([]byte, size)
	// and fill it with some junk
	for i := 0; i < len(data); i++ {
		data[i] = byte(i % 256)
	}
	log.Debugf("Junk data initialized, data size is %v bytes", len(data))
	// the key for the HSet
	key := strconv.FormatInt(time.Now().Unix(), 10)
	start := time.Now()
	for i := 0; i < amount; i++ {
		if err := client.StoreInHset(key, strconv.Itoa(i), data); err != nil {
			log.Errorf("Failed to store data in HSet: %v", err)
			log.Errorf("HSet key: %v, field: %v", key, i)
			log.Error("There were errors when trying to store data in redis, aborting test")
			return err
		}
	}

	duration := time.Since(start)
	log.Infof("Using client %v", client)
	log.Infof("Stored %v objects of %v bytes, took %v", amount, len(data), duration)
	return nil
}
