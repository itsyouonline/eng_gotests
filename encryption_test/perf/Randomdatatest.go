package perf

import (
	"bytes"
	"crypto/rand"
	"time"

	"docs.greenitglobe.com/despiegk/gotests/encryption_test/compressors"
	"docs.greenitglobe.com/despiegk/gotests/encryption_test/crypt"

	log "github.com/Sirupsen/logrus"
)

func RandomDataTestEverything(dataLength int) error {
	dataslice := make([][]byte, 1000000)

	log.Debug("Generating random data to encrypt")
	var err error
	for i := range dataslice {
		dataslice[i] = make([]byte, dataLength)
		_, err = rand.Read(dataslice[i])
	}
	log.Debug("Data generated")
	ciphers := []string{"aes", "3des", "twofish", "blowfish"}
	modes := []string{"cbc", "cfb", "ctr", "ofb", "gcm"}
	// comp := []string{"lz4", "gzip", "zlib"}

	key := make([]byte, 16)
	_, err = rand.Read(key)
	if err != nil {
		log.Fatal(err)
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

			log.Info("Verfiying decoded content")
			for i := range dataslice {
				if !bytes.Equal(dataslice[i], decryptedData[i]) {
					log.Fatal("Decoded output is different from the input.")
				}
			}
			if len(dataslice) != len(decryptedData) {
				log.Fatal("there are a different amount of decoded data and original data")
			}
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
	log.Infof("%v - %v: encoded entire slice of %v items, %v bytes each, in %v.", alg, mode, len(data), len(data[0]), time.Since(start))
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
	log.Infof("%v - %v: decoded entire slice of %v items, %v bytes each, in %v.", alg, mode, len(data), len(data[0]), time.Since(start))
	return decrypteddata, nil
}

func compress(data [][]byte, alg string) ([][]byte, error) {
	buf := bytes.NewBuffer(nil)
	log.Debug("Try to compress data with ", alg)
	var originalSize, newSize int
	compressor := compressors.NewCompressor(alg, buf)
	compressedData := make([][]byte, len(data))
	start := time.Now()
	for i := range data {
		originalSize += len(data[i])
		_, err := compressor.Write(data[i])
		if err != nil {
			log.Error("Failed to compress data: ", err)
			return nil, err
		}
		compressedData[i] = buf.Bytes()
		buf.Reset()
		newSize += len(compressedData[i])
	}
	log.Info("Compression done in ", time.Since(start))
	log.Infof("Compressed data with %v, original size: %v, new size: %v", alg, originalSize, newSize)
	return compressedData, nil
}

// func decrompress(data [][]byte, alg string) ([][]byte, error) {
//
// }
