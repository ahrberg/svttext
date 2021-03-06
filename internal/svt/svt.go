package svt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://www.svt.se/text-tv/api"
	userAgent      = "gosvttext"
)

type Client struct {
	BaseUrl    string
	UserAgent  string
	HTTPClient *http.Client
}

type svtTextRes struct {
	Status string `json:"status"`
	Data   struct {
		PageNumber string `json:"pageNumber"`
		PrevPage   string `json:"prevPage"`
		NextPage   string `json:"nextPage"`
		SubPages   []struct {
			SubPageNumber string `json:"subPageNumber"`
			GifAsBase64   string `json:"gifAsBase64"`
			ImageMap      string `json:"imageMap"`
			AltText       string `json:"altText"`
		} `json:"subPages"`
		Meta struct {
			Updated time.Time `json:"updated"`
		} `json:"meta"`
	} `json:"data"`
}

type SvtPage struct {
	Text       []string
	PageNumber string
	PrevPage   string
	NextPage   string
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

func (c *Client) GetNews(page string) (*SvtPage, error) {
	news, err := c.getPage(page)

	if err != nil {
		return nil, err
	}

	if news.Status != "success" {
		return &SvtPage{
			Text:       []string{},
			PageNumber: news.Data.PageNumber,
			PrevPage:   news.Data.PrevPage,
			NextPage:   news.Data.NextPage,
		}, nil
	}

	var texts = []string{}
	for _, page := range news.Data.SubPages {
		cleanText := cleanUp(page.AltText)
		texts = append(texts, cleanText)
	}

	return &SvtPage{
		Text:       texts,
		PageNumber: news.Data.PageNumber,
		PrevPage:   news.Data.PrevPage,
		NextPage:   news.Data.NextPage,
	}, nil
}

func (c *Client) getPage(page string) (*svtTextRes, error) {
	url := fmt.Sprintf("%s/%s", c.BaseUrl, page)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var s = new(svtTextRes)
	err = json.Unmarshal(body, &s)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func cleanUp(text string) string {
	ret := strings.ReplaceAll(text, "m ndag", "m??ndag")
	ret = strings.ReplaceAll(ret, "l rdag", "l??rdag")
	ret = strings.ReplaceAll(ret, "s ndag", "s??ndag")
	ret = centerPageNr(ret)

	return ret
}

// centerPageNr centers page numbers if the page number is the only
// charachters on the line
//
// can be of shape " 101 " or " 101-102 "
func centerPageNr(text string) string {
	lines := strings.Split(text, "\n")
	res := ""
	for i := 0; i < len(lines)-1; i++ {
		l := lines[i]
		match, err := regexp.Match(`(^ \d{3} $)|(^ \d{3}-\d{3} $)`, []byte(l))

		if err != nil {
			res += l + "\n"
			continue
		}

		if match != true {
			res += l + "\n"
			continue
		}

		res += "                  " + l + "\n"
	}
	return res
}
