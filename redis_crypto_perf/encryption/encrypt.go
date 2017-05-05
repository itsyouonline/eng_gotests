package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

var additionalData = []byte("something unrelated")

// Encrypt encrypts the given plaintext with the provided key. The nonce is randomly
// generated every run, and is prepended to the returned ciphertext.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	alg, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypter, err := cipher.NewGCM(alg)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, 0)
	nonce := make([]byte, encrypter.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	ciphertext := encrypter.Seal(dst, nonce, plaintext, additionalData)
	return append(nonce, ciphertext...), nil
}
