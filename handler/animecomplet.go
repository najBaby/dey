package handler

import (
	"github.com/gocolly/colly"
)

// AnimeCompletURL is
const AnimeCompletURL = "https://animecomplet.co/"

// AnimeCompletFunc is
func AnimeCompletFunc(c *colly.Collector) {

	c.Visit(AnimeCompletURL)
}

func homePage() {
	
}
