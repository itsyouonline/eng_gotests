package perf

import (
	"bytes"
	"crypto/rand"
	"io"
	"time"

	capnp "docs.greenitglobe.com/despiegk/gotests/encryption_test/capnp"

	"docs.greenitglobe.com/despiegk/gotests/encryption_test/compressors"
	"docs.greenitglobe.com/despiegk/gotests/encryption_test/crypt"

	log "github.com/Sirupsen/logrus"
)

func TestIndividual(dataSet, dataLength int) error {
	capnp.SetDataSize(dataLength)

	log.Debug("Generating random data to compress and encrypt")
	start := time.Now()
	dataslice := capnp.GenerateBlocks(dataSet / 100)
	log.Debug("Data generated in ", time.Since(start))
	return runTest(dataslice, false)
}

func TestCombined(dataSet, dataLength int) error {
	capnp.SetDataSize(dataLength)

	log.Debug("Generating random data to compress and encrypt")
	start := time.Now()
	dataslice := capnp.GenerateBlocks(dataSet)
	log.Debug("Data generated in ", time.Since(start))
	return runTest(dataslice, true)
}

func TestList(dataSet, dataLength int) error {
	capnp.SetDataSize(dataLength)

	log.Debug("Generating random data to compress and encrypt")
	start := time.Now()
	dataSlice := make([][]byte, 1)
	dataSlice[0] = capnp.GenerateList(dataSet)
	log.Debug("Data generated in ", time.Since(start))
	return runTest(dataSlice, false)
}

func runTest(dataslice [][]byte, combineSlice bool) error {
	ciphers := []string{"aes", "3des", "twofish", "blowfish"}
	modes := []string{"cbc", "cfb", "ctr", "ofb", "gcm"}
	comp := []string{"lz4", "gzip", "zlib"}

	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}

	// concetenate a slice of byteslices
	if combineSlice {
		allData := make([]byte, 0)
		start := time.Now()
		for i := range dataslice {
			allData = append(allData, dataslice[i]...)
		}
		log.Debug("Concatenated data in ", time.Since(start))
		dataslice = make([][]byte, 1)
		dataslice[0] = allData
	}

	// test our compression methods
	for _, c := range comp {
		compressed, err := compress(dataslice, c)
		if err != nil {
			log.Errorf("Failed to compress data with %v: %v", c, err)
			return err
		}
		dec, err := decompress(compressed, c)
		if err != nil {
			log.Errorf("Failed to decompress data with %v: %v", c, err)
			return err
		}

		// verify decoded data
		for i := range dataslice {
			if !bytes.Equal(dataslice[i], dec[i]) {
				log.Errorf("Decompressed block %v does not match its original.", i)
			}
		}
	}

	// now test the encryption methods
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

		// in all the supported modes
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
	// create our cipher
	ciph, err := crypt.NewBlockCipher(alg, key)
	if err != nil {
		return nil, err
	}
	// supply mode and iv
	encrypter := crypt.NewEncrypter(mode, ciph, iv)
	// create the output buffer
	encrypteddata := make([][]byte, len(data))
	// start timing
	start := time.Now()
	// encrypt all the data
	for i := range data {
		encrypteddata[i] = encrypter.Encrypt(data[i])
	}
	log.Infof("%v - %v: encoded entire slice of %v items, %v bytes total, in %v.", alg, mode, len(data), getUsedSpace(data), time.Since(start))
	return encrypteddata, nil
}

func decrypt(data [][]byte, key []byte, iv []byte, alg string, mode string) ([][]byte, error) {
	// create our cipher
	ciph, err := crypt.NewBlockCipher(alg, key)
	if err != nil {
		return nil, err
	}
	// supply mode and iv
	decrypter := crypt.NewDecrypter(mode, ciph, iv)
	// create the output buffer
	decrypteddata := make([][]byte, len(data))
	// timing
	start := time.Now()
	// decrypt all the data
	for i := range data {
		decrypteddata[i] = decrypter.Decrypt(data[i])
	}
	log.Infof("%v - %v: decoded entire slice of %v items, %v bytes total, in %v.", alg, mode, len(data), getUsedSpace(data), time.Since(start))
	return decrypteddata, nil
}

func compress(data [][]byte, alg string) ([][]byte, error) {
	// declare input and output buffers
	buf := bytes.NewBuffer(nil)
	input := bytes.NewReader(nil)
	log.Debug("Try to compress data with ", alg)
	// create compressor of the right type
	compressor := compressors.NewCompressor(alg, buf)
	// the final output destination
	compressedData := make([][]byte, len(data))
	// timers
	var spend time.Duration
	var start time.Time
	for i := range data {
		// compress the i'th piece of data
		input.Reset(data[i])
		// time compression
		start = time.Now()
		_, err := io.Copy(compressor, input)
		if err != nil {
			log.Error("Failed to compress data: ", err)
			return nil, err
		}
		// write header and flush
		err = compressor.Close()
		if err != nil {
			log.Error("Failed to close compressor: ", err)
			return nil, err
		}
		// add the time it took to compress to our timer
		spend += time.Since(start)
		// copy the compressed data
		compressedData[i] = buf.Bytes()
		// clear the buffer
		buf.Reset()
		// reset the compressor since we closed it
		compressor.Reset(buf)
		if i%1000 == 999 {
			log.Debugf("Compressed %v items", i+1)
		}
	}
	log.Info("Compression done in ", spend)
	log.Infof("Compressed data with %v, original size: %v, new size: %v", alg, getUsedSpace(data), getUsedSpace(compressedData))
	return compressedData, nil
}

func decompress(data [][]byte, alg string) ([][]byte, error) {
	// create the final output buffer
	decompressedData := make([][]byte, len(data))
	// intermediate buffers
	targetbuf := bytes.NewBuffer(nil)
	buf := bytes.NewReader(data[0])
	// create a new decompressor of the right type
	decompressor, err := compressors.NewDecompressor(alg, buf)
	if err != nil {
		log.Error("Failed to create new decompressor: ", err)
		return nil, err
	}
	log.Debug("Try to decompress data using ", alg)
	// timers
	var spend time.Duration
	var start time.Time
	for i := range data {
		if i%1000 == 999 {
			log.Debugf("Decompressing %v'th item", i+1)
		}
		// load the i'th piece of data
		buf.Reset(data[i])
		// reset the decompressor
		err = decompressor.Reset(buf)
		if err != nil {
			log.Error("Failed to reset the decompressor: ", err)
			return nil, err
		}
		// decompress and measure time
		start = time.Now()
		check, err := io.Copy(targetbuf, decompressor)
		if err != nil {
			log.Errorf("Failed to decompress data at %v after writing %v bytes: %v", i, check, err)
			return decompressedData, err
		}
		// add the decompression time to our timer
		spend += time.Since(start)
		// store decompressed data
		decompressedData[i] = targetbuf.Bytes()
		// clear buffer holding the data for i'th block
		targetbuf.Reset()
	}
	log.Info("Decompression done in ", spend)
	log.Infof("Decompressed data with %v, compressed size: %v, new size: %v", alg, getUsedSpace(data), getUsedSpace(decompressedData))
	return decompressedData, nil
}

func getUsedSpace(data [][]byte) int {
	size := 0
	for i := 0; i < len(data); i++ {
		size += len(data[i])
	}
	return size
}
