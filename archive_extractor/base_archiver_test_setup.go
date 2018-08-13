package archive_extractor

import (
	"errors"
	"fmt"
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
		"archveData": &ArchiveData{},
	}
}

func processingFunc(header *ArchiveHeader, params map[string]interface{}) error {
	if len(params) == 0 {
		return errors.New("Advance processing params are missing")
	}
	var ok bool
	var archiveData *ArchiveData
	if archiveData, ok = params["archiveData"].(*ArchiveData); !ok {
		return errors.New("Advance processing archveData param is missing")
	}
	archiveData.Name = header.Name
	archiveData.ModTime = header.ModTime
	archiveData.Size = header.Size
	archiveData.IsFolder = header.IsFolder
	fmt.Print(archiveData)
	return nil
}
