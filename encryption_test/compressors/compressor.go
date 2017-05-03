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
	Flush() error
}

// NewCompressor returns a new compressor for the given algorithm
func NewCompressor(name string, dst io.Writer) Compressor {
	switch strings.ToLower(name) {
	case "lz4":
		w := lz4.NewWriter(dst)
		// For some obscure reason, we have to enable block dependancy. Failure to do so
		// causes the lz4 library to write corrupted frames when compressing a capnp list
		// of sufficient size. Ofcourse, this behaviour isn't documented anywhere...
		w.Header.BlockDependency = true
		return w
	case "zlib":
		return zlib.NewWriter(dst)
	case "gzip":
		return gzip.NewWriter(dst)
	default:
		log.Fatal("Unrecognized compression algorithm")
		return nil
	}
}
