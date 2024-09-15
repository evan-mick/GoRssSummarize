package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	// "github.com/go-sql-driver/mysql"
	// "github.com/go-sql-driver/mysql"
)

// https://cloud.google.com/sdk/gcloud/reference/sql/connect
// how to connect
// gcloud sql connect summary-entries --user=root
// then table is entries and check for pass

// type Timestamp time.Time

/*
create table if not exists summaries (
  url VARCHAR(3000),
  summary VARCHAR(65535),
  timeAdded timestamp,
  timePublished timestamp,
  PRIMARY KEY (`url`)
);

FROM GPT
CREATE TABLE IF NOT EXISTS entries (
    url VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255),
    from VARCHAR(255),
    summary TEXT,
    time_added TIMESTAMP,
    time_published TIMESTAMP
);
*/

/*
INSERT IGNORE INTO summaries (url, summary, time)
VALUES (?, ?, ?)
*/

/*

 */

// struct
type SummaryEntry struct {
	Url           string    `json:"url"` // Also acts as an ID
	Title         string    `json:"title"`
	FromWeb       string    `json:"fromWeb"`
	Summary       string    `json:"summary"`
	TimeAdded     time.Time `json:"timeAdded"`
	TimePublished time.Time `json:"timePublished"`
	PhotoUrl      string    `json:"photoUrl"`

	FullText string `json:"text"` // should I do this??
	Score    int    `json:"score"`
}

const TimeLayout = "2006-01-02 15:04:05"

type DatabaseInfo struct {
	DB   *sql.DB
	Init bool
}

var Database = DatabaseInfo{
	DB:   nil,
	Init: false,
}

func InitDB() {

	if Database.Init {
		fmt.Println("Database already initialized")
		return
	}

	var (
		databaseName       = os.Getenv("DB_DBNAME")
		user               = os.Getenv("DB_USER")
		password           = os.Getenv("DB_PASS")
		instanceConnection = os.Getenv("DB_INSTANCENAME")
	)
	cfg := mysql.Cfg(instanceConnection, user, password)
	cfg.DBName = databaseName
	db, err := mysql.DialCfg(cfg)

	if err != nil {
		log.Fatal("Could not connect to database")
		return
	} else {
		fmt.Println("DB Opened")
	}
	Database.DB = db

	err = db.Ping()

	if err != nil {
		log.Fatal("DB ping error: " + err.Error())
	} else {
		fmt.Println("Pinged successfully, connected")
	}

	// _, err = db.Exec(`DROP TABLE entries`)

	// Important that the if not exists there
	// cause otherwise it will throw error and below logic won't work
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS entries (
		url VARCHAR(255) PRIMARY KEY,
		title VARCHAR(255),
		fromWeb VARCHAR(255),
		summary TEXT,
		timeAdded TIMESTAMP,
		timePublished TIMESTAMP,
		photoUrl VARCHAR(255)
	);	
	`)

	/*InsertSummary(SummaryEntry{
		Url:           "test",
		FromWeb:       "NPR",
		Summary:       "TESTTSKTJKLAJFSLA",
		TimeAdded:     time.Now(),
		Title:         "SICKTITLE",
		TimePublished: time.Now(),
	})*/

	if err != nil {
		db.Close()
		log.Fatal("Issue with entries table " + err.Error())
		return
	}

	Database.Init = true
	fmt.Println("Post exec successfully")

}

// / FOR THE LOVE OF GOD DO NOT GIVE FRONTEND ACCESS TO THIS
// specifically for backend CLI
func DirectSQLCMD(cmd string) {
	res, err := Database.DB.Exec(cmd)

	if err != nil {
		fmt.Print("Direct SQL error: " + err.Error() + "\n")
		return
	}

	rows, _ := res.RowsAffected()
	ins, _ := res.LastInsertId()

	fmt.Printf("SQL SUCCESSS: rows affected: %d, last insert ID: %d", rows, ins)
}

func InsertSummary(entry SummaryEntry) {

	/*cmd := fmt.Sprintf(`INSERT INTO entries (url, title, fromweb, summary, timeAdded, timePublished)
	VALUES ('%s', '%s', '%s', '%s', '%s', '%s')`,
		entry.Url, entry.Title, entry.FromWeb, entry.Summary, entry.TimeAdded.Format(TimeLayout), entry.TimePublished.Format(TimeLayout))*/

	cmd := `INSERT INTO entries (url, title, fromweb, summary, timeAdded, timePublished, photoUrl) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := Database.DB.Exec(cmd, entry.Url, entry.Title, entry.FromWeb, entry.Summary, entry.TimeAdded.Format(TimeLayout), entry.TimePublished.Format(TimeLayout), entry.PhotoUrl)

	if err != nil {
		fmt.Println("Error with insertion into database " + err.Error())
	}
}
func IsInDB(check string) bool {

	cmd := "SELECT 1 FROM entries WHERE url = ?;"

	var result int
	err := Database.DB.QueryRow(cmd, check).Scan(&result)

	if err != nil {
		if err != sql.ErrNoRows {
			// actual error
			fmt.Println("Check in DB error: " + err.Error())
		}
		return false
	}
	return true
}

