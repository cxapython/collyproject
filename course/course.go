package main

import (
	"collyproject/dbhelper/mongodb"
	"collyproject/fetcher"
	"context"
	"encoding/json"
	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strings"
)

// Course stores information about a coursera course
type Course struct {
	Title       string `bson:"title"`
	Description string `bson:"description"`
	Creator     string `bson:"author"`
	URL         string `bson:"url"`
}


func main() {
	db := mongodb.GetConnection()

	// Instantiate default collector
	c := fetcher.CreateCollector()
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

		db.Collection("course").UpdateOne(context.TODO(),bson.M{"url": course.URL}, bson.M{"$set": course},options.Update().SetUpsert(true))


		log.Println("now data", course)

	})

	// Start scraping on http://coursera.com/browse
	c.Visit("https://coursera.org/browse")



}
