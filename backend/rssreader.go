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

func FullRSSCycle() {
	const dur = time.Duration(time.Minute * 60)

	var websites = [...]Website{NPR} //, BBC}

	var rssWG sync.WaitGroup

	var text []string
	//var links []string

	var mut sync.Mutex
	var entries []SummaryEntry
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
			OneScrapeCycle(web, &mut, &text, &entries, &itemsChecked)
			rssWG.Done()
		}(web)

	}

	rssWG.Wait()
	fmt.Printf("RSS done, parsed items: %d\n", itemsChecked)

	//fmt.Println(links[1])
	//fmt.Println(text[1])

	summarizeAndInsertEntries(entries, text)
}

func OneScrapeCycle(web Website, mut *sync.Mutex, text *[]string, entries *[]SummaryEntry, itemsChecked *int) {
	rss := GetRSSDataFromLink(web.RSSLink)

	var w sync.WaitGroup
	for _, item := range rss.Channel.Items {
		// for i := 0; i < 3; i++ {
		w.Add(1)
		go func(item Item) {
			defer w.Done()

			if IsInDB(item.Link) {
				fmt.Println(item.Link + " already in DB!")
				return
			}

			// item := rss.Channel.Items[i]
			scrape := web.Scrape(item.Link)

			parsedPublishedTime, err := time.Parse(time.RFC1123Z, item.Published)

			if err != nil {
				fmt.Println("Could not get parsed time: " + err.Error())
				return
			}

			if scrape.allText != "" {
				mut.Lock()
				//links = append(links, item.Link)
				*entries = append(*entries, SummaryEntry{
					Url:           item.Link,
					FromWeb:       NPR.Name,
					Summary:       "",
					TimeAdded:     time.Now().UTC(),
					Title:         item.Title,
					TimePublished: parsedPublishedTime,
					PhotoUrl:      scrape.photoUrl,
				})
				*itemsChecked++
				*text = append(*text, scrape.allText)
				mut.Unlock()
			}
		}(item)
	}

	w.Wait()
	fmt.Println(web.Name + " routine finished")

}

func summarizeAndInsertEntries(entries []SummaryEntry, text []string) {

	var insertGroup sync.WaitGroup
	for i, entry := range entries {
		insertGroup.Add(1)
		go func(i int, entry SummaryEntry) {
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

			fmt.Println("Inserting")
			InsertSummary(entry)
			insertGroup.Done()
		}(i, entry)
	}
	insertGroup.Wait()
}

func GetRSSDataFromLink(link string) Rss {
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
		panic(err)
	}

	var rss Rss
	err = xml.Unmarshal(body, &rss)

	if err != nil {
		panic(err)
	}
	return rss
}
