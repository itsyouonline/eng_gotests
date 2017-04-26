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
		return des.NewTripleDESCipher(key)
	case "blowfish":
		return blowfish.NewCipher(key)
	case "twofish":
		return twofish.NewCipher(key)
	default:
		return nil, fmt.Errorf("Cipher %v not recognized", block)
	}
}
