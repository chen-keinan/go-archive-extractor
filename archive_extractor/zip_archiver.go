package archive_extractor

import (
	"archive/zip"
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
	for _, archiveEntry := range r.File {
		rc, err := archiveEntry.Open()
		if err != nil {
			rc.Close()
			return err
		}
		archiveHeader := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.Modified.Unix(), archiveEntry.FileInfo().Size())
		err = processingFunc(archiveHeader, params)
		if err != nil {
			rc.Close()
			return err
		}
		rc.Close()
	}
	return nil
}
