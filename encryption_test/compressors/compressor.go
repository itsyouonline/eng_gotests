package compressors

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/pierrec/lz4"
)

type Compressor interface {
	io.WriteCloser
	Reset(io.Writer)
}

// NewCompressor returns a new compressor for the given algorithm
func NewCompressor(name string, dst io.Writer) Compressor {
	switch strings.ToLower(name) {
	case "lz4":
		return lz4.NewWriter(dst)
	case "zlib":
		return zlib.NewWriter(dst)
	case "gzip":
		return gzip.NewWriter(dst)
	default:
		log.Fatal("Unrecognized compression algorithm")
		return nil
	}
}
