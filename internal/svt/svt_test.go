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
	cases := []struct {
		text     string
		expected string
	}{
		{"hej\n 102 \n", "hej\n                   102 \n"},
		{"hej\n 102-103 \n", "hej\n                   102-103 \n"},
		{"hej\n1033\n 1012-103 \n", "hej\n1033\n 1012-103 \n"},
	}
	for _, tc := range cases {
		res := centerPageNr(tc.text)

		if res != tc.expected {
			t.Errorf("Test page number test failed, expected:`%s`, got:`%s`", tc.expected, res)
		}
	}
}
