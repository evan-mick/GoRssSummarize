package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
)

const templateOutput = `./frontend/`

//../frontend/static/`

type htmlDat struct {
	Title           string
	Entries         []SummaryEntry
	MainListDefined bool
}

func OutputMainPage() {

	t1 := template.New("Main Page")

	read, err := os.ReadFile("frontend_data/mainTemplate.html")
	if err != nil {
		log.Printf("Error with reading file: " + err.Error())
		return
	}
	temp, err := t1.Parse(string(read))

	if err != nil {
		log.Printf("Error with template parsing: " + err.Error())
		return
	}

	file, err := os.Create(templateOutput + "index.html")

	if err != nil {
		log.Printf("Error with file parsing: " + err.Error())
		return
	}

	var entries []SummaryEntry = []SummaryEntry{}
	var entriesLoaded bool = true
	entries, err = LoadLocalCache() //SelectAllRows() // SelectNRows(5, 0)
	// If we can't find cache, we're loading
	if err != nil {
		entriesLoaded = false
		// log.Printf("Error with getting files: " + err.Error())
	}
	// log.Print(entries)

	toDisplay := htmlDat{
		Title:           "The Headliner",
		Entries:         entries,
		MainListDefined: entriesLoaded,
	}

	temp.Execute(file, toDisplay)

}

func LoadLocalCache() ([]SummaryEntry, error) {

	file, err := os.ReadFile("frontend_data/data.json")
	if err != nil {
		fmt.Println("Local file read error " + err.Error())
		return nil, err
	}
	var retVal []SummaryEntry

	err = json.Unmarshal(file, &retVal)
	if err != nil {
		fmt.Println("Local file unmarshal error " + err.Error())
		return nil, err
	}

	return retVal, nil
}

func StoreEntriesLocally(writeEntries []SummaryEntry) {

	file, err := os.Create("frontend_data/data.json")
	if err != nil {
		fmt.Println("Local file create error " + err.Error())
		return
	}
	dat, err := json.Marshal(writeEntries)
	if err != nil {
		log.Printf("Error with marshal: " + err.Error())
		return
	}
	file.Write(dat)
	if err != nil {
		log.Printf("Error with write file: " + err.Error())
		return
	}
}

/*func CreateLocalCache() {
	file, err := os.Create("frontend_data/data.json")
	if err != nil {
		fmt.Println("Local file create error " + err.Error())
		return
	}

	writeEntries, err := SelectNRows(18, 0)
	if err != nil {
		log.Printf("Error with getting files: " + err.Error())
		return
	}
	dat, err := json.Marshal(writeEntries)
	if err != nil {
		log.Printf("Error with marshal: " + err.Error())
		return
	}
	file.Write(dat)
	if err != nil {
		log.Printf("Error with write file: " + err.Error())
		return
	}
}*/
