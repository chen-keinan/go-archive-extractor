package archive_extractor

import (
	"github.com/JFrogDev/go-rpm"
	"github.com/chen-keinan/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/chen-keinan/go-archive-extractor/compression"
	"github.com/deoxxa/gocpio"
	"io"
	"os"
)

type RpmArchvier struct {
	RpmPackage rpm.Package
}

func (za *RpmArchvier) ExtractArchive(path string, processingFunc func(header *ArchiveHeader, params map[string]interface{}) error,
	params map[string]interface{}) error {
	rpm, err := rpm.OpenPackageFile(path)
	if err != nil {
		return archiver_errors.New(err)
	}
	za.RpmPackage = rpm
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	//Read content of cpio archive which starts after headers
	headerEnd := rpm.Headers[1].End
	archiveHead := make([]byte, 6)

	_, err = file.ReadAt(archiveHead, int64(headerEnd))
	if err != nil {
		return err
	}

	//rewind to start of the file
	file.Seek(int64(headerEnd), 0)
	fileReader, err := compression.CreateCompressionFromBytes(archiveHead).GetReader(file)
	defer fileReader.Close()
	if err != nil {
		return archiver_errors.New(err)
	}
	if fileReader == nil {
		return err
	}
	rc := cpio.NewReader(fileReader)
	var count = 0
	for {
		archiveEntry, err := rc.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		count++
		//skip trailer
		if archiveEntry.IsTrailer() {
			break
		}
		if archiveEntry.Mode != cpio.TYPE_DIR {
			if archiveEntry != nil {
				archiveHeader := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.Mtime, archiveEntry.Size)
				err = processingFunc(archiveHeader, params)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
