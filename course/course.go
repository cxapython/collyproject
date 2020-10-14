package main

import (
	"collyproject/dbhelper/mongodb"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strings"
)

// Course stores information about a coursera course
type Course struct {
	Title       string
	Description string
	Creator     string
	URL         string
}

func main() {
	db := mongodb.GetS()
	defer db.Close()

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: coursera.org, www.coursera.org
		colly.AllowedDomains("coursera.org", "www.coursera.org"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./coursera_cache"),
	)

	// mongo backend
	//storage := &mongo.Storage{
	//	Database: "course_colly",
	//	URI:      "mongodb://127.0.0.1:27017",
	//}
	//if err := c.SetStorage(storage); err != nil {
	//	panic(err)
	//}

	// Create another collector to scrape course details
	detailCollector := c.Clone()


	c.OnXML("//a[@href]", func(e *colly.XMLElement) {
		link := e.Attr("href")
		if !strings.HasPrefix(link, "/browse") || strings.Index(link, "=signup") > -1 || strings.Index(link, "=login") > -1 {
			return
		}
		if strings.Index(link, "language-learning") > -1 {
			//only scraping language-learning
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
	})
	c.OnResponse(func(r *colly.Response) {
		log.Println("visiting:", r.Ctx.Get("url"))
		//source:=string(r.Body)
		//log.Println("get source",source)
	})
	c.OnXML("//a[contains(@href,'learn')]", func(e *colly.XMLElement) {
		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("courseURL", courseURL)
		if strings.Index(courseURL, "coursera.org/learn") != -1 {
			detailCollector.Visit(courseURL)
		}
	})

	detailCollector.OnXML(`//div[@id="rendered-content"]`, func(e *colly.XMLElement) {
		log.Println("Course found", e.Request.URL)
		title := e.ChildText(".//div[@data-test='banner-title-container']")
		if title == "" {
			log.Println("No title found", e.Request.URL)
		}
		course := Course{
			Title:       title,
			URL:         e.Request.URL.String(),
			Description: e.ChildText(".//div[@class='content']"),
			Creator:     e.ChildText(".//div//h3[contains(@class,'instructor-name')]"),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		// Dump json to the standard output
		enc.Encode(course)
		db.GetC("course").Upsert(bson.M{"url": course.URL}, bson.M{"$set": course})


		log.Println("now data", course)

	})

	// Start scraping on http://coursera.com/browse
	c.Visit("https://coursera.org/browse")



}
