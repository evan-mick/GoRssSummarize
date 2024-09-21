package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	// Consider gocron?
	"time"
)

// https://feeds.npr.org/1001/rss.xml

// FLOW OF PROGRAM
// Main parts are the API responses, the RSS reading cycle,
// database communications, then summarizing with Google's Gemini

var quit bool = false
var loopTime time.Duration = time.Hour * 6

// var currentLoopTimer time.Duration = loopTime
var lastLoopCheck time.Time = time.Now()
var toLoopTime time.Time = time.Now()

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

	defer CloseDB()
	go InitAPIServer()

	scanner := bufio.NewScanner(os.Stdin)
	// CLI
	for {
		if !scanner.Scan() {
			continue
		}

		input := scanner.Text()
		input = strings.TrimSuffix(input, "\n")

		if input == "refresh" || input == "r" {
			OutputMainPage()
			fmt.Println("REFRESHED")
		} else if input == "cycle" || input == "c" {
			// InitDB()
			// FullRSSCycle()
			CollectAllLocal()
			fmt.Println("COLLECT COMPLETE")
		} else if input == "s" {
			SummarizeLocalCache()
			fmt.Println("SUMMARIZE COMPLETE")
		} else if input == "b" {
			fmt.Println("BEGINNING FULL TIMED LOOP")
			go MainLoop()
		} else if input == "q" {
			quit = true
			break
		} else if input == "t" {
			fmt.Printf("%f minutes left\n", (toLoopTime.Sub(time.Now())).Minutes())
		} else if input == "st" {
			fmt.Println("Enter time (in minutes) to set current loop to")
			scanner.Scan()
			in := scanner.Text()
			if i, err := strconv.Atoi(in); err == nil {
				toLoopTime = time.Now().Add(time.Minute * time.Duration(i))
				continue
			}
			fmt.Println("Invalid time input")

		} else if input == "stl" {
			fmt.Println("Enter time (in minutes) to set loop time to")
			scanner.Scan()
			in := scanner.Text()
			if i, err := strconv.Atoi(in); err == nil {
				loopTime = time.Minute * time.Duration(i)
				continue
			}
			fmt.Println("Invalid time input")
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

func MainLoop() {

	for {
		if quit {
			break
		}

		time.Sleep(time.Second)
		//currentLoopTimer += time.Second

		if time.Now().Compare(toLoopTime) > -1 {
			//		}
			//		if currentLoopTimer >= loopTime {
			fmt.Println("TIMER COMPLETE, BEGINNING FULL REFRESH")
			go RunOneFullRefresh()
			toLoopTime = time.Now().Add(loopTime)
		}
	}
}

func RunOneFullRefresh() {

	CollectAllLocal()
	fmt.Println("FULL LOOP: COLLECT COMPLETE")
	fmt.Println("FULL LOOP: SUMMARIZE START")
	SummarizeLocalCache()
	fmt.Println("FULL LOOP: SUMMARIZE END")
	OutputMainPage()
	fmt.Println("FULL LOOP: OUTPUT COMPLETE")
	fmt.Println("FULL LOOP: FULL REFRESH COMPLETE")

}

func SummarizeLocalCache() {

	entries, err := LoadLocalCache()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%d to summarize\n", len(entries))
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
