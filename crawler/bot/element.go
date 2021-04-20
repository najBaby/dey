package bot

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Element struct {
	Document *goquery.Selection
	Request  *http.Request
	HTML     string
}

func NewElementFromResponse(response *http.Response) *Element {
	document, _ := goquery.NewDocumentFromReader(response.Body)
	html, _ := document.Html()
	return &Element{
		HTML:     html,
		Request:  response.Request,
		Document: document.Selection,
	}
}

func NewElementFromNode(node *html.Node, request *http.Request) *Element {
	document := goquery.NewDocumentFromNode(node)
	html, _ := document.Html()
	return &Element{
		HTML:     html,
		Request:  request,
		Document: document.Selection,
	}
}

func (element *Element) ChildText(selector string) string {
	return strings.TrimSpace(element.Document.Find(selector).Text())
}

func (element *Element) ChildTexts(selector string) (values []string) {
	element.Document.Find(selector).Each(func(_ int, selection *goquery.Selection) {
		values = append(values, strings.TrimSpace(selection.Text()))
	})
	return values
}

func (element *Element) Attribute(k string) string {
	if value, ok := element.Document.Attr(k); ok {
		return value
	}
	return ""
}

func (element *Element) ChildAttribute(selector, name string) string {
	if attr, ok := element.Document.Find(selector).Attr(name); ok {
		return strings.TrimSpace(attr)
	}
	return ""
}

func (element *Element) ChildAttributes(selector, attrName string) (result []string) {
	element.Document.Find(selector).Each(func(_ int, s *goquery.Selection) {
		if attr, ok := s.Attr(attrName); ok {
			result = append(result, strings.TrimSpace(attr))
		}
	})
	return
}

func (element *Element) ForEach(selector string, callback func(int, *Element)) {
	element.Document.Find(selector).Each(func(_ int, selection *goquery.Selection) {
		for index, node := range selection.Nodes {
			callback(index, NewElementFromNode(node, element.Request))
		}
	})
}

func (element *Element) ForEachWithBreak(selector string, callback func(int, *Element) bool) {
	element.Document.Find(selector).EachWithBreak(func(_ int, selection *goquery.Selection) bool {
		for index, node := range selection.Nodes {
			if callback(index, NewElementFromNode(node, element.Request)) {
				return true
			}
		}
		return false
	})
}
