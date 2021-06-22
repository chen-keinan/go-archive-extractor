package archive_extractor

import (
	"github.com/chen-keinan/go-archive-extractor/utils"
	"io"
)

type Archiver interface {
	Extract(path string) ([]ArchiveHeader, error)
}

type ArchiveHeader struct {
	Name    string
	ModTime int64
	Size    int64
	Sha1    string
	Sha2    string
	PkgMeta map[string]interface{}
}

func NewArchiveHeader(archiveReader io.Reader, name string, modTime int64, size int64) (*ArchiveHeader, error) {
	b, err := io.ReadAll(archiveReader)
	if err != nil {
		return nil, err
	}
	return &ArchiveHeader{Sha1: utils.NewSHA1(b), Sha2: utils.NewSHA2(b), Name: name, ModTime: modTime, Size: size}, nil
}
