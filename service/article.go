package service

import (
	"context"
	"deyforyou/dey/crawler/bot"
	"deyforyou/dey/crawler/colorizer"
	"deyforyou/dey/schema"
	"net/http"
	"net/url"
	"regexp"
	"sync"
)

var articles = make([]*schema.Article, 0)
var expression = regexp.MustCompile(`(http(s|)://.*.mp4)`)

func init() {
	var browser = bot.NewBot()
	var serie = &Serie{
		Bot:            browser,
		url:            "https://www.enstream.cc/series/page-1.html",
		serieTitleTag:  ".short-content h4",
		serieListTag:   "#dle-content .shortstory-in",
		saisonListTag:  ".season-list .shortstory-in",
		episodeListTag: ".episode-list .saision_LI2",
	}
	go serie.Visit(
		serie.url,
		func(serieListDoc *bot.Element) { // Visiting SerieList

			serieListDoc.ForEach(
				serie.serieListTag,
				func(i int, serieElement *bot.Element) {
					serieHref := serieElement.ChildAttribute("a", "href")
					serieTitle := serieElement.ChildText(serie.serieTitleTag)
					serie.Visit(serieHref, func(serieDoc *bot.Element) { // Visiting SaisonList
						group := &sync.WaitGroup{}
						serieDoc.ForEach(
							serie.saisonListTag,
							func(i int, saisonElement *bot.Element) {
								i++
								image := saisonElement.ChildAttribute("img", "src")
								article := &schema.Article{
									Title: serieTitle,
									Index: uint32(i),
									Image: image,
								}
								saisonHref := saisonElement.ChildAttribute("a", "href")
								group.Add(1)
								go serie.Visit(saisonHref, func(saisonDoc *bot.Element) { // Visiting EpisodeList
									saisonDoc.ForEach(
										serie.episodeListTag,
										func(i int, episodeElement *bot.Element) {
											episodeHref := episodeElement.ChildAttribute("a", "href")
											i++
											articleEpisode := &schema.Article_Episode{
												Index: uint32(i),
												Title: article.Title,
												Image: article.Image,
											}
											serie.Visit(episodeHref, func(episodeDoc *bot.Element) {
												episodeDoc.ForEachWithBreak("li.streamer:not(:first-child) > div", func(i int, e *bot.Element) bool {
													link := parseUrl(serie.url, e.Attribute("data-url"))
													var video string

													serie.RequestByURL(link, func(request *http.Request) {
														request.Header.Set("Referer", episodeDoc.Request.URL.String())
													})
													serie.Visit(link,
														func(iframeDoc *bot.Element) {
															if videoUrl := expression.FindString(iframeDoc.HTML); videoUrl != "" {
																articleEpisode.Referer = iframeDoc.Request.URL.String()
																colorizer.Warn("Video", videoUrl)
																colorizer.Error("Referer", link)
																video = videoUrl
																articleEpisode.Video = videoUrl
															}
														})

													return video != ""
												})
											})
											article.Episodes = append(article.Episodes, articleEpisode)
										})
									group.Done()

								})
								articles = append(articles, article)
							})

						group.Wait()

					})
				})

		})
}

// ArticleServiceServer is
type ArticleServiceServer struct {
	articles []*schema.Article
	schema.UnimplementedArticleServiceServer
}

// NewArticleServiceServer is
func NewArticleServiceServer() *ArticleServiceServer {
	return &ArticleServiceServer{
		articles:                          make([]*schema.Article, 0),
		UnimplementedArticleServiceServer: schema.UnimplementedArticleServiceServer{},
	}
}

// ListArticles is
func (ass *ArticleServiceServer) ListArticles(
	context context.Context,
	request *schema.ListArticlesRequest,
) (*schema.ListArticlesResponse, error) {

	// browser.Visit(
	// 	"https://www.enstream.cc/series/page-1.html",
	// 	func(doc *bot.Element) {
	// 		doc.ForEach(
	// 			"#dle-content .shortstory-in",
	// 			func(i int, post *bot.Element) {
	// 				link := post.ChildAttribute("a", "href")
	// 				title := post.ChildText(".short-content h4")

	// 				browser.Visit(
	// 					link,
	// 					func(doc *bot.Element) {
	// 						doc.ForEach(
	// 							".season-list .shortstory-in",
	// 							func(i int, post *bot.Element) {
	// 								image := post.ChildAttribute("img", "src")
	// 								i++
	// 								article := &schema.Article{
	// 									Gender: fmt.Sprintf("Saison %v", i),
	// 									Title:  title,
	// 									Image:  image,
	// 								}
	// 								articles = append(articles, article)
	// 							})
	// 					})
	// 			})
	// 	})

	return &schema.ListArticlesResponse{Articles: articles}, nil
}

type Video interface {
	article()
	articleList()
}

type Film struct {
	film     string
	filmList string
	*bot.Bot
}

type Serie struct {
	url             string
	serieListTag    string
	serieTitleTag   string
	saisonListTag   string
	saisonTitleTag  string
	episodeListTag  string
	episodeTitleTag string
	*bot.Bot
}

func (serie *Serie) aricleList() (articles []*schema.Article) {

	return
}

func (Film) aricleList(query string) {

}

func (Serie) aricle() {

}

func (Film) aricle() {

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
