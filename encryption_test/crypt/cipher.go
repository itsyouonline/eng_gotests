package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"strings"

	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/twofish"
)

func NewBlockCipher(block string, key []byte) (cipher.Block, error) {
	switch strings.ToLower(block) {
	case "aes":
		return aes.NewCipher(key)
	case "3des":
		// since 3DES requires a 24 byte key and our use case supplies 16 bit keys, just repeat it a bit
		// this isn't secure, only use this in tests
		newkey := make([]byte, 24)
		for i := range newkey {
			newkey[i] = key[i%len(key)]
		}
		return des.NewTripleDESCipher(newkey)
	case "blowfish":
		return blowfish.NewCipher(key)
	case "twofish":
		return twofish.NewCipher(key)
	default:
		return nil, fmt.Errorf("Cipher %v not recognized", block)
	}
}
