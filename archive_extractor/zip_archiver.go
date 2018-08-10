package archive_extractor

import (
	"archive/zip"
)

type ZipArchvier struct {
}

func (za ZipArchvier) ExtractArchive(archivePath string, advanceProcessing func(header *ArchiveHeader, advanceProcessingParams map[string]interface{}) error,
	advanceProcessingParams map[string]interface{}) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		if err.Error() == "zip: not a valid zip file" {
			return nil
		}
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
		err = advanceProcessing(archiveHeader, advanceProcessingParams)
		if err != nil {
			rc.Close()
			return err
		}
		rc.Close()
	}
	return nil
}
