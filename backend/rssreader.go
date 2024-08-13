package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Published   string `xml:"pubDate"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Rss struct {
	Channel Channel `xml:"channel"`
}

/*InsertSummary(SummaryEntry{
	Url:           "test",
	FromWeb:       "NPR",
	Summary:       "TESTTSKTJKLAJFSLA",
	TimeAdded:     time.Now(),
	Title:         "SICKTITLE",
	TimePublished: time.Now(),
})*/

func FullRSSCycle() (unsummarizedEntries []SummaryEntry, text []string) {
	const dur = time.Duration(time.Minute * 60)

	var websites = [...]Website{NPR, AP} //, BBC}

	var rssWG sync.WaitGroup

	// var text []string
	//var links []string

	var mut sync.Mutex
	// var unsummarizedEntries []SummaryEntry
	itemsChecked := 0

	// For every RSS feed
	// Make a go routine to get rss data from the links
	// Once links acquired, make more sub goroutines
	// to then go to each website and scrape
	for _, web := range websites {
		rssWG.Add(1)

		// RSS Feed routine
		go func(web Website) {
			fmt.Println("Routine spun for: " + web.Name)
			// web.RSSLink
			newText, newEntries, newItemsChecked, err := OneScrapeCycle(web)
			if err != nil {
				rssWG.Done()
				return
			}

			mut.Lock()
			text = append(text, newText...)
			unsummarizedEntries = append(unsummarizedEntries, newEntries...)
			itemsChecked += newItemsChecked
			mut.Unlock()
			rssWG.Done()
			fmt.Println("CYCLE DONE FOR " + web.Name)
		}(web)

	}

	rssWG.Wait()
	fmt.Printf("RSS done, parsed items: %d\n", itemsChecked)
	return unsummarizedEntries, text

	//fmt.Println(links[1])
	//fmt.Println(text[1])
}

func attemptTimeParse(checkTime string) (time.Time, error) {
	strs := [...]string{time.ANSIC, time.RFC1123Z, "Mon, 02 Jan 2006 15:04:05 MST"}

	var toRet time.Time
	var err error

	for _, layout := range strs {
		toRet, err = time.Parse(layout, checkTime)

		if err == nil {
			return toRet, nil
		}
	}

	return toRet, err
}

func OneScrapeCycle(web Website) (text []string, entries []SummaryEntry, checked int, err_ret error) {
	rss, err := GetRSSDataFromLink(web.RSSLink)

	if err != nil {
		fmt.Printf("Error getting rss data for %s: %s\n", web, err.Error())
		return nil, nil, 0, err
	}

	var mut sync.Mutex
	var w sync.WaitGroup
	checked = 0

	for _, item := range rss.Channel.Items {
		// for i := 0; i < 3; i++ {
		w.Add(1)
		go func(item Item) {
			defer w.Done()

			if Database.Init && IsInDB(item.Link) {
				fmt.Println(item.Link + " already in DB!")
				return
			}

			// item := rss.Channel.Items[i]
			scrape := web.Scrape(item.Link)

			// Try to parse time until we really can't lmao
			// parsedPublishedTime, err := time.Parse(time.RFC1123Z, item.Published)
			parsedPublishedTime, err := attemptTimeParse(item.Published)

			if err != nil {
				fmt.Println("Could not get parsed time: " + err.Error())
				return
			}
			/*if err != nil {
				err = nil
				parsedPublishedTime, err = time.Parse(time.ANSIC, item.Published)
				if err != nil {

					err = nil
					parsedPublishedTime, err = time.Parse(time.UnixDate, item.Published)
					if err != nil {
						fmt.Println("Could not get parsed time: " + err.Error())
						return
					}
				}
			}*/

			if scrape.allText != "" {
				mut.Lock()
				//links = append(links, item.Link)
				entries = append(entries, SummaryEntry{
					Url:           item.Link,
					FromWeb:       web.Name,
					Summary:       "",
					TimeAdded:     time.Now().UTC(),
					Title:         item.Title,
					TimePublished: parsedPublishedTime,
					PhotoUrl:      scrape.photoUrl,
				})
				checked++
				text = append(text, scrape.allText)
				mut.Unlock()
			}
		}(item)
	}

	w.Wait()
	fmt.Println(web.Name + " routine finished")
	return text, entries, checked, nil
}

func summarizeEntries(entries []SummaryEntry, text []string) (newEntries []SummaryEntry) {

	var mut sync.Mutex
	var insertGroup sync.WaitGroup
	for i, entry := range entries {
		insertGroup.Add(1)
		go func(i int, entry SummaryEntry) {
			defer insertGroup.Done()
			// why sleep? to space out google requests
			// gemini free if you have <60 requests a minute
			time.Sleep(time.Duration(float64(i) * 1.25 * float64(time.Second)))
			var err error
			entry.Summary, err = googleRequest(text[i])
			entry.Summary = strings.ReplaceAll(entry.Summary, "'", "")
			if err != nil {
				fmt.Println("SUM INSERT ERR " + err.Error())
				return
				// continue
			}

			mut.Lock()
			newEntries = append(newEntries, entry)
			mut.Unlock()
			//fmt.Println("Inserting")
			//InsertSummary(entry)
		}(i, entry)
	}
	insertGroup.Wait()

	return newEntries
}

func GetRSSDataFromLink(link string) (Rss, error) {
	// Get data from rss
	res, err := http.Get(link)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// Convert that data into our structs
	body, err := io.ReadAll(res.Body)
	//fmt.Println(string(body))

	if err != nil {
		// fmt.printf("Error parsing rss: %s", err.Error())
		return Rss{}, err
	}

	var rss Rss
	err = xml.Unmarshal(body, &rss)

	if err != nil {
		// panic(err)
		// fmt.printf("Error getting rss: %s", err.Error())
		return Rss{}, err
	}
	return rss, nil
}
