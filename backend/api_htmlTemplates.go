package main

import (
	"html/template"
	"log"
	"os"
)

const templateOutput = `../frontend/static/`

type htmlDat struct {
	Title   string
	Entries []SummaryEntry
}

func OutputMainPage(dat SummaryEntry) {

	t1 := template.New("Main Page")

	read, err := os.ReadFile("mainTemplate.html")
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

	entries, err := SelectAllRows() // SelectNRows(5, 0)
	if err != nil {
		log.Printf("Error with getting files: " + err.Error())
		return

	}
	log.Print(entries)

	toDisplay := htmlDat{
		Title:   "Website",
		Entries: entries,
	}

	temp.Execute(file, toDisplay)

}
