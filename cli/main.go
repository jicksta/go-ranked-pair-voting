package main

import (
	"fmt"
	trp "github.com/jicksta/ranked-pairs-voting"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	startTime := time.Now()
	electionFilename := filenameFromArgs(os.Args)
	election := electionFromFile(electionFilename)
	results := election.Results()
	executionDuration := time.Now().Sub(startTime)

	fmt.Print("Results:\n\n")
	for n, group := range results.Winners() {
		fmt.Printf(" %2d. %s\n", n+1, strings.Join(group, " = "))
	}

	fmt.Printf(`
Number of choices: %d
Number of votes:   %d

Algorithm: Tideman Ranked Pairs (TRP)
Time to calculate: %s
Number of cyclical locked pairs: %d`,
		len(election.Choices),
		len(election.Ballots),
		executionDuration,
		len(results.RankedPairs.CyclicalLockedPairsIndices))

	fmt.Print("\n\n\nRanked Pairs Data:\n\n")
	results.RankedPairs.PrintTable(os.Stdout)

	fmt.Print("\n\nTally Data:\n\n")
	results.Tally.PrintTable(os.Stdout)

	fmt.Println(`
The tally data contains information about the 1:1 "runoff" elections that Ranked Pairs simulates, consistent with the
Condorcet Criterion. A table cell with the following text "A=3  B=2  (1)" would indicate that, in a runoff election
between A and B, there are 3 votes for A over B, 2 votes for B over A, and 1 tie. A and B refer to the names listed in
the row and column headers, respectively (see labels).`)

}

func filenameFromArgs(args []string) string {
	if len(args) == 1 {
		log.Fatal("Must supply a filename as a CLI argument")
	} else if len(args) == 2 {
		// Pass through. Go compiler apparently doesn't recognize os.Exit or log.Fatal as terminal like a return
		// statement, so this return must be moved to the bottom of the function.
	} else {
		log.Fatal("Too many params given to the CLI")
	}
	return args[1]
}

func electionFromFile(filename string) *trp.Election {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error: Could not open file at " + filename)
	}
	defer f.Close()
	if election, err := trp.ReadElection(filename, f); err == nil {
		return election
	} else {
		log.Fatal("Error: Unable to process " + filename)
		return nil
	}
}
