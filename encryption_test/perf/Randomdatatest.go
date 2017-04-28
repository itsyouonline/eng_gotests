package perf

import (
	"bytes"
	"crypto/rand"
	"io"
	mathrand "math/rand"
	"os"
	"runtime/pprof"
	"time"

	"docs.greenitglobe.com/despiegk/gotests/encryption_test/compressors"
	"docs.greenitglobe.com/despiegk/gotests/encryption_test/crypt"

	log "github.com/Sirupsen/logrus"
)

func RandomDataTestEverything(dataLength int) error {
	// seed the generator
	mathrand.Seed(time.Now().UnixNano())
	dataslice := make([][]byte, 100000)

	log.Debug("Generating random data to compress and encrypt")
	var err error
	start := time.Now()
	for i := range dataslice {
		dataslice[i] = generateRandomStringBytes(dataLength)
		// dataslice[i] = make([]byte, dataLength)
		// rand.Read(dataslice[i])
	}
	log.Debug("Data generated in ", time.Since(start))
	ciphers := []string{"aes", "3des", "twofish", "blowfish"}
	modes := []string{"cbc", "cfb", "ctr", "ofb", "gcm"}
	comp := []string{"lz4", "gzip", "zlib"}

	key := make([]byte, 16)
	_, err = rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}

	// dump profile in this file to later use with `go tool pprof`
	f, err := os.Create("app.cpuprof")
	if err != nil {
		// Exit if we cant create the profile
		log.Fatalf("failed to create profiling file: %v", err)
	}
	// close file when we are done
	defer f.Close()
	// start profile
	pprof.StartCPUProfile(f)
	// make sure we stop the profiling
	defer pprof.StopCPUProfile()

	for _, c := range comp {
		compressed, err := compress(dataslice, c)
		if err != nil {
			log.Errorf("Failed to compress data with %v: %v", c, err)
			return err
		}
		_, err = decompress(compressed, c)
		if err != nil {
			log.Errorf("Failed to decompress data with %v: %v", c, err)
			// return err
		}
	}

	for _, c := range ciphers {
		cipher, err := crypt.NewBlockCipher(c, key)
		if err != nil {
			log.Fatal(err)
		}

		iv := make([]byte, cipher.BlockSize())
		_, err = rand.Read(iv)
		if err != nil {
			log.Fatal(err)
		}

		for _, mode := range modes {
			if mode == "gcm" && (c == "blowfish" || c == "3des") {
				continue
			}
			encryptedData, err := encrypt(dataslice, key, iv, c, mode)
			if err != nil {
				log.Error("Error while encrypting data: ", err)
				return err
			}

			decryptedData, err := decrypt(encryptedData, key, iv, c, mode)
			if err != nil {
				log.Error("Error while decrypting data: ", err)
				return err
			}

			log.Debug("Verfiying decoded content")
			for i := range dataslice {
				if !bytes.Equal(dataslice[i], decryptedData[i]) {
					log.Fatal("Decoded output is different from the input.")
				}
			}
			if len(dataslice) != len(decryptedData) {
				log.Fatal("there are a different amount of decoded data and original data")
			}
			log.Debug("Decoded content verified")
		}
	}
	return err
}

func encrypt(data [][]byte, key []byte, iv []byte, alg string, mode string) ([][]byte, error) {
	ciph, err := crypt.NewBlockCipher(alg, key)
	if err != nil {
		return nil, err
	}
	encrypter := crypt.NewEncrypter(mode, ciph, iv)
	encrypteddata := make([][]byte, len(data))
	start := time.Now()
	for i := range data {
		encrypteddata[i] = encrypter.Encrypt(data[i])
	}
	log.Infof("%v - %v: encoded entire slice of %v items, %v bytes total, in %v.", alg, mode, len(data), getUsedSpace(data), time.Since(start))
	return encrypteddata, nil
}

