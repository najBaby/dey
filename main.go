package main

import (
	"movies_scraping/core"
	"time"
)

func main() {
	core.Run(time.Minute * 20)
}
