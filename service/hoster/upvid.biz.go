package hoster

import (
	"deyforyou/dey/browser/html"
	"deyforyou/dey/browser"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
)

type UpvidBiz struct {
	url     string
	browser *browser.Browser
}

func NewUpvidBiz(value string) *UpvidBiz {
	return &UpvidBiz{
		url:     value,
		browser: browser.New(),
	}
}

func (hoster *UpvidBiz) Video() {
	browser := hoster.browser
	browser.Visit(hoster.url, func(element *html.Element) {
		firstIframeSrc := element.ChildAttribute("iframe", "src")
		browser.RequestByURL(firstIframeSrc, func(request *http.Request) {
			request.Header.Set("Referer", hoster.url)
		})

		browser.Visit(firstIframeSrc, func(element *html.Element) {
			secondIframeSrc := element.ChildAttribute("iframe", "src")
			clone := browser.Clone()
			clone.RequestByURL(secondIframeSrc, func(request *http.Request) {
				request.Header.Set("Referer", firstIframeSrc)
			})
			clone.ResponseByURL(secondIframeSrc, func(response *http.Response) {
				f, _ := os.Create("upvid.biz.html")
				io.Copy(f, response.Body)
			})
			clone.Visit(secondIframeSrc, func(element *html.Element) {
				data, _ := base64.StdEncoding.DecodeString(element.ChildAttribute("input#code", "value"))
				log.Println(string(data))
			})
		})
	})
}

func Decode(r, o string) string {
	e := make([]int, 0)
	for f := 0; f < 256; f++ {
		e = append(e, f)
	}

	var n int
	for f := 0; f < 256; f++ {
		n = (n + e[f] + int(r[f%len(r)])) % 256
		var t = e[f]
		e[f] = e[n]
		e[n] = t
	}

	var a string
	length := len(o)
	for i, n, f := 0, 0, 0; i < length; i++ {
		f = f + 1
		n = (n + e[f%256]) % 256
		var t = e[f]
		e[f] = e[n]
		e[n] = t
		a += string(rune(i ^ e[(e[f]+e[n])%256]))
	}
	return a
}
