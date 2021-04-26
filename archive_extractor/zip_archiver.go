package archive_extractor

import (
	"archive/zip"
	"fmt"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
)

type ZipArchiver struct{}

func (ZipArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {

	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer r.Close()
	var multiArchiveErr error
	for _, archiveEntry := range r.File {
		rc, err := archiveEntry.Open()
		if err != nil {
			if rc != nil {
				rc.Close()
			}
			multiArchiveErr = archiver_errors.Append(multiArchiveErr, fmt.Errorf("failed to open %s: %v", path, err))
			continue
		}
		archiveHeader := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.ModTime().Unix(), archiveEntry.FileInfo().Size())
		err = processingFunc(archiveHeader, params)
		if err != nil {
			if rc != nil {
				rc.Close()
			}
			return err
		}
		rc.Close()
	}
	if multiArchiveErr != nil {
		return archiver_errors.New(multiArchiveErr)
	}
	return nil
}
