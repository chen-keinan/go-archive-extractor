package archive_extractor

import (
	"github.com/JFrogDev/go-rpm"
	"github.com/deoxxa/gocpio"
	"github.com/go-archive-extractor/compression"
	"io"
	"os"
)

type RpmArchvier struct {
}

func (za RpmArchvier) ExtractArchive(archivePath string, advanceProcessing func(header *ArchiveHeader, advanceProcessingParams map[string]interface{}) error,
	advanceProcessingParams map[string]interface{}) error {
	rpm, err := rpm.OpenPackageFile(archivePath)
	if err != nil {
		return nil
	}
	file, err := os.Open(archivePath)
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
		return nil
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
				err = advanceProcessing(archiveHeader, advanceProcessingParams)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
