package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"sharefinder/api"
	"sharefinder/model"
)

// ensureHolding returns an existing holding or creates a new one to return as required
func ensureHolding(name string, isFund bool, holdings map[string]*model.Holding) *model.Holding {
	holding, found := holdings[name]
	if !found {
		holding = &model.Holding{Name: name, IsFund: isFund, Holdings: []*model.Holding{}}
		holdings[name] = holding
	}

	// Latch to being a fund
	holding.IsFund = holding.IsFund || isFund

	return holding
}

// processRecords takes a flat list of fund records and converts them into
func processRecords(fundRecords []api.FundRecord) map[string]*model.Holding {
	// Create a map of all the holdings by name
	holdings := map[string]*model.Holding{}

	// Find the appropriate attachments for each fundRecord
	for _, fundRecord := range fundRecords {
		fundHolding := ensureHolding(fundRecord.Name, true, holdings)
		for _, holdingRecord := range fundRecord.Holdings {
			holding := ensureHolding(holdingRecord.Name, false, holdings)
			fundHolding.Holdings = append(fundHolding.Holdings, holding)
		}
	}

	return holdings
}

func findShares(holding *model.Holding, shares *[]string) {
	for _, h := range holding.Holdings {
		if h.IsFund {
			findShares(h, shares)
		} else {
			*shares = append(*shares, h.Name)
		}
	}
}

func getShares(primaryHoldingName string, holdings map[string]*model.Holding) ([]string, error) {
	primaryHolding, found := holdings[primaryHoldingName]
	if !found {
		return nil, fmt.Errorf("holding %s was not found", primaryHoldingName)
	}

	shares := []string{}
	findShares(primaryHolding, &shares)
	return shares, nil
}

func main() {
	// Read from a fixed file
	jsonFile, err := os.Open("data/example.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer jsonFile.Close()

	// Read the file data into memory
	var fundRecords []api.FundRecord
	bytes, _ := io.ReadAll(jsonFile)
	json.Unmarshal(bytes, &fundRecords)

	// Structure the data so we can traverse the relationships
	holdings := processRecords(fundRecords)

	// Now we can recurse through the holdings and find the shares
	shares, err := getShares("Ethical Global Fund", holdings)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Print the result to stdout in json format
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(shares)
}
