package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/nine-monsters/crawler/collect"
	"github.com/nine-monsters/crawler/proxy"
	"regexp"
	"time"
)

var (
	headerRe = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)
)

func main() {
	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	switcher, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		fmt.Println("RoundRobinProxySwitcher failed")
	}

	url := "https://google.com"
	fetch := collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   switcher,
	}
	body, err := fetch.Get(url)

	if err != nil {
		fmt.Printf("read content failed:%v", err)
		return
	}
	fmt.Println(string(body))

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("read content failed:%v", err)
	}

	doc.Find("div.news_li h2 a[target=_blank]").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}
