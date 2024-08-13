package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
)

const templateOutput = `./frontend/static/`

//../frontend/static/`

type htmlDat struct {
	Title   string
	Entries []SummaryEntry
}

func OutputMainPage(dat SummaryEntry) {

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

	entries, err := LoadLocalCache() //SelectAllRows() // SelectNRows(5, 0)
	if err != nil {
		// log.Printf("Error with getting files: " + err.Error())
		return

	}
	// log.Print(entries)

	toDisplay := htmlDat{
		Title:   "Website",
		Entries: entries,
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
