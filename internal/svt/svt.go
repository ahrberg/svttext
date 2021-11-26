package svt

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	defaultBaseURL = "https://www.svt.se/text-tv"
	userAgent      = "gosvttext"
)

type Client struct {
	BaseUrl    string
	UserAgent  string
	HTTPClient *http.Client
}

func NewClient() *Client {

	c := &http.Client{
		Timeout: time.Second * 10,
	}

	return &Client{
		BaseUrl:    defaultBaseURL,
		UserAgent:  userAgent,
		HTTPClient: c,
	}
}

func (c *Client) GetNews(page string) (string, error) {
	body, err := c.getPage(page)

	if err != nil {
		return "", err
	}

	readable, err := parseReadable(body)

	if err != nil {
		return "", err
	}

	readable = cleanUp(readable)

	return readable, nil
}

func (c *Client) getPage(page string) (io.Reader, error) {
	url := fmt.Sprintf("%s/%s", c.BaseUrl, page)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func parseReadable(body io.Reader) (string, error) {

	// Find the content of div with class name
	const screenReaderClassMatch = "screenreaderOnly"

	doc, err := html.Parse(body)
	if err != nil {
		return "", nil
	}
	var f func(*html.Node) string

	f = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && strings.Contains(a.Val, screenReaderClassMatch) {
					return n.FirstChild.Data
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			readable := f(c)

			if readable != "" {
				return readable
			}
		}
		return ""
	}

	readable := f(doc)

	return readable, nil
}

func cleanUp(text string) string {
	ret := strings.ReplaceAll(text, "m ndag", "måndag")
	ret = strings.ReplaceAll(ret, "l rdag", "lördag")
	ret = strings.ReplaceAll(ret, "s ndag", "söndag")

	return ret
}
