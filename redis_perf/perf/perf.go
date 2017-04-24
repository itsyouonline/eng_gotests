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
	log.Debug("Key is ", key)
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

// storeDataHSetRandom stores a defined amount of objects with a given size in a redis hset.
// the key is generated from the current unix timestamp. Pipes are used to reduce total rtt.
// Pipes buffer the specified amount of commands before being executed.
func StoreDataHSetPipeRandom(amount, pipelength, size int, client redis.RedisClient) error {
	log.Infof("Storing %v objects in a redis HSet with pipes, %v bytes per object.", amount, size)
	// make the byteslice to hold our data
	data := make([]byte, size)
	// and fill it with some junk
	for i := 0; i < len(data); i++ {
		data[i] = byte(i % 256)
	}
	log.Debugf("Junk data initialized, data size is %v bytes", len(data))
	// the key for the HSet
	key := strconv.FormatInt(time.Now().Unix(), 10)
	log.Debug("Key is ", key)
	start := time.Now()
	pipe := client.StartPipe()
	var pipeCounter int
	for i := 0; i < amount; i++ {
		if err := pipe.StoreInHset(key, strconv.Itoa(i), data); err != nil {
			log.Errorf("Failed to store data in HSet: %v", err)
			log.Errorf("HSet key: %v, field: %v", key, i)
			log.Error("There were errors when trying to store data in redis, aborting test")
			return err
		}
		// increment the counter if we add a statement
		pipeCounter++
		// execute the pipe after the selected amount of statements.
		if pipeCounter%pipelength == 0 && i != 0 {
			_, err := pipe.Execute()
			if err != nil {
				log.Error("Error while executing pipe: ", err)
				log.Errorf("HSet key: %v, field: %v, pipe counter: %v", key, i, pipeCounter)
				log.Error("There were errors when trying to store data in redis, aborting test")
				return err
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
			return err
		}
	}

	duration := time.Since(start)
	log.Infof("Using client %v, with pipe", client)
	log.Infof("Stored %v objects of %v bytes, took %v", amount, len(data), duration)
	return nil
}
