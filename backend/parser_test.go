package main

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	fmt.Println("test running")
	//testJson := UnmarshalScoringInformation([]byte(`
	//	{
	//		"D" : [ "disaster" ]
	//	}
	//	`))
	// Reliant on global json state, which ain't ideal
	pts := GetWebsiteScore("test, test, DISASTER, disASTER, disaster", "non-title")

	shouldBePts := rankingToPoints["D"] * 3

	if shouldBePts != pts {
		t.Errorf("Parser points wrong, output %d should be %d", pts, shouldBePts)
	}
}

func TestParserAdvanced(t *testing.T) {
	fmt.Println("test running")

	testJson := UnmarshalScoringInformation([]byte(`
		{
			"D" : [ "disaster" ],
			"H" : [ "test" ],
			"N" : [ "title" ]
		}
		`))

	// Reliant on global json state, which ain't ideal
	pts := GetWebsiteScoreFull("test, test, DISASTER, disASTER, disaster", "non-title", testJson)

	shouldBePts := rankingToPoints["D"]*3 +
		rankingToPoints["H"]*2 + rankingToPoints["N"]*titleMulti

	if shouldBePts != pts {
		t.Errorf("Parser points wrong, output %d should be %d", pts, shouldBePts)
	}
}
