package fetcher

import (
	"collyproject/config"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

var Conf = config.Config

func Random(strings []string, length int) (string, error) {
	if len(strings) <= 0 {
		return "", errors.New("the length of the parameter strings should not be less than 0")
	}

	if length <= 0 || len(strings) <= length {
		return "", errors.New("the size of the parameter length illegal")
	}

	for i := len(strings) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}

	str := ""
	for i := 0; i < length; i++ {
		str += strings[i]
	}
	return str, nil
}
func getProxies() string{
	proxyHost := "http-dyn.abuyun.com"
	proxyPort := "9020"

	var proxyList = [] string{
		//your user and password
		"H01234567890123P:0123456789012345",
		"H01234567890123S:SSSSDDDDDDDDDDDDD",

	}
	proxyInfo,err := Random(proxyList,1)
	if err!=nil{
		panic(err)
	}
	arr:=strings.Split(proxyInfo,":")
	proxyUser:=arr[0]
	proxyPass:=arr[1]
	proxyMeta := fmt.Sprintf("%s:%s@%s:%s" ,
		proxyUser,
		proxyPass,
		proxyHost,
		proxyPort,
	)
	return proxyMeta
}
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
	if Conf.GetBool("CRAWLER.USE_PROXY"){
		proxyStr:=getProxies()
		log.Println("Get Proxy:",proxyStr)
		rp, err := proxy.RoundRobinProxySwitcher(fmt.Sprintf("http://%s",proxyStr))
		if err != nil {
			log.Fatal(err)
		}
		c.SetProxyFunc(rp)
	}
	c.OnError(func(r *colly.Response, e error) {
		if r.StatusCode != 200 {
			log.Println("访问失败")
		}
	})
}
