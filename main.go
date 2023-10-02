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
	investments := model.NewInvestmentsDag(fundData)

	// Lets demonstrate basic recursion here
	shares, err := finder.GetSharesRecurse("Ethical Global Fund", investments)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// NOTE
	// Have a look at the unit tests for another approach based on a memoization pattern,
	// which is a common way to tackle dynamic programming problems
	//
	// The benchmarks can be run with `go test -bench . -cpuprofile cpu.prof -count 5`
	// The profiler results can be viewed from `go tool pprof -http=":3000" cpu.prof`
	//
	// The results demonstrate that less time is spent recursing through the data set and we
	// would expect this to scale better as the data set size increases.

	// Print the result to stdout in json format
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(shares)
}
