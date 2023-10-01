package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"sharefinder/api"
	"sharefinder/finder"
	"sharefinder/model"
)

// loadFundData loads api.FundData from a file
func loadFundData(path string) (api.FundData, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	bytes, _ := io.ReadAll(jsonFile)
	return api.NewFundData(bytes), nil
}

func main() {
	// Keep it simple and load directly from a fixed file in a relative location
	// We could load from a path in a command line option, or from stdin or from an API
	fundData, err := loadFundData("testdata/example.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Structure the data into a DAG (directed acyclic graph) so we can traverse the relationships
	holdings := model.NewHoldingsDag(fundData)

	// Lets demonstrate basic recursion to start with
	//
	// NOTE
	// Simple recursion would be ok for small datasets, however, if there was very deeply nested DAG,
	// then we could overflow the stack because golang does not support proper tail calls
	shares, err := finder.GetSharesRecurse("Ethical Global Fund", holdings)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Print the result to stdout in json format
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(shares)
}
