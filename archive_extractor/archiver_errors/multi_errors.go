package archiver_errors

import (
	"bytes"
	"fmt"
)

type MultiError struct {
	Errors []error
}

func Append(err error, newError error) *MultiError {
	if merr, ok := err.(*MultiError); ok {
		if merr == nil {
			merr = &MultiError{}
		}
		merr.Errors = append(merr.Errors, newError)
		return merr
	} else {
		return &MultiError{Errors: []error{newError}}
	}
}

func (m *MultiError) Error() string {
	if m == nil || len(m.Errors) == 0 {
		return ""
	}

	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%d errors:\n", len(m.Errors)))
	for i, e := range m.Errors {
		buf.WriteString(fmt.Sprintf("%d.%v\n", i+1, e))
	}
	return buf.String()
}
