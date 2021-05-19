package browser

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"deyforyou/dey/browser/html"
	"deyforyou/dey/browser/remote"

	"github.com/PuerkitoBio/goquery"
)

// var (
// 	header = http.Header{
// 		"Connection":      {"keep-alive"},
// 		"Accept-Language": {"fr,fr-FR;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6"},
// 		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
// 		"User-Agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"},
// 	}
// )

type (
	ElementCallBack func(element *html.Element)

	Browser struct {
		*remote.Remote
		elementCallBacks []ElementCallBack
	}
)

func New() *Browser {
	jar, _ := cookiejar.New(nil)
	browser := &Browser{
		Remote: remote.New(&http.Client{
			Jar:       jar,
			Transport: &http.Transport{},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				fmt.Println(req.URL)
				return nil
			},
		}),
		elementCallBacks: make([]ElementCallBack, 0),
	}
	return browser
}

func (browser *Browser) Clone() *Browser {
	return &Browser{
		Remote:           browser.Remote.Clone(),
		elementCallBacks: browser.elementCallBacks,
	}
}

func (browser *Browser) Document(callBack ElementCallBack) {
	browser.elementCallBacks = append(browser.elementCallBacks, callBack)
}

func (browser *Browser) Visit(url string, callBack ElementCallBack) {
	log.Println(url)
	response, err := browser.Get(url)
	if err != nil {
		log.Println(err)
	} else {
		document, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			panic(err)
		}
		if callBack != nil {
			callBack(html.NewElement(document.Selection))
		}
	}
}

func (browser *Browser) Form(url string, values url.Values, callBack ElementCallBack) {
	log.Println(url)
	response, err := browser.PostForm(url, values)
	if err != nil {
		log.Println(err)
	} else {
		document, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			panic(err)
		}
		if callBack != nil {
			callBack(html.NewElement(document.Selection))
		}
	}
}
