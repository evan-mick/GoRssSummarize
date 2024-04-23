package main

import "github.com/gocolly/colly"

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
		return s
	}
	return s
}

// BELOW HERE ARE GLOBALS
var NPR = Website{RSSLink: "https://feeds.npr.org/1001/rss.xml", Name: "NPR", scrapeFunc: func(htmlstring string) (ScrapeReturn, error) {
	c := colly.NewCollector()

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

// NOT DONE
// https://feeds.bbci.co.uk/news/rss.xml
var BBC = Website{RSSLink: "hi", scrapeFunc: func(htmlstring string) (ScrapeReturn, error) {
	c := colly.NewCollector()

	var ret ScrapeReturn

	c.OnHTML(".storytext", func(e *colly.HTMLElement) {
		ret.allText = e.ChildText("p")
	})

	err := c.Visit(htmlstring)

	if err != nil {
		return ret, err
	}

	return ret, nil
}}
