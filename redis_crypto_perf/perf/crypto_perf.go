package perf

import (
	"bytes"
	"io/ioutil"
	"os"
	"time"

	"docs.greenitglobe.com/despiegk/gotests/redis_crypto_perf/compression"
	"docs.greenitglobe.com/despiegk/gotests/redis_crypto_perf/encryption"
	"docs.greenitglobe.com/despiegk/gotests/redis_perf/redis"

	log "github.com/Sirupsen/logrus"
)

// CryptoTest gets timing and data size on different actions while storing dat in redis
func CryptoTest(amount, size int, lz4 bool, rc redis.RedisClient) error {
	// make all our slices to hold intermediate values
	md5a := make([][]byte, amount)
	objects := make([][]byte, amount)
	md5b := make([][]byte, amount)
	var err error
	var usedSpaceData, usedSpaceCompressed, usedSpaceCiphertext int
	var usedSpacePlaintext, usedSpaceDecompressed int

	// generate the test data
	start := time.Now()
	testData := generateDataFromBinary(amount, size)
	log.Debugf("Generated %v objects of %v bytes each in %v, total size %v bytes",
		amount, size, time.Since(start), getUsedSpace(testData))
	usedSpaceData = getUsedSpace(testData)

	// hash the test data
	start = time.Now()
	for i := range testData {
		md5a[i] = encryption.Hash(testData[i])
	}
	duration := time.Since(start)
	log.Infof("Hashed %v objects in %v, average %v per object", amount, duration, time.Duration(int(duration)/amount))

	// compress the test data
	start = time.Now()
	for i := range testData {
		if lz4 {
			objects[i], err = compression.Lz4Compress(testData[i])
		} else {
			objects[i], err = compression.GzipCompress(testData[i])
		}
		if err != nil {
			log.Errorf("Error during compression of item %v: %v", i, err)
			return err
		}
		if i%1000 == 999 {
			log.Debugf("Compressed %v'th item", i+1)
		}
	}
	duration = time.Since(start)
	usedSpaceCompressed = getUsedSpace(objects)
	log.Infof("Compressed %v objects in %v, average %v per object", amount, duration, time.Duration(int(duration)/amount))
	log.Infof("Old size total: %v, new size total: %v", usedSpaceData, usedSpaceCompressed)

	// encrypt the compressed data
	start = time.Now()
	for i := range objects {
		objects[i], err = encryption.Encrypt(md5a[i], objects[i])
		if err != nil {
			log.Errorf("Error during encryption of item %v: %v", i, err)
			return err
		}
	}
	duration = time.Since(start)
	usedSpaceCiphertext = getUsedSpace(objects)
	log.Infof("Encrypted %v objects in %v, average %v per object", amount, duration, time.Duration(int(duration)/amount))
	log.Infof("Old size total: %v, new size total: %v", usedSpaceCompressed, usedSpaceCiphertext)

	// hash the encrypted data
	start = time.Now()
	for i := range objects {
		md5b[i] = encryption.Hash(objects[i])
	}
	duration = time.Since(start)
	log.Infof("Hashed %v objects in %v, average %v per object", amount, duration, time.Duration(int(duration)/amount))

	// get redis memory usage before we store the data
	initialMemUsage, err := rc.GetMemUsage()
	if err != nil {
		log.Error("Failed to get redis memory usage: ", err)
		return err
	}

	log.Infof("Initial redis memory consumption: %v bytes", initialMemUsage)

	// store data in redis
	start = time.Now()
	for i := range objects {
		err = rc.StoreInHset(string(md5b[i]), string(md5a[i]), objects[i])
		if err != nil {
			log.Errorf("Failed to store item %v in HSET with key %v in field %vv: %v",
				i, string(md5b[i]), string(md5a[i]), err)
		}
	}
	duration = time.Since(start)
	log.Infof("Stored %v objects in %v, average %v per object", amount, duration, time.Duration(int(duration)/amount))

	// get redis memory usage now that everything is stored
	storedMemUsage, err := rc.GetMemUsage()
	if err != nil {
		log.Error("Failed to get redis memory usage: ", err)
		return err
	}

	log.Infof("Redis memory consumption after storing data: %v bytes", storedMemUsage)

	redisMemUsage := storedMemUsage - initialMemUsage
	log.Infof("Redis used %v bytes, to store %v bytes worth of data (%.2f%%)",
		redisMemUsage, usedSpaceCiphertext, float64(100)*(float64(redisMemUsage)/float64(usedSpaceCiphertext)))
	// TODO: RELOAD DATA FROM REDIS

	// decrypt data
	start = time.Now()
	for i := range objects {
		objects[i], err = encryption.Decrypt(md5a[i], objects[i])
		if err != nil {
			log.Errorf("Error during encryption of item %v: %v", i, err)
			return err
		}
	}
	duration = time.Since(start)
	usedSpacePlaintext = getUsedSpace(objects)
	log.Infof("Decrypted %v objects in %v, average %v per object", amount, duration, time.Duration(int(duration)/amount))
	log.Infof("Old size total: %v, new size total: %v", usedSpaceCiphertext, usedSpacePlaintext)

	// decompress plaintext
	start = time.Now()
	for i := range objects {
		if lz4 {
			objects[i], err = compression.Lz4Decompress(objects[i])
		} else {
			objects[i], err = compression.GzipDecompress(objects[i])
		}
		if err != nil {
			log.Errorf("Error during decompression of item %v: %v", i, err)
			return err
		}
		if i%1000 == 999 {
			log.Debugf("Decompressed %v'th item", i+1)
		}
	}
	duration = time.Since(start)
	usedSpaceDecompressed = getUsedSpace(objects)
	log.Infof("Decompressed %v objects in %v, average %v per object", amount, duration, time.Duration(int(duration)/amount))
	log.Infof("Old size total: %v, new size total: %v", usedSpacePlaintext, usedSpaceDecompressed)

	// verify data
	start = time.Now()
	for i := range objects {
		if !bytes.Equal(objects[i], testData[i]) {
			log.Errorf("Data section %v does not match the original", i)
		}
	}
	duration = time.Since(start)
	log.Infof("Verified %v objects in %v, average %v per object", amount, duration, time.Duration(int(duration)/amount))
	log.Info("All objects are correct")

	return nil
}

// generateDataFromBinary generates our data from the current binary
func generateDataFromBinary(amount, size int) [][]byte {
	data := make([][]byte, amount)
	// read the binary into memory
	binary, err := ioutil.ReadFile(os.Args[0])
	if err != nil {
		log.Fatal("Failed to open binary: ", err)
	}
	if len(binary) < size {
		log.Fatal("Binary too small")
	}
	// chop up the binary to fit into our test data, repeat segments if required
	for i, ix := 0, 0; i < len(data); i++ {
		// make sure we still have a full sized chunk to copy
		if (ix+1)*size > len(binary) {
			ix = 0
		}
		// copy the data
		data[i] = make([]byte, size)
		copy(data[i], binary[ix*size:(ix+1)*size])
		ix++
	}
	return data
}

func getUsedSpace(data [][]byte) int {
	size := 0
	for i := 0; i < len(data); i++ {
		size += len(data[i])
	}
	return size
}
