package compression

import (
	"bytes"
	"io"

	"github.com/pierrec/lz4"
)

// Decompress decompress data in the lz4 format
func Decompress(data []byte) ([]byte, error) {
	outputBuffer := bytes.NewBuffer(nil)
	inputReader := bytes.NewReader(data)
	decomp := lz4.NewReader(inputReader)
	_, err := io.Copy(outputBuffer, decomp)
	return outputBuffer.Bytes(), err
}
