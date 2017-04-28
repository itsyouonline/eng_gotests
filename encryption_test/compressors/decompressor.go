package compressors

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/pierrec/lz4"
)

type Decompressor interface {
	io.ReadCloser
	Reset(r io.Reader) error
}

type lz4Decompressor struct {
	dec *lz4.Reader
}

func (lz lz4Decompressor) Close() error {
	return nil
}

func (lz lz4Decompressor) Reset(r io.Reader) error {
	lz.dec.Reset(r)
	return nil
}

func (lz lz4Decompressor) Read(buf []byte) (int, error) {
	return lz.dec.Read(buf)
}

type zlibDecompressor struct {
	dec io.ReadCloser
}

func (zl *zlibDecompressor) Reset(r io.Reader) error {
	var err error
	zl.dec, err = zlib.NewReader(r)
	return err
}

func (zl *zlibDecompressor) Close() error {
	return zl.dec.Close()
}

func (zl *zlibDecompressor) Read(buf []byte) (int, error) {
	return zl.dec.Read(buf)
}

func NewDecompressor(name string, src io.Reader) (Decompressor, error) {
	switch strings.ToLower(name) {
	case "lz4":
		return lz4Decompressor{
			lz4.NewReader(src),
		}, nil
	case "zlib":
		dec, err := zlib.NewReader(src)
		return &zlibDecompressor{
			dec: dec,
		}, err
	case "gzip":
		return gzip.NewReader(src)
	default:
		log.Fatal("Unrecognized compression algorithm")
		return nil, nil
	}
}
