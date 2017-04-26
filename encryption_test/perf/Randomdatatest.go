package perf

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"

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
	aes, err := aes.NewCipher(md5a[:])
	if err != nil {
		log.Fatal(err)
	}
	iv := make([]byte, aes.BlockSize())
	_, err = rand.Read(iv)
	if err != nil {
		log.Fatal(err)
	}
	cbl := cipher.NewCBCEncrypter(aes, iv)

	//PKCS7 padding
	paddingsize := aes.BlockSize() - len(data)%aes.BlockSize()
	padding := make([]byte, paddingsize)
	// set all the padding bytes to the amount of bytes we add as per the PKCS7
	for i := range padding {
		padding[i] = byte(paddingsize)
	}
	data = append(data, padding...)
	dst := make([]byte, len(data))
	cbl.CryptBlocks(dst, data)
	log.Warn("IV: ", iv)
	log.Warn("Original data: ", data)
	log.Warn("Encoded data: ", dst)
	cbld := cipher.NewCBCDecrypter(aes, iv)
	dec := make([]byte, len(dst))
	cbld.CryptBlocks(dec, dst)
	log.Warn("Decoded data: ", dec)

	// strip padding
	// read the last byte of the slice
	paddingBytes := int(dec[len(dec)-1])
	padding = dec[len(dec)-paddingBytes:]
	log.Warn("Padding: ", padding)
	for _, v := range padding {
		if int(v) != paddingBytes {
			log.Fatalf("Wrong value for padding byte, got %v, expected %v", int(v), paddingBytes)
		}
	}
	dec = dec[:len(dec)-paddingBytes]
	log.Warn("Decoded data with padding stripped: ", dec)
	log.Info("Decoded message: ", string(dec))
	return err
}
