package compressors

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/pierrec/lz4"
)

func NewDecompressor(name string, src io.Reader) (io.Reader, error) {
	switch strings.ToLower(name) {
	case "lz4":
		return lz4.NewReader(src), nil
	case "zlib":
		return zlib.NewReader(src)
	case "gzip":
		return gzip.NewReader(src)
	default:
		log.Fatal("Unrecognized compression algorithm")
		return nil, nil
	}
}