func SelectOneRow() SummaryEntry {

	cmd := fmt.Sprintf(`SELECT * FROM entries ORDER BY timePublished LIMIT 1`)
	row := Database.DB.QueryRow(cmd)

	if row.Err() != nil {
		fmt.Println("Error: " + row.Err().Error())
		return SummaryEntry{}
	}

	var nxtEntry SummaryEntry
	var (
		published string
		added     string
	)
	row.Scan(&nxtEntry.Url, &nxtEntry.Title,
		&nxtEntry.FromWeb, &nxtEntry.Summary,
		&added, &published, &nxtEntry.PhotoUrl)
	//&nxtEntry.TimeAdded, &nxtEntry.TimePublished)
	var err error
	nxtEntry.TimePublished, err = time.Parse(TimeLayout, published)
	nxtEntry.TimeAdded, err = time.Parse(TimeLayout, added)

	if err != nil {
		fmt.Println("select parse error")
	}

	//fmt.Print(nxtEntry)
	return nxtEntry
	// json.Unmarshal(data, v)
}

func parseRowsToSummary(rows *sql.Rows) ([]SummaryEntry, error) {
	var returnEntries []SummaryEntry

	var err error = nil
	// var nxtEntry SummaryEntry
	var (
		published string
		added     string
	)
	for rows.Next() {
		fmt.Println("Iterating through entry")
		var nxtEntry SummaryEntry

		err = rows.Scan(&nxtEntry.Url, &nxtEntry.Title,
			&nxtEntry.FromWeb, &nxtEntry.Summary,
			&added, &published, &nxtEntry.PhotoUrl)

		if err != nil {
			fmt.Println("Trouble parsing one of the scans " + err.Error())
			continue
		}

		nxtEntry.TimePublished, err = time.Parse(TimeLayout, published)

		if err != nil {
			fmt.Println("Trouble parsing time published: " + err.Error())
			continue
		}

		nxtEntry.TimeAdded, err = time.Parse(TimeLayout, added)

		if err != nil {
			fmt.Println("Trouble parsing time added: " + err.Error())
			continue
		}

		returnEntries = append(returnEntries, nxtEntry)
	}

	return returnEntries, nil
	// json.Unmarshal(data, v)
}

func SelectNRows(number int, startAt int) ([]SummaryEntry, error) {

	cmd := `SELECT * FROM entries ORDER BY timePublished DESC LIMIT ? OFFSET ?`
	rows, err := Database.DB.Query(cmd, number, number+startAt-1)
	// 2024-04-03T16:02:30Z

	if err != nil {
		fmt.Println("Query error")
		return nil, err
	}
	parsed, err := parseRowsToSummary(rows)
	return parsed, err
}

func SelectAllRows() ([]SummaryEntry, error) {
	cmd := `SELECT * FROM entries ORDER BY timePublished DESC`
	rows, err := Database.DB.Query(cmd)
	if err != nil {
		fmt.Println("Query error")
		return nil, err
	}
	parsed, err := parseRowsToSummary(rows)
	return parsed, err
}

func CloseDB() {
	if !Database.Init {
		return
	}
	Database.Init = false
	Database.DB.Close()
	Database.DB = nil

}
