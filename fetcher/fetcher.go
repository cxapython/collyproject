package fetcher

import (
	"collyproject/config"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"log"
	"net"
	"net/http"
	"time"
)

var Conf = config.Config

func CreateCollector() *colly.Collector {
	debugger := &debug.LogDebugger{}
	c := colly.NewCollector(
		colly.Async(Conf.GetBool("CRAWLER.Async")),
		colly.Debugger(debugger),
		colly.IgnoreRobotsTxt(),
		colly.AllowURLRevisit(),//是否允许相同的url被二次访问。
		colly.AllowedDomains("coursera.org", "www.coursera.org"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./coursera_cache"),
	)
	//disable http keep alives
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   120 * time.Second, // 超时时间
			KeepAlive: 30 * time.Second,  // keepAlive 超时时间
		}).DialContext,
	})
	setDefaultsSetting(c)
	extensions.RandomUserAgent(c)
	return c
}
func setDefaultsSetting(c *colly.Collector) {
	delay := time.Duration(Conf.GetInt("CRAWLER.DELAY"))
	randomDelay := time.Duration(Conf.GetInt("CRAWLER.RANDOM_DELAY"))
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: Conf.GetInt("CRAWLER.Parallelism"),
		Delay: delay * time.Second, RandomDelay: randomDelay * time.Second})
	c.OnError(func(r *colly.Response, e error) {
		if r.StatusCode != 200 {
			log.Println("访问失败")
		}
	})
}
