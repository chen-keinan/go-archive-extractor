package archiver_errors

import (
	"fmt"
)

type ArchiverExtractorError struct {
	archiverError error
	message       string
}

func New(e error) ArchiverExtractorError {
	return ArchiverExtractorError{archiverError: e}
}

func NewWithMessage(msg string, e error) ArchiverExtractorError {
	return ArchiverExtractorError{archiverError: e, message: msg}
}

func (aee ArchiverExtractorError) Error() string {
	if aee.message != "" {
		return fmt.Sprintf("Archive extractor error, message:%s, err:%s", aee.message, aee.archiverError.Error())
	}
	return fmt.Sprintf("Archive extractor error, %s", aee.archiverError.Error())
}
