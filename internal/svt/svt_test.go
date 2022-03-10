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

func TestCenterPageNr(t *testing.T) {
	// Arrange
	text := "hej\n 102 \n"
	expected := "hej\n                   102 \n"

	// Act
	res := centerPageNr(text)

	// Assert
	if res != expected {
		t.Errorf("Page number not centered, expected:`%s`, got:`%s`", expected, res)
	}

	// Arrange
	text = "hej\n 1043\n103\n 11\ntest\n"

	// Act
	res = centerPageNr(text)

	// Assert
	if res != text {
		t.Error("Result not expected to be modified")
	}
}
