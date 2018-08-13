package archive_extractor

import (
	"errors"
	"fmt"
	"github.com/blakesmith/ar"
	"io"
	"os"

	"github.com/chen-keinan/go-archive-extractor/utils"
)

type DebArchvier struct {
}

func (za DebArchvier) ExtractArchive(path string, processingFunc func(header *ArchiveHeader, params map[string]interface{}) error,
	params map[string]interface{}) error {
	debFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer debFile.Close()
	rc := ar.NewReader(debFile)
	if rc == nil {
		return errors.New(fmt.Sprintf("Failed to open deb file : %s", path))
	}
	for {
		archiveEntry, err := rc.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if archiveEntry == nil {
			return errors.New(fmt.Sprintf("Failed to open file : %s", path))
		}
		if !utils.IsFolder(archiveEntry.Name) {
			archiveHeader := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.ModTime.Unix(), archiveEntry.Size)
			err = processingFunc(archiveHeader, params)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
