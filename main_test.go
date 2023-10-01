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

func TestFindingShares(t *testing.T) {
	// First create the holdings DAG
	fundData, _ := api.LoadFundData("testdata/example.json")
	holdings := model.NewHoldingsDag(fundData)

	// Define the test cases
	tests := []struct {
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

	for _, test := range tests {
		// Find the shares for the given fund name and check for an error
		shares, err := finder.GetSharesRecurse(test.fundName, holdings)
		assert.Nil(t, err, "error should be nil")

		// Use a diff on the sorted arrays so ordering does not matter
		diff := cmp.Diff(shares, test.shares, cmpopts.SortSlices(func(a, b string) bool { return a < b }))
		assert.Empty(t, diff, "shares should contain same elements in any order")
	}
}

func BenchmarkRecursion(b *testing.B) {
	for i := 0; i < 1000; i++ {
		fundData, _ := api.LoadFundData("testdata/example.json")
		holdings := model.NewHoldingsDag(fundData)
		finder.GetSharesRecurse("Ethical Global Fund", holdings)
	}
}
