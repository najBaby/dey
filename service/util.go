package service

import (
	"deyforyou/dey/browser"
	"deyforyou/dey/browser/html"
	"deyforyou/dey/service/hoster"
	"net/http"
	"net/url"
)

type Serie struct {
	url             string
	serieListTag    string
	serieTitleTag   string
	saisonListTag   string
	saisonTitleTag  string
	episodeListTag  string
	episodeTitleTag string
	hosterListTag   string
	hosterTitleTag  string
	*browser.Browser
}

func (serie *Serie) Clone() *Serie {
	return &Serie{
		url:             serie.url,
		serieListTag:    serie.serieListTag,
		serieTitleTag:   serie.serieListTag,
		saisonListTag:   serie.serieListTag,
		saisonTitleTag:  serie.serieListTag,
		episodeListTag:  serie.serieListTag,
		episodeTitleTag: serie.serieListTag,
		hosterListTag:   serie.hosterListTag,
		hosterTitleTag:  serie.hosterTitleTag,
		Browser:         serie.Browser.Clone(),
	}
}

func (serie *Serie) pageSerieList(element *html.Element) {
	element.ForEach(
		serie.serieListTag,
		func(i int, serieElement *html.Element) {
			saisonListHref := serieElement.ChildAttribute("a", "href")
			serie.Visit(saisonListHref, serie.pageSaisonList)
		})
}

func (serie *Serie) pageSaisonList(element *html.Element) {
	element.ForEach(
		serie.saisonListTag,
		func(i int, saisonElement *html.Element) {
			episodeListHref := saisonElement.ChildAttribute("a", "href")
			go serie.Visit(episodeListHref, serie.pageEpisodeList)
		})
}

func (serie *Serie) pageEpisodeList(element *html.Element) {

	element.ForEach(
		serie.episodeListTag,
		func(i int, episodeElement *html.Element) {
			hosterListHref := episodeElement.ChildAttribute("a", "href")
			serie.Visit(hosterListHref, serie.pageHosterList)
		})
}

func (serie *Serie) pageHosterList(element *html.Element) {

	element.ForEach(
		serie.hosterListTag,
		func(i int, hosterElement *html.Element) {
			hosterUrl := parseUrl(serie.url, hosterElement.ChildAttribute("div", "data-url"))
			clone := serie.Clone()
			clone.ResponseByURL(hosterUrl, func(response *http.Response) {
				hoster.New(response.Request.URL.String()).Video()
			})
			clone.Visit(hosterUrl, nil)
		})
}

func parseUrl(dirUrl, baseUrl string) string {
	dir, err := url.Parse(dirUrl)
	if err != nil {
		panic(err)
	}
	base, err := url.Parse(baseUrl)
	if err != nil {
		panic(err)
	}
	dir.Path = base.Path
	return dir.String()
}
