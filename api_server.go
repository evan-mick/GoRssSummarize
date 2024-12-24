/// The actual API server

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Ping is a request handler that simple serves a test html file.
func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Println("PINGED")
	//http.ServeFile(w, req, "../frontend/dev/ping.html")
}

var fs http.Handler

// ServeMain is a requesthandler that serves the main static website
func serveMain(w http.ResponseWriter, req *http.Request) {
	ip := req.Header.Get("X-FORWARDED-FOR")
	fmt.Println("MAIN SERVED IP: " + ip)
	//http.ServeFile(w, req, "./frontend/static")
	fs.ServeHTTP(w, req)

}

// ReturnOneEntry is an API request handler that returns a single row from the database
func returnOneEntry(w http.ResponseWriter, req *http.Request) {
	fmt.Println("ONE ENTRY")
	data := SelectOneRow()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(data)
	//http.ServeFile(w, req, "./frontend/static")
}

func returnManyEntry(w http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()
	startS := params.Get("start")
	numberS := params.Get("number")

	number := 10
	start := 0

	if startS != "" {
		val, err := strconv.Atoi(startS)
		if err == nil {
			start = val
		}
	}

	if numberS != "" {
		val, err := strconv.Atoi(numberS)
		if err == nil {
			number = val
		}
	}

	fmt.Printf("MANY ENTRY, ARGS: %d %d", number, start)
	data, err := SelectNRows(number, start)

	//fmt.Println(data)

	if err != nil {
		// TODO: better handling, should return something else
		fmt.Println("Error with selection: " + err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
	//http.ServeFile(w, req, "./frontend/static")
}

// InitAPIServer initializes and runs the HTTP server
// Because it runs ListenAndServe which runs forever, it should
// either be called last or in a goroutine
func InitAPIServer() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/requests/entry/", returnOneEntry)

	http.HandleFunc("/requests/entries/", returnManyEntry)

	// Make sure this one is last
	// http.HandleFunc("/", serveMain)
	fs = http.FileServer(http.Dir(templateOutput))
	http.HandleFunc("/", serveMain)
	//http.Handle("/", serveMain)

	// for many entries, needs to have args for what section to give back out
	// make sure the bounds of such are checked
	// http.HandleFunc("/requests/entries")

	fmt.Println("API server initialized")

	http.ListenAndServe(":"+port, nil)

}
