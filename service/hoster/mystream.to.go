package hoster

import (
	"deyforyou/dey/browser"
	"io"
	"net/http"
	"os"
)

type MyStream struct {
	url     string
	browser *browser.Browser
}

func NewMyStream(value string) *MyStream {
	return &MyStream{
		url:     value,
		browser: browser.New(),
	}
}

func (hoster *MyStream) Video() {
	browser := hoster.browser
	browser.ResponseByURL(hoster.url, func(response *http.Response) {
		f, _ := os.Create("mystream.html")
		io.Copy(f, response.Body)
	})
	browser.Visit(hoster.url, nil)
}
