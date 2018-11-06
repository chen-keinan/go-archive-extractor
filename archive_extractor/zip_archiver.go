package archive_extractor

import (
	"archive/zip"
	"fmt"
	"github.com/chen-keinan/go-archive-extractor/archive_extractor/archiver_errors"
)

type ZipArchvier struct {
}

func (za ZipArchvier) ExtractArchive(path string, processingFunc func(header *ArchiveHeader, params map[string]interface{}) error,
	params map[string]interface{}) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer r.Close()
	var multiErr error
	for _, archiveEntry := range r.File {
		rc, err := archiveEntry.Open()
		if err != nil {
			if rc != nil {
				rc.Close()
			}
			multiErr = archiver_errors.Append(multiErr, fmt.Errorf("failed to open %s: %v", path, err))
			continue
		}
		archiveHeader := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.ModTime().Unix(), archiveEntry.FileInfo().Size())
		err = processingFunc(archiveHeader, params)
		if err != nil {
			if rc != nil {
				rc.Close()
			}
			multiErr = archiver_errors.Append(multiErr, fmt.Errorf("failed to process %s: %v", path, err))
			continue
		}
		rc.Close()
	}
	return multiErr
}
