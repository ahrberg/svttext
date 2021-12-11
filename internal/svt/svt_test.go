package svt

import (
	"testing"
)

func TestGetNews(t *testing.T) {
	client := NewClient()
	text, err := client.GetNews("100")

	if err != nil {
		t.Error("No error expected")
	}

	if len(text.Text) < 1 {
		t.Error("Could not find screen reader content")
	}
}

func TestGetNewsNotFound(t *testing.T) {
	client := NewClient()
	text, err := client.GetNews("this_should_not_exists")

	if err != nil {
		t.Error("No error expected")
	}

	if len(text.Text) > 0 {
		t.Error("Expecting empty text for non existing page")
	}
}
