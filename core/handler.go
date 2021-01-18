package core

import (
	"github.com/gocolly/colly"
)

// HandlerFunc is
type HandlerFunc func(*colly.Collector)

// Handlers is
var Handlers map[string]HandlerFunc

func init() {
	Handlers = make(map[string]HandlerFunc)
}
