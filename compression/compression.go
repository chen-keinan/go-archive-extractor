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
	//BZ2 const
	BZ2 = ".bz2"
	//GZ const
	GZ = ".gz"
	//TGZ const
	TGZ = ".tgz"
	//LZW const
	LZW = ".Z"
	//INFL const
	INFL = ".infl"
	//Zlibe const
	Zlibe = ".xp3"
	//Xz const
	Xz = ".xz"

	//Here comes the magic
	gzipID1  = 0x1f
	gzipID2  = 0x8b
	bzip2ID1 = 0x42
	bzip2ID2 = 0x5a
)

var lZMAAloneMagic = []byte{0x5d, 0x00, 0x00}
var xzMagic = []byte{0xfd, 0x37, 0x7a}

//Compression interface
type Compression interface {
	GetReader(reader io.Reader) (io.ReadCloser, error)
}

//CreateCompression create new compression object
//accept compressed file path
//return compression object
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

//CreateCompressionFromBytes create compression object from bytes
// accept compressed file byte
// return compression object
func CreateCompressionFromBytes(magicBytes []byte) Compression {
	//TODO: replace with switch
	if magicBytes[0] == gzipID1 && magicBytes[1] == gzipID2 {
		return new(Gzip)
	}
	if magicBytes[0] == bzip2ID1 && magicBytes[1] == bzip2ID2 {
		return new(Bzip2)
	}
	if bytes.Equal(magicBytes[:3], xzMagic) {
		return new(XZ)
	}
	if bytes.Equal(magicBytes[:3], lZMAAloneMagic) {
		return new(Lzw)
	}
	return new(NoCompression)
}

//Bzip2 object
type Bzip2 struct {
}

//GetReader return bzip2 reader
//accept io.reader
func (comp Bzip2) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader := bzip2.NewReader(reader)
	cReader := ioutil.NopCloser(fileReader)
	return cReader, nil
}

//Flate object
type Flate struct {
}

//GetReader return flate reader
//accept io.reader
func (comp Flate) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader := flate.NewReader(reader)
	return fileReader, nil
}

//Gzip object
type Gzip struct {
}

//GetReader return gzip reader
//accept io.reader
func (comp Gzip) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return fileReader, nil
}

//Lzw object
type Lzw struct {
}

//GetReader return lzw reader
//accept io.reader
func (comp Lzw) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader := lzw.NewReader(reader, lzw.LSB, 100)
	return fileReader, nil
}

//Zlib object
type Zlib struct {
}

//GetReader return zlib reader
//accept io.reader
func (comp Zlib) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return fileReader, nil
}

//NoCompression object
type NoCompression struct {
}

//GetReader return NoCompression reader
//accept io.reader
func (comp NoCompression) GetReader(reader io.Reader) (io.ReadCloser, error) {
	return nil, nil
}

//XZ object
type XZ struct {
}

//GetReader return XZ reader
//accept io.reader
func (comp XZ) GetReader(reader io.Reader) (io.ReadCloser, error) {
	fileReader, err := xz.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return XZReaderCloser{fileReader}, nil
}

//XZReaderCloser object
type XZReaderCloser struct {
	*xz.Reader
}

//Close close XZReaderCloser
func (xzrc XZReaderCloser) Close() error {
	return nil
}
