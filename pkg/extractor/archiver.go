package extractor

import (
	"github.com/chen-keinan/go-archive-extractor/utils"
	"io"
)

//Archive enum
type Archive int

const (
	//Zip const
	Zip Archive = iota
	//Tar const
	Tar
	//SevenZip const
	SevenZip
	//Deb const
	Deb
	//Rpm const
	Rpm
)

//Archiver interface
type Archiver interface {
	Extract(path string) ([]*ArchiveHeader, error)
}

//ArchiveHeader archive headers object
type ArchiveHeader struct {
	Name    string
	ModTime int64
	Size    int64
	Sha1    string
	Sha2    string
	PkgMeta map[string]interface{}
}

//NewArchiveHeader return new archiver header metadata object
// accept header data
// return headers metadata object
func NewArchiveHeader(archiveReader io.Reader, name string, modTime int64, size int64) (*ArchiveHeader, error) {
	b, err := io.ReadAll(archiveReader)
	if err != nil {
		return nil, err
	}
	return &ArchiveHeader{Sha1: utils.NewSHA1(b), Sha2: utils.NewSHA2(b), Name: name, ModTime: modTime, Size: size}, nil
}

//New instantiate new archive
func New(arc Archive) Archiver {
	switch arc {
	case Zip:
		return new(zipArchvier)
	case Rpm:
		return new(rpmArchvier)
	case Deb:
		return new(debArchvier)
	case SevenZip:
		return new(sevenZipArchvier)
	case Tar:
		return new(tarArchvier)
	default:
		return new(zipArchvier)
	}
}
