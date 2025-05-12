
This is the source code for [The Headliner](https://theheadliner.news)!

Notably, a .env file is needed with a Google Gemini API key. It needs the input

```
GOOGLE_AI_KEY = YOUR_KEY_HERE
```

Once that is in place, to run it, you can do a simple

```
go build
```

In the terminal, and run

```
./rsssummarize
```


# Docker

This project can be built as a docker image. It also can be uploaded to fly.io with the fly.toml file.
Make sure the docker daemon is running, then run

```
docker build --platform=linux/amd64 --tag YOUR_NAME/headliner:latest .
```

To make a docker build with YOUR_NAME and headliner:latest as the tag name. You can run the docker image and its output will be exposed on port 8080. 


# How it works

Upon start up, The Headliner goes through a series of RSS feeds, mainly from the BBC and NPR. It compiles all the articles, then web scraps their contents. Mainly their article content and their titles. Then, all the content is parsed through, and scored based on what words are in their content. The (up to) top 15 postively scored articles are sent to Google Gemini for summarizing, then once received turned into a one-page website based on a template with the post information added. From there, when a user connects, they will see the rendered website. 


# File Overview

### main.go
The entry point of the program, handles the REPL interface and the initialization of articles

### rssreader.go
Given an RSS link, gets a list of articles, scraps them, summarizes them, then stores them. A lot of key high level logic occurs here.

### api_server.go
The logic for serving web pages. Also code for getting some number of articles though that is disabled on the live website.

### api_requests.go
Logic for api requests. Right now, it is only the Google gemini API function. 

### api_htmlTemplates.go
Creates and stores the html pages from the html templates. Requires that data is stored locally to be put onto the site. 

### parser.go
Parses and ranks articles based on their article content and titles. 

### scraper.go
All the logic for scraping websites. In here is the logic for scraping each individual news site given a link to one of their articles. 

### database.go
Mostly unused. But some stuff for connect, adding, and removing this from a database goes here. Was originally planning on storing article data but realized it was not super necessary if the templated things get stored. 

# Website ranking

**points.json** is a file for storing keywords. Each of the rankings adds a value to a score for each of the articles. 
The dictionary can be found in parser for each of the conversions. Here is a snippet

```

var rankingToPoints = map[string]int{
	// H -> Hight, D -> (natural) Disaster, M -> Medium priority, L -> Low priority
	"H": 10,
	"D": 8,
	"M": 5,
	"L": 3,

	// B -> bad, no more acronyms these are arbitrary, W -> worse, N -> "No" worst, basically disable
	"B": -3,
	"W": -10,
	"N": -50,
}
```

For every word in each category that appears in the article, it gets the number of given points added. If its in the title, it gets a 10x multiplier for the given score. 
Its a bit of a naive scoring algorithm, but works well enough for my personal usage. 

