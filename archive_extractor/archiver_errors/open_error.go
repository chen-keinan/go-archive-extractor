package archiver_errors

import "fmt"

type OpenError struct {
	msg string
	err error
}

func NewOpenError(msg string, err error) *OpenError {
	return &OpenError{err: err, msg: msg}
}

func (op OpenError) Error() string {
	return fmt.Sprintf("Failed to open file, file:%s, err:%s", op.msg, op.err.Error())
}

func (op OpenError) Unwrap() error {
	return op.err
}
