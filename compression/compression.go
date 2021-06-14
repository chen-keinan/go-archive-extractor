package compression

import (
	"bytes"
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"errors"
	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	bz2Ext   = ".bz2"
	tbz2Ext  = ".tbz2"
	gzExt    = ".gz"
	tgzExt   = ".tgz"
	lzwExt   = ".Z"
	inflExt  = ".infl"
	zlibExt  = ".xp3"
	xzExt    = ".xz"
	txzExt   = ".txz"
	lzmaExt  = ".lzma"
	tlzmaExt = ".tlzma"

	maxMagicBytes = 6 // 6 is the biggest used here (xz)
)

var (
	gzipMagic = []byte{0x1F, 0x8B}
	bz2Magic  = []byte{0x42, 0x5A, 0x68}
	xzMagic   = []byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00}
	lzmaMagic = []byte{0x5D, 0x00, 0x00}
)

func NewReaderSkipBytes(filePath string, skip int64) (io.ReadCloser, error) {
	return newReader(&fileArgs{path: filePath, skipBytes: skip})
}

func NewReader(filePath string) (io.ReadCloser, error) {
	return newReader(&fileArgs{path: filePath})
}

type fileArgs struct {
	path      string
	skipBytes int64
}

func (fa *fileArgs) open() (*os.File, error) {
	f, err := os.Open(fa.path)
	if err != nil {
		return nil, err
	}
	if fa.skipBytes > 0 {
		_, err := f.Seek(fa.skipBytes, io.SeekStart)
		if err != nil {
			f.Close()
			return nil, err
		}
	}
	return f, nil
}

func newReader(fa *fileArgs) (io.ReadCloser, error) {
	ext := filepath.Ext(fa.path)
	//these types has no defined magic bytes
	switch ext {
	case lzwExt:
		return initReader(fa, lzwReader)
	case inflExt:
		return initReader(fa, flateReader)
	case zlibExt:
		return initReader(fa, zlibReader)
	}
	// if possible init by magic bytes
	if magic, err := getMagicBytes(fa); err == nil {
		switch {
		case bytes.HasPrefix(magic, bz2Magic):
			return initReader(fa, bz2Reader)
		case bytes.HasPrefix(magic, gzipMagic):
			return initReader(fa, gzipReader)
		case bytes.HasPrefix(magic, xzMagic):
			return initReader(fa, xzReader)
		case bytes.HasPrefix(magic, lzmaMagic):
			return initReader(fa, lzmaReader)
		}
	}
	// fallback to init by extension
	switch ext {
	case bz2Ext, tbz2Ext:
		return initReader(fa, bz2Reader)
	case gzExt, tgzExt:
		return initReader(fa, gzipReader)
	case xzExt, txzExt:
		return initReader(fa, xzReader)
	case lzmaExt, tlzmaExt:
		return initReader(fa, lzmaReader)
	default:
		// no compression format found
		return initReader(fa, fileReader)
	}
}

func getMagicBytes(fa *fileArgs) ([]byte, error) {
	f, err := fa.open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b := make([]byte, maxMagicBytes)
	if _, err = f.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// compression readers

type cReader struct {
	reader io.ReadCloser
	file   *os.File
}

func (cr *cReader) Read(p []byte) (int, error) {
	return cr.reader.Read(p)
}

func (cr *cReader) Close() error {
	if err := cr.file.Close(); err != nil {
		return err
	}
	if err := cr.reader.Close(); err != nil {
		return err
	}
	return nil
}

func initReader(fa *fileArgs, getReader func(io.Reader) (io.ReadCloser, error)) (io.ReadCloser, error) {

	f, err := fa.open()
	if err != nil {
		return nil, err
	}
	r, err := getReader(f)
	if err != nil {
		f.Close()
		return nil, &ErrGetReader{err}
	}

	return &cReader{reader: r, file: f}, nil
}

type ErrGetReader struct {
	err error
}

func (e *ErrGetReader) Error() string {
	return e.err.Error()
}

func IsGetReaderError(err error) bool {
	for e := err; e != nil; e = errors.Unwrap(e) {
		if _, ok := e.(*ErrGetReader); ok {
			return true
		}
	}
	return false
}

func bz2Reader(reader io.Reader) (io.ReadCloser, error) {
	return ioutil.NopCloser(bzip2.NewReader(reader)), nil
}

func flateReader(reader io.Reader) (io.ReadCloser, error) {
	return flate.NewReader(reader), nil
}

func gzipReader(reader io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(reader)
}

func lzwReader(reader io.Reader) (io.ReadCloser, error) {
	return lzw.NewReader(reader, lzw.LSB, 100), nil
}

func zlibReader(reader io.Reader) (io.ReadCloser, error) {
	return zlib.NewReader(reader)
}

func xzReader(reader io.Reader) (io.ReadCloser, error) {
	r, err := xz.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(r), nil
}

func lzmaReader(reader io.Reader) (io.ReadCloser, error) {
	r, err := lzma.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(r), nil
}

func fileReader(reader io.Reader) (io.ReadCloser, error) {
	return ioutil.NopCloser(reader), nil
}
