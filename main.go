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
	// We'd probably want to convince ourselves that the maximum DAG depth wouldn't overflow the stack.
	// golang doesn't have comprehensive tail call optimisation and this approach can certainly result in
	// stack overflows in extreme cases. This is easy to show by introducing a loop in example.json
	//
	// We may decide to attempt some form of trampolining, e.g. https://github.com/kandu/go_tailcall
	shares, err := finder.GetSharesRecurse("Ethical Global Fund", holdings)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// NOTE
	// Have a look at the unit tests for another approach based on a memento pattern,
	// which is a common way to tackle dynamic programming like problems
	//
	// The benchmarks can be run with `go test -bench . -cpuprofile cpu.prof -count 5`
	//
	// With the tiny dataset we're using, it doesn't make much difference because the
	// map copy operations counteract the gains at this data scale.
	//
	// However the profiler results from `go tool pprof -http=":3000" cpu.prof`
	// demonstrate that less time is spent recursing through the data set and we
	// would expect this to scale better as the data set size increases
	//
	// Follow-up work could be to assess the scale of real-world datasets and then
	// decide whether an optimisation like the memento pattern should be pursued

	// Convert the ShareSet to a slice, we already know the capacity so preallocate it
	sharesForPrint := make([]string, 0, len(shares))
	for k := range shares {
		sharesForPrint = append(sharesForPrint, k)
	}

	// Print the result to stdout in json format
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(sharesForPrint)
}