func decrypt(data [][]byte, key []byte, iv []byte, alg string, mode string) ([][]byte, error) {
	ciph, err := crypt.NewBlockCipher(alg, key)
	if err != nil {
		return nil, err
	}
	decrypter := crypt.NewDecrypter(mode, ciph, iv)
	decrypteddata := make([][]byte, len(data))
	start := time.Now()
	for i := range data {
		decrypteddata[i] = decrypter.Decrypt(data[i])
	}
	log.Infof("%v - %v: decoded entire slice of %v items, %v bytes total, in %v.", alg, mode, len(data), getUsedSpace(data), time.Since(start))
	return decrypteddata, nil
}

func compress(data [][]byte, alg string) ([][]byte, error) {
	buf := bytes.NewBuffer(nil)
	input := bytes.NewReader(nil)
	log.Debug("Try to compress data with ", alg)
	compressor := compressors.NewCompressor(alg, buf)
	compressedData := make([][]byte, len(data))
	start := time.Now()
	for i := range data {
		input.Reset(data[0])
		// _, err := compressor.Write(data[i])
		_, err := io.Copy(compressor, input)
		if err != nil {
			log.Error("Failed to compress data: ", err)
			return nil, err
		}
		compressor.Close()
		compressedData[i] = buf.Bytes()
		buf.Reset()
		compressor.Reset(buf)
		if i%10000 == 0 {
			log.Debugf("Compressed %v items", i)
		}
	}
	log.Info("Compression done in ", time.Since(start))
	log.Infof("Compressed data with %v, original size: %v, new size: %v", alg, getUsedSpace(data), getUsedSpace(compressedData))
	return compressedData, nil
}

func compresscombined(data [][]byte, alg string) ([][]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	log.Debug("Try to compress all data in one with ", alg)
	compressor := compressors.NewCompressor(alg, buf)
	compressedData := make([][]byte, 1)
	allData := make([]byte, 0)
	start := time.Now()
	for i := range data {
		allData = append(allData, data[i]...)
	}
	log.Debug("Concatenated data in ", time.Since(start))
	start = time.Now()
	_, err := compressor.Write(allData)
	if err != nil {
		log.Error("Failed to compress data: ", err)
		return nil, err
	}
	err = compressor.Close()
	if err != nil {
		log.Error("Failed to close compressor: ", err)
		return nil, err
	}
	compressedData[0] = buf.Bytes()
	log.Info("Compression done in ", time.Since(start))
	log.Infof("Compressed data with %v, original size: %v, new size: %v", alg, len(allData), getUsedSpace(compressedData))
	return compressedData, nil
}

func decompress(data [][]byte, alg string) ([][]byte, error) {
	decompressedData := make([][]byte, len(data))
	targetbuf := bytes.NewBuffer(nil)
	buf := bytes.NewReader(data[0])
	decompressor, err := compressors.NewDecompressor(alg, buf)
	if err != nil {
		log.Error("Failed to create new decompressor: ", err)
		return nil, err
	}
	log.Debug("Try to decompress data using ", alg)
	start := time.Now()
	for i := range data {
		buf.Reset(data[0])
		err = decompressor.Reset(buf)
		if err != nil {
			log.Error("Failed to reset the decompressor: ", err)
			return nil, err
		}
		_, err = io.Copy(targetbuf, decompressor)
		if err != nil {
			log.Error("Failed to decompress data: ", err)
			return nil, err
		}
		decompressedData[i] = targetbuf.Bytes()
		targetbuf.Reset()
	}
	log.Info("Decompression done in ", time.Since(start))
	log.Infof("Decompressed data with %v, compressed size: %v, new size: %v", alg, getUsedSpace(data), getUsedSpace(decompressedData))
	return decompressedData, nil
}

func generateRandomStringBytes(length int) []byte {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 "
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	b := make([]byte, length)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := length-1, mathrand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = mathrand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func getUsedSpace(data [][]byte) int {
	size := 0
	for i := 0; i < len(data); i++ {
		size += len(data[i])
	}
	return size
}
