package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

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

var titleMulti = 10

// Map of score name to list of words
// setup like "h" : [ "word1", "word2" ]
var _scoringMap map[string][]string = nil
var _jsonInfoPath string = "./points.json"

func RankEntries(entries *[]SummaryEntry) {

	if entries == nil {
		return
	}

	for i, entry := range *entries {

		sentences := CheckNumberOfSentences(entry.FullText)

		// Probably not even a full article
		if sentences < 6 {
			(*entries)[i].Score = -1000
			continue
		}

		score := GetWebsiteScore(entry.FullText, entry.Title)

		if score > 0 {
			score += sentences * 3
		}

		(*entries)[i].Score = score
	}

	sort.Slice((*entries)[:], func(i, j int) bool {
		return (*entries)[i].Score > (*entries)[j].Score
	})
}

// Default get website score
// Stateful, uses global maping and json path variables
func GetWebsiteScore(body string, title string) int {
	if _scoringMap == nil {
		data, err := os.ReadFile(_jsonInfoPath)
		if err != nil {
			fmt.Println("ERROR READING POINT JSON FILE " + err.Error())
			return 0
		}
		_scoringMap = UnmarshalScoringInformation(data)
	}

	if _scoringMap == nil {
		return -1
	}

	return GetWebsiteScoreFull(body, title, _scoringMap)
}

func GetWebsiteScoreFull(body string, title string, scoreMap map[string][]string) int {
	var returnScore int = 0
	for scoreIndex, scoreList := range scoreMap {

		scoreMulti, found := rankingToPoints[scoreIndex]

		if !found || len(scoreList) == 0 {
			continue
		}

		returnScore += CountWordsInList(scoreList, body) * scoreMulti
		returnScore += CountWordsInList(scoreList, title) * scoreMulti * titleMulti
	}

	return returnScore

}

func GetPointsFromText(scoring map[string]int, text string) int {
	strs := strings.Split(text, " ")
	var points int = 0
	for _, str := range strs {
		val, ok := scoring[str]
		if !ok {
			continue
		}
		points += val
	}
	return points
}

func UnmarshalScoringInformation(jsonString []byte) map[string][]string {
	var unmarshalTo map[string][]string
	err := json.Unmarshal(jsonString, &unmarshalTo)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return unmarshalTo
}

// Need function that builds regex from json map
// get json map, unmarshal into string map or something
// for each ranking, regex into the string,
// for every occurance you find, multiply that by ranking points
// return that, and add it to the total points
func CountWordsInList(words []string, targetString string) int {
	pattern := fmt.Sprintf(`(?i)\b(%s)\b`, strings.Join(words, "|"))
	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0
	}
	matches := re.FindAllString(targetString, -1)
	return len(matches)
}

func CheckNumberOfSentences(text string) int {
	return strings.Count(text, ".") + strings.Count(text, "?")
}
