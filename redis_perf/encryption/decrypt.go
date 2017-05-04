package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
)

// Decrypt decrypts a ciphertext with the given key. It expects the nonce to be prepended
// to the ciphertext
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	alg, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypter, err := cipher.NewGCM(alg)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, 0)
	// get the nonce we prepended to the ciphertext
	nonce := ciphertext[:decrypter.NonceSize()]
	plaintext, err := decrypter.Open(dst, nonce, ciphertext[decrypter.NonceSize():], additionalData)
	return plaintext, err
}
