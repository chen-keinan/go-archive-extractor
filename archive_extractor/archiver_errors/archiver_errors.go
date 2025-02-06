package archiver_errors

import (
	"fmt"
)

type ArchiverExtractorError struct {
	err error
	msg string
}

func New(e error) *ArchiverExtractorError {
	return &ArchiverExtractorError{err: e}
}

func NewArchiverExtractorError(msg string, e error) *ArchiverExtractorError {
	return &ArchiverExtractorError{err: e, msg: msg}
}

func (aee ArchiverExtractorError) Error() string {
	if aee.msg != "" {
		return fmt.Sprintf("Archive extractor error, message:%s, err:%s", aee.msg, aee.err.Error())
	}
	return fmt.Sprintf("Archive extractor error, %s", aee.err.Error())
}

func (aee ArchiverExtractorError) Unwrap() error {
	return aee.err
}
