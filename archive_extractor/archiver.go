package archive_extractor

import (
	"io"
	"jfrog.com/xray/utils"
)

type Archiver interface {
	ExtractArchive(archivePath string, advanceProcessing func(header *ArchiveHeader) error) error
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
