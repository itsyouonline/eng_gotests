package compression

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/pierrec/lz4"
)

// Lz4Decompress decompress data in the lz4 format
func Lz4Decompress(data []byte) ([]byte, error) {
	// outputBuffer will hold the decompressed data
	outputBuffer := bytes.NewBuffer(nil)
	// wrap the input data in an *io.Reader
	inputReader := bytes.NewReader(data)
	// get a new lz4 decompressor
	decomp := lz4.NewReader(inputReader)
	// decompress
	_, err := io.Copy(outputBuffer, decomp)
	return outputBuffer.Bytes(), err
}

// GzipDecompress decompresses data in the gzip format
func GzipDecompress(data []byte) ([]byte, error) {
	// outputBuffer will hold the decompressed data
	outputBuffer := bytes.NewBuffer(nil)
	// wrap the input data in an *io.Reader
	inputReader := bytes.NewReader(data)
	// get a new gzip decompressor
	decomp, err := gzip.NewReader(inputReader)
	if err != nil {
		return outputBuffer.Bytes(), err
	}
	// decompress
	_, err = io.Copy(outputBuffer, decomp)
	if err != nil {
		return outputBuffer.Bytes(), err
	}
	// flush and close
	err = decomp.Close()
	return outputBuffer.Bytes(), err
}
