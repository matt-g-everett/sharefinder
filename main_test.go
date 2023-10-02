package main

import (
	"testing"

	"sharefinder/api"
	"sharefinder/finder"
	"sharefinder/model"

	"github.com/stretchr/testify/assert"
)

const benchmarkIterations = 100000

// Define the test cases
var finderTestCases = []struct {
	fundName string
	shares   finder.ShareSet
}{
	{
		fundName: "Ethical Global Fund",
		shares: finder.ShareSet{
			"GreenCo":       struct{}{},
			"SolarCorp":     struct{}{},
			"SpaceY":        struct{}{},
			"BeanzRUS":      struct{}{},
			"GoldenGadgets": struct{}{},
			"MicroFit":      struct{}{},
			"GrapeCo":       struct{}{},
		},
	},
	{
		fundName: "Fund B",
		shares: finder.ShareSet{
			"GreenCo":  struct{}{},
			"GrapeCo":  struct{}{},
			"MicroFit": struct{}{},
		},
	},
	{
		fundName: "Fund D",
		shares: finder.ShareSet{
			"SolarCorp": struct{}{},
			"GrapeCo":   struct{}{},
			"SpaceY":    struct{}{},
			"BeanzRUS":  struct{}{},
		},
	},
}

func TestFindingShares(t *testing.T) {
	// First create the holdings DAG
	fundData, _ := api.LoadFundData("testdata/example.json")
	holdings := model.NewHoldingsDag(fundData)

	for _, test := range finderTestCases {
		// Find the shares for the given fund name and check for an error
		shares, err := finder.GetSharesRecurse(test.fundName, holdings)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, shares, test.shares, "shares should contain same elements")
	}
}

func TestFindingSharesMemento(t *testing.T) {
	// First create the holdings DAG and the memento
	fundData, _ := api.LoadFundData("testdata/example.json")
	holdings := model.NewHoldingsDag(fundData)
	memento := make(finder.FinderMemento)
	for _, test := range finderTestCases {
		// Find the shares for the given fund name and check for an error
		shares, err := finder.GetSharesMemento(test.fundName, holdings, memento)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, shares, test.shares, "shares should contain same elements")
	}
}

func BenchmarkRecursion(b *testing.B) {
	fundData, _ := api.LoadFundData("testdata/example.json")
	holdings := model.NewHoldingsDag(fundData)
	for i := 0; i < benchmarkIterations; i++ {
		finder.GetSharesRecurse("Ethical Global Fund", holdings)
	}
}

func BenchmarkMemento(b *testing.B) {
	fundData, _ := api.LoadFundData("testdata/example.json")
	holdings := model.NewHoldingsDag(fundData)
	memento := make(finder.FinderMemento)
	for i := 0; i < benchmarkIterations; i++ {
		finder.GetSharesMemento("Ethical Global Fund", holdings, memento)
	}
}
