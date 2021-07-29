package archive_extractor

import (
	"github.com/jfrog/go-archive-extractor/utils"
	"io"
	"os"
)

type Archiver interface {
	ExtractArchive(path string, processingFunc func(header *ArchiveHeader, params map[string]interface{}) error, params map[string]interface{}) error
}

type ArchiveHeader struct {
	ArchiveReader io.Reader
	IsFolder      bool
	Name          string
	ModTime       int64
	Size          int64
}

func NewArchiveHeader(archiveReader io.Reader, name string, modTime int64, size int64) *ArchiveHeader {
	return &ArchiveHeader{ArchiveReader: archiveReader, IsFolder: utils.IsFolder(name), Name: name, ModTime: modTime, Size: size}
}

func maxBytesLimit(path string, maxCompressRation int64) (int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size() * maxCompressRation, nil
}
