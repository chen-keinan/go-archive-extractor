package compression

import (
	"bytes"
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"github.com/ulikunitz/xz"
	"io"
	"io/ioutil"
	"path/filepath"
)

const (
	BZ2   = ".bz2"
	GZ    = ".gz"
	TGZ   = ".tgz"
	LZW   = ".Z"
	INFL  = ".infl"
	Zlibe = ".xp3"
	Xz    = ".xz"

	//Here comes the magic
	gzipID1  = 0x1f
	gzipID2  = 0x8b
	bzip2ID1 = 0x42
	bzip2ID2 = 0x5a
)

var LZMA_alone_magic = []byte{0x5d, 0x00, 0x00}
var XZ_magic = []byte{0xfd, 0x37, 0x7a}

type Compression interface {
	GetReader(reader io.Reader) (io.ReadCloser, error)
}

func CreateCompression(fileName string) Compression {
	ext := filepath.Ext(fileName)
	switch ext {
	case BZ2:
		return new(Bzip2)
	case GZ, TGZ:
		return new(Gzip)
	case LZW:
		return new(Lzw)
	case INFL:
		return new(Flate)
	case Zlibe:
		return new(Zlib)
	case Xz:
		return new(XZ)
	default:
		return new(NoCompression)
	}
}

func CreateCompressionFromBytes(magicBytes []byte) Compression {
	//TODO: replace with switch
	if magicBytes[0] == gzipID1 && magicBytes[1] == gzipID2 {
		return new(Gzip)
	}
	if magicBytes[0] == bzip2ID1 && magicBytes[1] == bzip2ID2 {
		return new(Bzip2)
	}
	if bytes.Equal(magicBytes[:3], XZ_magic) {
		return new(XZ)
	}
	if bytes.Equal(magicBytes[:3], LZMA_alone_magic) {
	}
	return new(NoCompression)
}

type Bzip2 struct {
}

func (comp Bzip2) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader := bzip2.NewReader(reader)
	cReader := ioutil.NopCloser(fileReader)
	return cReader, nil
}

type Flate struct {
}

func (comp Flate) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader := flate.NewReader(reader)
	return fileReader, nil
}

type Gzip struct {
}

func (comp Gzip) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return fileReader, nil
}

type Lzw struct {
}

func (comp Lzw) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader := lzw.NewReader(reader, lzw.LSB, 100)
	return fileReader, nil
}

type Zlib struct {
}

func (comp Zlib) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return fileReader, nil
}

type NoCompression struct {
}

func (comp NoCompression) GetReader(reader io.Reader) (io.ReadCloser, error) {
	return nil, nil
}

type XZ struct {
}

func (comp XZ) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader, err := xz.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return XZReaderCloser{fileReader}, nil
}

type XZReaderCloser struct {
	*xz.Reader
}

func (xzrc XZReaderCloser) Close() error {
	return nil
}
