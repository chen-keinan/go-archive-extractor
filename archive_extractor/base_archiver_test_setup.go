package archive_extractor

import (
	"errors"
	"io"
)

type ArchiveData struct {
	ArchiveReader io.Reader
	IsFolder      bool
	Name          string
	ModTime       int64
	Size          int64
}

func params() map[string]interface{} {
	return map[string]interface{}{
		"archiveData": &ArchiveData{},
	}
}

func processingFunc(header *ArchiveHeader, params map[string]interface{}) error {
	if len(params) == 0 {
		return errors.New("Advance processing params are missing")
	}
	var ok bool
	var archiveData *ArchiveData
	if archiveData, ok = params["archiveData"].(*ArchiveData); !ok {
		return errors.New("Failed to read param")
	}
	archiveData.Name = header.Name
	archiveData.ModTime = header.ModTime
	archiveData.Size = header.Size
	archiveData.IsFolder = header.IsFolder
	archiveData.ArchiveReader = header.ArchiveReader
	return nil
}
