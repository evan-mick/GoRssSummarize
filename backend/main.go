package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// https://feeds.npr.org/1001/rss.xml

// FLOW OF PROGRAM
// Main parts are the API responses, the RSS reading cycle,
// database communications, then summarizing with Google's Gemini

func main() {
	err := godotenv.Load(".env")
	fmt.Println("STARTING UP")

	// scr := AP.Scrape("https://apnews.com/article/house-speaker-jeffries-johnson-marjorie-taylor-greene-41bf396eca6b0ef3b2bfb71a3cf1fc91")
	// OneScrapeCycle(AP)

	// fmt.Println(txt)

	// return
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dat := SummaryEntry{
		Title: "Test",
	}

	// InitAPIServer()
	// For all items in our rss, save the html to a local file

	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Print("DB connected!")
	// }

	//InsertSummary()
	defer CloseDB()
	// go
	// SelectOneRow()
	go InitAPIServer()

	//FullRSSCycle()
	//OutputMainPage(dat)
	// reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(os.Stdin)
	// CLI
	for {

		// ReadString will block until the delimiter is entered
		// input, err := reader.ReadString('\n')
		// scanner.Scan

		if !scanner.Scan() {
			continue
		}

		input := scanner.Text()
		// if err != nil {
		// 	fmt.Println("An error occured while reading input. Please try again", err)
		// 	continue
		// }
		input = strings.TrimSuffix(input, "\n")

		if input == "refresh" || input == "r" {
			OutputMainPage(dat)
			fmt.Println("REFRESHED")
		} else if input == "cycle" || input == "c" {
			// InitDB()
			// FullRSSCycle()
			CollectAllLocal()
			fmt.Println("COLLECT COMPLETE")
		} else if input == "s" {
			SummarizeLocalCache()
			fmt.Println("SUMMARIZE COMPLETE")
		}

		/*else if input == "DELETEitALlBIGBOi" {
			InitDB()
			DirectSQLCMD("DELETE FROM entries")
			fmt.Println("DELETED")
		}*/

		/*else if input == "s" {
			InitDB()
			CreateLocalCache()
			fmt.Println("CACHE CREATED")
		}*/
	}
}

func SummarizeLocalCache() {

	entries, err := LoadLocalCache()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%d to summarize", len(entries))
	summEntry := summarizeEntries(entries)
	StoreEntriesLocally(summEntry)
}

func CollectAllLocal() {
	allEntries := FullRSSCycle()

	RankEntries(&allEntries)

	for i, entry := range allEntries {
		fmt.Printf("%d: %d %s \n", i, entry.Score, entry.Title)
	}
	//summEntry := summarizeEntries(allEntries)
	StoreEntriesLocally(allEntries[:15])
}
