package svt

import (
	"fmt"
	"testing"
)

func TestColorPageNumbers(t *testing.T) {

	expected := "This is a test " + fmt.Sprintf(string(boldBold), "100") + " this is a test"

	res := ColorPage("This is a test 100 this is a test")

	if res != expected {
		t.Error("Page Numbers not colorized?")
	}
}
