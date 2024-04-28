package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

// Type
// BBC, NPR are all websites
// They have their own RSS feeds that give you links to websites that can then be scraped
// so struct with rss feed, and function for scrapping specific website

type Website struct {
	RSSLink    string
	Name       string
	scrapeFunc func(htmlstring string) (ScrapeReturn, error)
}

type ScrapeReturn struct {
	allText  string
	photoUrl string
}

func (b *Website) Scrape(htmlstring string) ScrapeReturn {

	s, err := b.scrapeFunc(htmlstring)

	if err != nil {
		fmt.Println("SCRAPE ERROR: " + err.Error())
		return s
	}
	return s
}

func getDefaultCollector() *colly.Collector {
	c := colly.NewCollector()
	c.CheckHead = true
	c.DisableCookies()
	c.IgnoreRobotsTxt = true
	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36")
		r.Headers.Set("User-Agent", "THE HEADLINER (NON-COMMERCIAL) WEB-SCRAPER (THIS BOT COLLECTS ARTICLE DATA TO SUMMARIZE)")
	})
	return c
}

// BELOW HERE ARE GLOBALS
var NPR = Website{RSSLink: "https://feeds.npr.org/1001/rss.xml", Name: "NPR", scrapeFunc: func(htmlstring string) (ScrapeReturn, error) {
	c := getDefaultCollector()

	var ret ScrapeReturn

	c.OnHTML(".storytext", func(e *colly.HTMLElement) {
		ret.allText = e.ChildText("p")
		//ret.photoUrl = "https://media.npr.org/chrome_svg/npr-logo.svg"
	})

	var def bool
	// THIS NOT WORKING
	c.OnHTML("div#storytext div.bucketwrap picture img.img", func(e *colly.HTMLElement) {
		first := e.DOM.First()
		var exists bool
		// Extract the `src` attribute value

		if !def {
			ret.photoUrl, exists = first.Attr("src")

			if !exists {
				ret.photoUrl = "https://media.npr.org/chrome_svg/npr-logo.svg"
			}
			def = true
		}

	})

	err := c.Visit(htmlstring)

	if err != nil {
		return ret, err
	}

	return ret, nil
}}

// MUST ADD PHOTO
var AP = Website{RSSLink: "https://apnews.com/index.rss", scrapeFunc: func(htmlstring string) (ScrapeReturn, error) {
	c := getDefaultCollector()

	var ret ScrapeReturn

	c.OnHTML(".RichTextStoryBody", func(e *colly.HTMLElement) {
		fmt.Println("ON HTML")
		ret.allText = e.ChildText("p")
	})

	err := c.Visit(htmlstring)

	if err != nil {
		return ret, err
	}

	return ret, nil
}}

// MUST ADD PHOTO
// RSS DOESN'T WORK, NEED ALT SCRIPT
var Reuters = Website{RSSLink: "https://www.reutersagency.com/feed/?best-sectors=economy&post_type=best", scrapeFunc: func(htmlstring string) (ScrapeReturn, error) {
	c := getDefaultCollector()

	var ret ScrapeReturn

	c.OnHTML(".article-body__content__17Yit", func(e *colly.HTMLElement) {
		fmt.Println("ON HTML")
		ret.allText = e.ChildText("div")
	})

	err := c.Visit(htmlstring)

	if err != nil {
		return ret, err
	}

	return ret, nil
}}

// NOT WORKING
// https://feeds.bbci.co.uk/news/rss.xml
var BBC = Website{RSSLink: "https://feeds.bbci.co.uk/news/world/rss.xml", scrapeFunc: func(htmlstring string) (ScrapeReturn, error) {
	c := getDefaultCollector()
	c.Async = true

	var ret ScrapeReturn

	c.OnHTML("p", func(e *colly.HTMLElement) {
		fmt.Println("IN HTML ")
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Page fully loaded and scraped.")
	})

	err := c.Visit(htmlstring)

	if err != nil {
		return ret, err
	}
	time.Sleep(2 * time.Second)

	c.Wait()

	return ret, nil
}}
