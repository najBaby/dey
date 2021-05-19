package service

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"deyforyou/dey/browser"
	"deyforyou/dey/browser/html"
	"deyforyou/dey/schema"

	"github.com/PuerkitoBio/goquery"
)

var expVideo = regexp.MustCompile(`(http(s|)://.*.mp4)`)

type StreamComplet3 struct {
	*browser.Browser
	FilmsUrl string
	BaseUrl  string
}

func NewStreamComplet3() *StreamComplet3 {
	return &StreamComplet3{
		FilmsUrl: "https://w3.streamcomplet3.net/film/page/%v/",
		BaseUrl:  "https://w3.streamcomplet3.net/",
		Browser:  browser.New(),
	}
}

func (site *StreamComplet3) NewsMovies(offset int64) []*schema.Movie {
	movies := make([]*schema.Movie, 0)
	site.Visit(
		fmt.Sprintf(site.FilmsUrl, offset),
		func(element *html.Element) {
			movies = site.ListMovies(element)
		})
	return movies
}

func (site *StreamComplet3) Search(query string) []*schema.Movie {
	movies := make([]*schema.Movie, 0)
	site.Form(
		site.BaseUrl,
		url.Values{
			"subaction": {"search"},
			"do":        {"search"},
			"story":     {query},
		}, func(element *html.Element) {
			movies = site.ListMovies(element)
		})
	return movies
}

func (site *StreamComplet3) ListMovies(element *html.Element) []*schema.Movie {
	movies := make([]*schema.Movie, 0)
	element.ForEach(
		"#dle-content .th-item",
		func(i int, e *html.Element) {
			image := parseUrl(site.BaseUrl, e.ChildAttribute(".th-img img", "src"))
			title := strings.TrimSuffix(e.ChildAttribute(".th-img img", "alt"), "Film Streaming Complet")
			source := parseUrl(site.BaseUrl, e.ChildAttribute(".th-in", "href"))

			movie := new(schema.Movie)
			movie.Title = title
			movie.Image = image
			movie.Source = source
			movie.Subtitle = "FILM"

			movies = append(movies, movie)
		})
	return movies
}

func (site *StreamComplet3) Movie(url string) *schema.Movie {
	result := new(schema.Movie)
	site.Visit(url, func(element *html.Element) {
		movie := new(schema.Movie)

		ready := make(chan bool)
		iframe := element.ChildAttribute(".iframe iframe", "src")
		site.RequestByURL(iframe, func(request *http.Request) {
			request.Header.Set("Referer", url)
		})
		site.ResponseByURL(iframe, func(response *http.Response) {
			url := response.Request.URL
			movie.Hoster = fmt.Sprintf("%s://%s", url.Scheme, url.Host)
		})
		go site.Visit(iframe, func(element *html.Element) {
			video := expVideo.FindString(element.Content)
			if video != "" {
				movie.Video = video
				// image := expImage.FindString(element.Content)
				// movie.Image = image
			}
			ready <- true
		})

		// title := element.ChildText("article .fone h1")
		synopsis := element.ChildText("article .full-text")
		image := parseUrl(url, element.ChildAttribute(".ftwo .fleft .fposter img", "src"))

		// movie.Title = title
		element.Selection.Find("article .short-info:nth-child(4)").
			Each(func(i int, s *goquery.Selection) {
				s.Find("a").Each(func(i int, s *goquery.Selection) {
					movie.Genres = append(movie.Genres, s.Text())
				})
			})
		movie.Image = image
		movie.Subtitle = "FILM"
		movie.Category = schema.Movie_FILM
		movie.Synopsis = strings.TrimSpace(synopsis)
		movie.Synopsis = element.Selection.Find("article .full-text").Contents().Text()

		<-ready

		// if episodes := movie.GetVideo(); len(episodes) > 0 {
		result = movie
		// }
	})
	return result
}
