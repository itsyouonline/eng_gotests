package crypt

import (
	"crypto/cipher"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func BlockEncrypt(blockMode string, block cipher.Block, src, iv []byte) []byte {
	dst := make([]byte, len(src))
	switch strings.ToLower(blockMode) {
	case "cbc":
		src = AddPadding(src, block.BlockSize())
		dst = make([]byte, len(src))
		cipher.NewCBCEncrypter(block, iv).CryptBlocks(dst, src)
		break
	case "cfb":
		cipher.NewCFBEncrypter(block, iv).XORKeyStream(dst, src)
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
