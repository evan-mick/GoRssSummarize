package main

import (
	"fmt"

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
	//c.CheckHead = true
	c.DisableCookies()
	c.IgnoreRobotsTxt = true
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36")
		//r.Headers.Set("User-Agent", "THE HEADLINER (NON-COMMERCIAL) WEB-SCRAPER (THIS BOT COLLECTS ARTICLE DATA TO SUMMARIZE)")
	})
	return c
}

// BELOW HERE ARE GLOBALS
var NPR = Website{RSSLink: "https://feeds.npr.org/1002/rss.xml", Name: "NPR", scrapeFunc: func(htmlstring string) (ScrapeReturn, error) {
	c := getDefaultCollector()

	var ret ScrapeReturn

	c.OnHTML(".storytext", func(e *colly.HTMLElement) {
		ret.allText = e.ChildText("p")
		//ret.photoUrl = "https://media.npr.org/chrome_svg/npr-logo.svg"
	})

	var def bool = false
	c.OnHTML(".storytext picture img.img", func(e *colly.HTMLElement) {
		first := e.DOM.First()
		var exists bool

		if !def {
			ret.photoUrl, exists = first.Attr("src")

			if !exists {
				ret.photoUrl = "https://media.npr.org/chrome_svg/npr-logo.svg"
			}
			def = true
		}

	})

	err := c.Visit(htmlstring)

	if !def {
		ret.photoUrl = "https://media.npr.org/chrome_svg/npr-logo.svg"
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}}

// MUST ADD PHOTO
var AP = Website{RSSLink: "https://apnews.com/index.rss", Name: "AP", scrapeFunc: func(htmlstring string) (ScrapeReturn, error) {
	c := getDefaultCollector()

	var ret ScrapeReturn

	c.OnHTML(".RichTextStoryBody", func(e *colly.HTMLElement) {
		// fmt.Println("ON HTML")
		ret.allText = e.ChildText("p")
	})

	c.OnHTML(".Page-main .CarouselSlide-media img.Image", func(e *colly.HTMLElement) {
		// first := e.DOM.First()
		// Extract the `src` attribute value
		if ret.photoUrl == "" {
			ret.photoUrl = e.Attr("src")

			// ret.photoUrl = strings.Split(ret.photoUrl, "\n")[0]
			// fmt.Println(ret.photoUrl)
		}

	})
	err := c.Visit(htmlstring)

	if err != nil {
		return ret, err
	}

	if ret.photoUrl == "" {
		ret.photoUrl = "https://dims.apnews.com/dims4/default/6e4b276/2147483647/strip/true/crop/640x236+0+0/resize/640x236!/quality/90/?url=https%3A%2F%2Fassets.apnews.com%2Fc3%2F4c%2F65482a7b452db66043542c093eaf%2Fpromo-2x.png" //"https://assets.apnews.com/fa/ba/9258a7114f5ba5c7202aaa1bdd66/aplogo.svg"
	}

	return ret, nil
}}

// MUST ADD PHOTO
// RSS DOESN'T WORK, NEED ALT SCRIPT

// By Reuters - <a rel="nofollow" class="external text" href="http://www.reuters.com">www.reuters.com</a>, Public Domain, <a href="https://commons.wikimedia.org/w/index.php?curid=149082496">Link</a>
var Reuters = Website{RSSLink: "https://www.reutersagency.com/feed/?best-sectors=economy&post_type=best", scrapeFunc: func(htmlstring string) (ScrapeReturn, error) {
	c := getDefaultCollector()

	var ret ScrapeReturn
	ret.photoUrl = "https://commons.wikimedia.org/w/index.php?curid=149082496"

	c.OnHTML(".article-body__content__17Yit", func(e *colly.HTMLElement) {
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
// https://feeds.bbci.co.uk/news/world/rss.xml"
var BBC = Website{RSSLink: "https://feeds.bbci.co.uk/news/rss.xml", Name: "BBC", scrapeFunc: func(htmlstring string) (ScrapeReturn, error) {
	c := getDefaultCollector()
	// c.Async = true

	var ret ScrapeReturn

	c.OnHTML("article", func(e *colly.HTMLElement) {

		ret.allText = e.ChildText("p")
		//fmt.Println(ret.allText)
	})

	var exists bool = false
	c.OnHTML(".sc-a34861b-1 img", func(e *colly.HTMLElement) {
		ele := e.DOM.First()

		ret.photoUrl, exists = ele.Attr("src")

		if !exists {
			ret.photoUrl = "https://1000logos.net/wp-content/uploads/2016/10/BBC-Logo.jpg"
		}
	})

	if !exists {
		ret.photoUrl = "https://1000logos.net/wp-content/uploads/2016/10/BBC-Logo.jpg"
	}

	/*c.OnScraped(func(r *colly.Response) {
		fmt.Println("Page fully loaded and scraped.")
	})*/

	err := c.Visit(htmlstring)

	if err != nil {
		return ret, err
	}
	// time.Sleep(2 * time.Second)

	// c.Wait()

	return ret, nil
}}
