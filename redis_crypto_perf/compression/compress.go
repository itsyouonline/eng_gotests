package compression

import (
	"bytes"
	"io"

	"github.com/pierrec/lz4"
)

// Compress compress data in the lz4 format
func Compress(data []byte) ([]byte, error) {
	// buffer for our output
	outputBuffer := bytes.NewBuffer(nil)
	// wrap the input in an *io.Reader
	inputReader := bytes.NewReader(data)
	// make or compressor with the underlying *io.Writer as output
	comp := lz4.NewWriter(outputBuffer)
	// Set block dependancy to true
	comp.Header.BlockDependency = true
	// compress the data
	_, err := io.Copy(comp, inputReader)
	if err != nil {
		return outputBuffer.Bytes(), err
	}
	// flush and close the compressor
	err = comp.Close()
	// return the valuo of our output buffer and a possible error
	return outputBuffer.Bytes(), err
}
