package perf

import (
	"math/rand"
	"time"

	"docs.greenitglobe.com/despiegk/gotests/redis_perf/compression"
	"docs.greenitglobe.com/despiegk/gotests/redis_perf/encryption"

	log "github.com/Sirupsen/logrus"
)

func init() {
	//seed the generate
	rand.Seed(time.Now().UnixNano())
}

func CryptoTest(amount, size int) error {
	// make all our slices to hold intermediate values
	testData := make([][]byte, amount)
	md5a := make([][]byte, amount)
	compressed := make([][]byte, amount)
	ciphertext := make([][]byte, amount)
	md5b := make([][]byte, amount)
	var err error
	// generate the test data
	start := time.Now()
	for i := range testData {
		testData[i] = generateRandomStringBytes(size)
	}
	log.Debugf("Generated %v objects of %v bytes each in %v", amount, size, time.Since(start))
	// hash the test data
	start = time.Now()
	for i := range testData {
		md5a[i] = encryption.Hash(testData[i])
	}
	log.Infof("Hashed %v objects in %v", amount, time.Since(start))
	// compress the test data
	start = time.Now()
	for i := range testData {
		if i < 10 {
			log.Warn(string(testData[i]))
		}
		compressed[i], err = compression.Compress(testData[i])
		if err != nil {
			log.Errorf("Error during compression of item %v: %v", i, err)
			return err
		}
		if i < 10 {
			log.Warn(compressed[i])
		}
		if i%1000 == 999 {
			log.Debugf("Compressed %v'th item", i+1)
		}
	}
	duration := time.Since(start)
	log.Infof("Compressed %v objects in %v, average %v per object", amount, duration, time.Duration(int(duration)/amount))
	log.Infof("Old size total: %v, new size total: %v", getUsedSpace(testData), getUsedSpace(compressed))
	// encrypt the compressed data
	start = time.Now()
	for i := range testData {
		ciphertext[i], err = encryption.Encrypt(md5a[i], compressed[i])
		if err != nil {
			log.Errorf("Error during encryption of item %v: %v", i, err)
			return err
		}
	}
	duration = time.Since(start)
	log.Infof("Encrypted %v objects in %v, average %v per object", amount, duration, time.Duration(int(duration)/amount))
	log.Infof("Old size total: %v, new size total: %v", getUsedSpace(compressed), getUsedSpace(ciphertext))
	// hash the encrypted data
	start = time.Now()
	for i := range testData {
		md5b[i] = encryption.Hash(ciphertext[i])
	}
	log.Infof("Hashed %v objects in %v", amount, time.Since(start))
	return nil
}

// generateRandomStringBytes generates a random string with the specified length
func generateRandomStringBytes(length int) []byte {
	// declare our constants
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 " // all valid letters
		letterIdxBits = 6                                                                 // 6 bits to represent a letter index (64 values only 63 letters)
		letterIdxMask = 1<<letterIdxBits - 1                                              // positive bitmask for letterIdxBits
		letterIdxMax  = 63 / letterIdxBits                                                // # of letter indices fitting in 63 bits
	)
	// byteslice that will represent our string
	b := make([]byte, length)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remainder := length-1, rand.Int63(), letterIdxMax; i >= 0; {
		// if we don't have enough uses from our cache anymore, refill it
		if remainder == 0 {
			cache, remainder = rand.Int63(), letterIdxMax
		}
		// if the int value represented by the last `letterIdxBits` is a valid index in the letterBytes,
		// add the letter to our string
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		// push the bits we used out of the cache
		cache >>= letterIdxBits
		// make sure we remove a use from the cache
		remainder--
	}
	// don't cast to a string as we are mainly working with byte slices anyway
	return b
}

func getUsedSpace(data [][]byte) int {
	size := 0
	for i := 0; i < len(data); i++ {
		size += len(data[i])
	}
	return size
}
