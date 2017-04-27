package crypt

import (
	"crypto/cipher"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Decrypter interface {
	// Decrypt decrypts src and returns the output
	Decrypt(src []byte) []byte
}

// BlockDecrypter represents a decrypter in block mode. Padding is handled automatically
type BlockDecrypter struct {
	dec cipher.BlockMode
}

func (bd BlockDecrypter) Decrypt(src []byte) []byte {
	dst := make([]byte, len(src))
	bd.dec.CryptBlocks(dst, src)
	dst = RemovePadding(dst)
	return dst
}

type StreamDecrypter struct {
	dec cipher.Stream
}

func (sd StreamDecrypter) Decrypt(src []byte) []byte {
	dst := make([]byte, len(src))
	sd.dec.XORKeyStream(dst, src)
	return dst
}

type AeadDecrypter struct {
	dec   cipher.AEAD
	nonce []byte
}

func (ad AeadDecrypter) Decrypt(src []byte) []byte {
	dst := make([]byte, 0)
	ret, err := ad.dec.Open(dst, ad.nonce, src, additionalData)
	if err != nil {
		log.Fatal("Failed to decrypt message: ", err)
	}
	return ret
}

// NewDecrypter returns a new Decrypter based on a given cipher.Block. The IV is
// required already. In production, IV should be random for every block and derived
// from the ciphertext (assuming it is appended or prepended there)
func NewDecrypter(mode string, alg cipher.Block, iv []byte) Decrypter {
	switch strings.ToLower(mode) {
	case "cbc":
		return &BlockDecrypter{
			dec: cipher.NewCBCDecrypter(alg, iv),
		}
	case "cfb":
		return &StreamDecrypter{
			dec: cipher.NewCFBDecrypter(alg, iv),
		}
	case "ctr":
		return &StreamDecrypter{
			dec: cipher.NewCTR(alg, iv),
		}
	case "ofb":
		return &StreamDecrypter{
			dec: cipher.NewOFB(alg, iv),
		}
	case "gcm":
		aead, err := cipher.NewGCM(alg)
		if err != nil {
			log.Fatal("Failed to create decrypter in gcm mode: ", err)
		}
		return &AeadDecrypter{
			dec:   aead,
			nonce: iv[:12],
		}
	default:
		log.Fatal("Unrecognized block mode")
	}
	return nil
}

func BlockDecrypt(blockMode string, block cipher.Block, src, iv []byte) []byte {
	dst := make([]byte, len(src))
	switch strings.ToLower(blockMode) {
	case "cbc":
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(dst, src)
		dst = RemovePadding(dst)
		break
	case "cfb":
		cipher.NewCFBDecrypter(block, iv).XORKeyStream(dst, src)
		break
	case "ctr":
		cipher.NewCTR(block, iv).XORKeyStream(dst, src)
		break
	case "ofb":
		cipher.NewOFB(block, iv).XORKeyStream(dst, src)
		break
	default:
		log.Fatal("Unrecognized block mode")
	}
	return dst
}
