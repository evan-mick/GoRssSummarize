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

	//scr := AP.Scrape("https://apnews.com/article/israel-hamas-war-news-04-27-2024-7ea816cac94138492f7dddf2865c1d2f")

	//fmt.Println(scr.allText)

	//return
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

	// CLI
	for {

		reader := bufio.NewReader(os.Stdin)
		// ReadString will block until the delimiter is entered
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		input = strings.TrimSuffix(input, "\n")

		if input == "refresh" || input == "r" {
			OutputMainPage(dat)
		} else if input == "cycle" || input == "c" {
			InitDB()
			FullRSSCycle()
			CreateLocalCache()
		} else if input == "DELETEitALlBIGBOi" {
			InitDB()
			DirectSQLCMD("DELETE FROM entries")
		}
	}

}
