package perf

import (
	"crypto/md5"
	"crypto/rand"

	"docs.greenitglobe.com/despiegk/gotests/encryption_test/crypt"

	log "github.com/Sirupsen/logrus"
)

// RandomDataTest does encryption and compression tests on random data
func RandomDataTest(dataLength int) error {
	// data := make([]byte, dataLength)
	// read, err := rand.Read(data)
	data := []byte("My super secret message")
	log.Info("Message: ", string(data))
	md5a := md5.Sum(data)
	log.Debug(len(md5a), " ", md5a)
	// log.Debugf("Read %v bytes of random data", read)
	cipher, err := crypt.NewBlockCipher("twofish", md5a[:])
	if err != nil {
		log.Fatal(err)
	}

	iv := make([]byte, cipher.BlockSize())
	_, err = rand.Read(iv)
	if err != nil {
		log.Fatal(err)
	}

	dst := crypt.BlockEncrypt("ofb", cipher, data, iv)
	log.Warn("IV: ", iv)
	log.Warn("Original data: ", data)
	log.Warn("Encoded data: ", dst)

	dec := crypt.BlockDecrypt("ofb", cipher, dst, iv)

	log.Warn("Decoded data: ", dec)

	log.Warn("Decoded data with padding stripped: ", dec)
	log.Info("Decoded message: ", string(dec))
	return err
}

func RandomDataTestEverything(dataLength int) error {
	data := []byte("My super secret test stringy.")
	ciphers := []string{"aes" /*"3des", */, "twofish", "blowfish"}
	modes := []string{"cbc", "cfb", "ctr", "ofb"}

	key := make([]byte, 16)
	_, err := rand.Read(key)
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
			log.Info(c, ": ", mode)
			ciphertext := crypt.BlockEncrypt(mode, cipher, data, iv)
			log.Warn(ciphertext)

			dec := crypt.BlockDecrypt(mode, cipher, ciphertext, iv)
			log.Warn(dec)
			log.Warn(string(dec))
		}
	}
	return err
}
