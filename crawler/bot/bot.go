package bot

import (
	"deyforyou/dey/crawler/colorizer"
	"net/http"
)

var (
	header = http.Header{
		"Connection":      {"keep-alive"},
		"Accept-Language": {"fr,fr-FR;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
		"User-Agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"},
	}
)

type ElementCallBack func(element *Element)

type Bot struct {
	*Remote
	elementCallBacks []ElementCallBack
}

func NewBot() *Bot {
	bot := &Bot{
		Remote: New(),
	}
	return bot
}

func (bot *Bot) Document(callBack ElementCallBack) {
	bot.elementCallBacks = append(bot.elementCallBacks, callBack)
}

func (bot *Bot) Visit(url string, callBack ElementCallBack) {
	response, err := bot.Get(url)
	if err != nil {
		colorizer.Error("Error", err)
		return
	}
	callBack(NewElementFromResponse(response))
}
