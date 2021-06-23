package aerrors

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	ar := New(fmt.Errorf("new error"))
	if ar.Error() != "Archive extractor error,new error" {
		t.Fatal("error do not match")
	}
}
