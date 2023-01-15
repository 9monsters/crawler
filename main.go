package main

import (
	"bytes"
	"crawler/collect"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"time"
)

var (
	headerRe = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)
)

func main() {
	url := "https://baidu.com/"
	fetch := collect.BrowserFetch{
		Timeout: 300 * time.Millisecond,
	}
	body, err := fetch.Get(url)

	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}
	fmt.Println(string(body))

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Println("read content failed:%v", err)
	}

	doc.Find("div.news_li h2 a[target=_blank]").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}
