package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// tag v0.0.5
var headerRe = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)

func main() {
	url := "https://www.thepaper.cn/"
	body, err := Fetch(url)

	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}

	numLinks := strings.Count(string(body), "<a")
	fmt.Printf("homepage has %d links!\n", numLinks)

	exist := strings.Contains(string(body), "疫情")
	fmt.Printf("是否存在疫情:%v\n", exist)

	matches := headerRe.FindAllSubmatch(body, -1)
	for _, m := range matches {
		fmt.Println("fetch card news:", string(m[1]))
	}

	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Println("htmlquery.Parse failed:%v", err)
	}
	nodes := htmlquery.Find(doc, `//div[@class="news_li"]/h2/a[@target="_blank"]`)
	for _, node := range nodes {
		fmt.Println("fetch card ", node.FirstChild.Data)
	}
}

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code:%d", resp.StatusCode)
	}
	reader := bufio.NewReader(resp.Body)
	determinEncoding := DeterminEncoding(reader)
	// 将 HTML 文本从特定编码转换为 UTF-8 编码
	newReader := transform.NewReader(reader, determinEncoding.NewDecoder())

	return io.ReadAll(newReader)
}

func DeterminEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		fmt.Println("fetch error:%v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
