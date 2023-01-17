package main

import (
	"github.com/nine-monsters/crawler/collect"
	"github.com/nine-monsters/crawler/engine"
	"github.com/nine-monsters/crawler/log"
	"github.com/nine-monsters/crawler/proxy"
	"go.uber.org/zap/zapcore"
	"regexp"
	"time"
)

var (
	headerRe = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)
)

func main() {
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	// proxy
	var proxyURLs []string
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed")
	}

	var f collect.Fetcher = &collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Logger:  logger,
		Proxy:   p,
	}
	seeds := make([]*collect.Task, 0, 1000)
	seeds = append(seeds, &collect.Task{
		Name:    "find_douban_sun_room",
		Fetcher: f,
	})

	s := engine.NewEngine(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithSeeds(seeds),
	)
	s.Run()
}
