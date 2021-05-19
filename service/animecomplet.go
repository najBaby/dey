package service

import (
	"deyforyou/dey/browser"
	"deyforyou/dey/browser/html"
)

type AnimeComplet struct {
	*browser.Browser
}

type Season struct {
}

type Film struct {
}

func (ac *AnimeComplet) Film(url string)     {}
func (ac *AnimeComplet) Serie(url string)    {}
func (ac *AnimeComplet) Season(url string)   {}
func (ac *AnimeComplet) Episode(url string)  {}
func (ac *AnimeComplet) Episodes(url string) {}

func (ac *AnimeComplet) Films(query string) {}
func (ac *AnimeComplet) Series(query string) {
	ac.Visit("", ac.recentPosts)
}

func (ac *AnimeComplet) recentPosts(element *html.Element) {
	element.ForEach(
		".recent-posts li",
		func(i int, e *html.Element) {
			articleLink := e.ChildAttribute(".meta-category a", "href")
			ac.Visit(articleLink, ac.article)
		})
}
 
func (ac *AnimeComplet) article(element *html.Element) {
	// article := new(schema.Article)
	element.ForEach(
		".recent-posts li",
		func(i int, e *html.Element) {
			articleLink := e.ChildAttribute("h2 a", "href")
			ac.Visit(articleLink, func(element *html.Element) {
			})
		})
}

func (ac *AnimeComplet) Seasons(query string) {}
