package main

import (
	"encoding/json"
	"fmt"
	"os"

	"sharefinder/api"
	"sharefinder/finder"
	"sharefinder/model"
)

func main() {
	// Keep it simple and load directly from a fixed file in a relative location
	// We could load from a path in a command line option, or from stdin or from an API
	fundData, err := api.LoadFundData("testdata/example.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Structure the data into a DAG (directed acyclic graph) so we can traverse the relationships
	holdings := model.NewHoldingsDag(fundData)

	// Lets demonstrate basic recursion here
	//
	// NOTE
	// We'd probably want to convince ourselves that for a very deeply nested DAG, we wouldn't overflow the stack.
	// golang doesn't have comprehensive tail call optimisation, but some cases it's ok with
	shares, err := finder.GetSharesRecurse("Ethical Global Fund", holdings)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// NOTE
	// Have a look at the unit tests for another approach based on a memento pattern
	// Try running the benchmarks with `go test -bench . -count 3`

	// Print the result to stdout in json format
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(shares)
}
