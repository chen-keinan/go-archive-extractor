package archive_extractor

import (
	"errors"
	"fmt"
	"github.com/blakesmith/ar"
	"io"
	"os"

	"github.com/go-archive-extractor/utils"
)

type DebArchvier struct {
}

func (za DebArchvier) ExtractArchive(archivePath string, advanceProcessing func(header *ArchiveHeader, advanceProcessingParams map[string]interface{}) error,
	advanceProcessingParams map[string]interface{}) error {
	debFile, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer debFile.Close()
	rc := ar.NewReader(debFile)
	if rc == nil {
		return errors.New(fmt.Sprintf("Failed to open deb file : %s", archivePath))
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
			return errors.New(fmt.Sprintf("Failed to open file : %s", archivePath))
		}
		if !utils.IsFolder(archiveEntry.Name) {
			archiveHeader := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.ModTime.Unix(), archiveEntry.Size)
			err = advanceProcessing(archiveHeader, advanceProcessingParams)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
