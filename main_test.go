package main

import (
	"testing"

	"sharefinder/api"
	"sharefinder/finder"
	"sharefinder/model"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

const benchmarkIterations = 100000

// Define the test cases
var finderTestCases = []struct {
	fundName string
	shares   []string
}{
	{
		fundName: "Ethical Global Fund",
		shares: []string{
			"GreenCo",
			"GrapeCo",
			"GrapeCo",
			"SolarCorp",
			"SpaceY",
			"BeanzRUS",
			"GrapeCo",
			"GoldenGadgets",
			"GrapeCo",
			"SolarCorp",
			"SpaceY",
			"BeanzRUS",
			"GrapeCo",
			"MicroFit",
		},
	},
	{
		fundName: "Fund B",
		shares: []string{
			"MicroFit",
			"GrapeCo",
			"GreenCo",
		},
	},
	{
		fundName: "Fund D",
		shares: []string{
			"SolarCorp",
			"GrapeCo",
			"SpaceY",
			"BeanzRUS",
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

		// Use a diff on the sorted arrays so ordering does not matter
		diff := cmp.Diff(shares, test.shares, cmpopts.SortSlices(func(a, b string) bool { return a < b }))
		assert.Empty(t, diff, "shares should contain same elements in any order")
	}
}

func TestFindingSharesMemento(t *testing.T) {
	// First create the holdings DAG and the memento
	fundData, _ := api.LoadFundData("testdata/example.json")
	holdings := model.NewHoldingsDag(fundData)
	memento := finder.FinderMemento{}
	for _, test := range finderTestCases {
		// Find the shares for the given fund name and check for an error
		shares, err := finder.GetSharesMemento(test.fundName, holdings, memento)
		assert.Nil(t, err, "error should be nil")

		// Use a diff on the sorted arrays so ordering does not matter
		diff := cmp.Diff(shares, test.shares, cmpopts.SortSlices(func(a, b string) bool { return a < b }))
		assert.Empty(t, diff, "shares should contain same elements in any order")
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
	memento := finder.FinderMemento{}
	for i := 0; i < benchmarkIterations; i++ {
		finder.GetSharesMemento("Ethical Global Fund", holdings, memento)
	}
}
