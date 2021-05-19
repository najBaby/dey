package hoster

import "net/url"

const (
	upvidBiz = "upvid.biz"
	mystream = "embed.mystream.to"
)

type Hoster interface {
	Video()
}

type empty struct{}

func (*empty) Video() {

}

func New(value string) Hoster {
	url, _ := url.Parse(value)
	switch url.Host {
	case upvidBiz:
		return NewUpvidBiz(value)
	case mystream:
		return NewMyStream(value)
	}
	return &empty{}
}
