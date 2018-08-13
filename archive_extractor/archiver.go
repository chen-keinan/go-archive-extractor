package archive_extractor

import (
	"github.com/chen-keinan/go-archive-extractor/utils"
	"io"
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
