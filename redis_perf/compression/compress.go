package compression

import (
	"bytes"
	"io"

	"github.com/pierrec/lz4"
)

// Compress compress data in the lz4 format
func Compress(data []byte) ([]byte, error) {
	outputBuffer := bytes.NewBuffer(nil)
	inputReader := bytes.NewReader(data)
	comp := lz4.NewWriter(outputBuffer)
	comp.Header.BlockDependency = true
	_, err := io.Copy(comp, inputReader)
	return outputBuffer.Bytes(), err
}
