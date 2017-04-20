package main

import (
	"bytes"
	"strconv"
	"time"

	capnp "zombiezen.com/go/capnproto2"

	log "github.com/Sirupsen/logrus"

	redis "github.com/go-redis/redis"
)

func storeAndReadCapnpInHset(client *redis.Client, amount int) {
	// get a random hset key, use the current unix timestamp to avoid collisions, format to string
	hkey := strconv.FormatInt(time.Now().Unix(), 10)
	// a new buffer
	buf := new(bytes.Buffer)

	log.Debug("HSet key is the current timestamp: ", string(hkey))
	for i := 0; i < amount; i++ {
		// reset the buffer
		buf.Reset()
		// create a new msg with a tlog block
		_, msg, err := generateSequenceBlock(i)
		if err != nil {
			log.Errorf("Failed to generate block %v: %v", i, err)
			return
		}
		// encode the block in the buffer
		capnp.NewEncoder(buf).Encode(msg)
		// store the encoded block in a redis hset in the field identified by its sequence
		client.HSet(string(hkey), strconv.Itoa(i), buf.Bytes())
	}
	log.Infof("Stored %v capnp messages in a redis HSet with key %v", amount, string(hkey))

	log.Debug("Get some of the messages to verify that we stored them correctly")

	for i := 0; i < amount; i++ {
		// only verify every 100th block
		if i%100 == 0 {
			// get the encoded block from the hset
			cmd := client.HGet(hkey, strconv.Itoa(i))
			if cmd.Err() != nil {
				log.Errorf("Failed to get the value stored in field %v: %v", i, cmd.Err())
				return
			}
			// decode the message
			msg, err := capnp.NewDecoder(bytes.NewBufferString(cmd.Val())).Decode()
			if err != nil {
				log.Infof("decode failed: %v", err)
				return
			}

			// get the block from the message
			block, err := ReadRootTlogBlock(msg)
			if err != nil {
				log.Errorf("decode failed to read root block: %v", err)
				return
			}
			// match the block sequence to the expected value
			if int(block.Sequence()) != i {
				log.Errorf("Decoded block has the wrong value (%v), expected %v", block.Sequence(), i)
				return
			}
			log.Debug("Correctly decoded block ", i)
		}
	}
	log.Info("Blocks stored successfully, all verification checks passed")
}
