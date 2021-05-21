package service

import (
	"deyforyou/dey/browser"
	"deyforyou/dey/browser/html"
	"deyforyou/dey/schema"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	// var serie = &Serie{
	// 	Browser:        browser.New(),
	// 	hosterTitleTag: ".streamer",
	// 	serieTitleTag:  ".short-content h4",
	// 	episodeListTag: ".episode-list .saision_LI2",
	// 	serieListTag:   "#dle-content .shortstory-in",
	// 	saisonListTag:  ".season-list .shortstory-in",
	// 	hosterListTag:  ".clearfix ul:first-child .streamer",
	// 	url:            "https://www.enstream.cc/series/page-1.html",
	// }
	// serie.Request(func(request *http.Request) {
	// 	fmt.Println(request.URL)
	// })
	// serie.Visit(serie.url, serie.pageSerieList)
}

type Enstream struct {
	*browser.Browser
	SeriesUrl string
	BaseUrl   string
}

func NewEnstream() *Enstream {
	return &Enstream{
		SeriesUrl: "https://www.enstream.cc/series/action/page-%v.html",
		BaseUrl:   "https://w3.Enstream.net/",
		Browser:   browser.New(),
	}
}

func (site *Enstream) NewsMovies(offset int64) []*schema.Movie {
	movies := make([]*schema.Movie, 0)
	site.Visit(
		fmt.Sprintf(site.SeriesUrl, offset),
		func(element *html.Element) {
			movies = site.ListMovies(element)
		})
	return movies
}

func (site *Enstream) Search(query string) []*schema.Movie {
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

func (site *Enstream) ListMovies(element *html.Element) []*schema.Movie {
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

func (site *Enstream) Movie(url string) *schema.Movie {
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
		movie.Synopsis = strings.TrimSpace(strings.TrimPrefix(synopsis, "Synopsis"))
		<-ready

		// if episodes := movie.GetVideo(); len(episodes) > 0 {
		result = movie
		// }
	})
	return result
}
