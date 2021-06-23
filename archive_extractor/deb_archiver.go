package archive_extractor

import (
	"errors"
	"fmt"
	"github.com/blakesmith/ar"
	"io"
	"os"
	"path/filepath"

	"github.com/chen-keinan/go-archive-extractor/utils"
)

type DebArchvier struct {
}

//Extract extract deb archive
//accept deb file path
//return file header metadata
func (za DebArchvier) Extract(path string) ([]*ArchiveHeader, error) {
	headers := make([]*ArchiveHeader, 0)
	debFile, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = debFile.Close()
		if err != nil {
			fmt.Print(err.Error())
		}
	}()
	rc := ar.NewReader(debFile)
	if rc == nil {
		return nil, errors.New(fmt.Sprintf("Failed to open deb file : %s", path))
	}
	for {
		archiveEntry, err := rc.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if archiveEntry == nil {
			return nil, errors.New(fmt.Sprintf("Failed to open file : %s", path))
		}
		if !utils.IsFolder(archiveEntry.Name) {
			archiveHeader, err := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.ModTime.Unix(), archiveEntry.Size)
			if err != nil {
				return nil, err
			}
			headers = append(headers, archiveHeader)
		}
	}
	return headers, nil
}
