package crypt

import (
	"crypto/cipher"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var additionalData = []byte("My little stringy")

type Encrypter interface {
	// Encrypt encrypts src and returns the output
	Encrypt(src []byte) []byte
}

// BlockEncrypter represents an encrypter in block mode. Padding is handled automatically
type BlockEncrypter struct {
	enc cipher.BlockMode
}

func (be *BlockEncrypter) Encrypt(src []byte) []byte {
	src = AddPadding(src, be.enc.BlockSize())
	dst := make([]byte, len(src))
	be.enc.CryptBlocks(dst, src)
	return dst
}

type StreamEncrypter struct {
	enc cipher.Stream
}

func (se *StreamEncrypter) Encrypt(src []byte) []byte {
	dst := make([]byte, len(src))
	se.enc.XORKeyStream(dst, src)
	return dst
}

type AeadEncrypter struct {
	enc   cipher.AEAD
	nonce []byte
}

func (ae AeadEncrypter) Encrypt(src []byte) []byte {
	dst := make([]byte, 0)
	ret := ae.enc.Seal(dst, ae.nonce, src, additionalData)
	return ret
}

// NewEncrypter returns a new Encrypter based on a given cipher.Block. The IV is
// required already. In production, a new random iv should be generated for every
// encryption.
func NewEncrypter(mode string, alg cipher.Block, iv []byte) Encrypter {
	switch strings.ToLower(mode) {
	case "cbc":
		return &BlockEncrypter{
			enc: cipher.NewCBCEncrypter(alg, iv),
		}
	case "cfb":
		return &StreamEncrypter{
			enc: cipher.NewCFBEncrypter(alg, iv),
		}
	case "ctr":
		return &StreamEncrypter{
			enc: cipher.NewCTR(alg, iv),
		}
	case "ofb":
		return &StreamEncrypter{
			enc: cipher.NewOFB(alg, iv),
		}
	case "gcm":
		aead, err := cipher.NewGCM(alg)
		if err != nil {
			log.Fatal("Failed to create encrypter in gcm mode: ", err)
		}
		return &AeadEncrypter{
			enc:   aead,
			nonce: iv[:12],
		}
	default:
		log.Fatal("Unrecognized block mode")
	}
	return nil
}
