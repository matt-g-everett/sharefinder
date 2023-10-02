package main

import (
	"fmt"
	"math"
	"testing"

	"sharefinder/api"
	"sharefinder/finder"
	"sharefinder/model"

	"github.com/stretchr/testify/assert"
)

const benchmarkIterations = 10000

// epsilon represents the smallest increment for a float close to zero
var epsilon = math.Nextafter(1, 2) - 1

// compareShareWeights determines whether two sets of ShareWeights are similar enough
func compareShareWeights(t *testing.T, control finder.ShareWeights, result finder.ShareWeights) {
	assert.Equal(t, len(control), len(result), "should be the same number of shares")
	for controlName, controlWeight := range control {
		resultWeight, found := result[controlName]
		assert.True(t, found, fmt.Sprintf("share %s should be present", controlName))

		// NOTE
		// Floats are a bit tricky to compare, this is a common way to do it
		weightDiff := math.Abs(controlWeight - resultWeight)
		assert.LessOrEqual(t, weightDiff, epsilon, fmt.Sprintf("share %s weight of %f should be close to %f", controlName, resultWeight, controlWeight))
	}
}

func checkSumOfWeights(t *testing.T, result finder.ShareWeights) {
	sum := 0.0
	for _, weight := range result {
		sum += weight
	}

	assert.LessOrEqual(t, sum-1.0, epsilon, "the sum of the weights should be close to 1.0")
}

// Define the test cases
var finderTestCases = []struct {
	fundName string
	shares   finder.ShareWeights
}{
	{
		fundName: "Ethical Global Fund",
		shares: finder.ShareWeights{
			"GreenCo":       0.06,
			"SolarCorp":     0.028,
			"SpaceY":        0.105,
			"BeanzRUS":      0.21,
			"GoldenGadgets": 0.15,
			"MicroFit":      0.1,
			"GrapeCo":       0.347,
		},
	},
	{
		fundName: "Fund B",
		shares: finder.ShareWeights{
			"GreenCo":  0.3,
			"GrapeCo":  0.2,
			"MicroFit": 0.5,
		},
	},
	{
		fundName: "Fund D",
		shares: finder.ShareWeights{
			"SolarCorp": 0.08,
			"GrapeCo":   0.02,
			"SpaceY":    0.3,
			"BeanzRUS":  0.6,
		},
	},
}

func TestFindingShares(t *testing.T) {
	// First create the holdings DAG
	fundData, _ := api.LoadFundData("testdata/example.json")
	investments := model.NewInvestmentsDag(fundData)

	for _, test := range finderTestCases {
		// Find the shares for the given fund name and check for an error
		shares, err := finder.GetSharesRecurse(test.fundName, investments)
		assert.Nil(t, err, "error should be nil")
		checkSumOfWeights(t, shares)
		compareShareWeights(t, test.shares, shares)
	}
}

func TestFindingSharesMemoized(t *testing.T) {
	// First create the holdings DAG and the memoir
	fundData, _ := api.LoadFundData("testdata/example.json")
	holdings := model.NewInvestmentsDag(fundData)
	memoir := make(finder.FinderMemoir)
	for _, test := range finderTestCases {
		// Find the shares for the given fund name and check for an error
		shares, err := finder.GetSharesMemoized(test.fundName, holdings, memoir)
		assert.Nil(t, err, "error should be nil")
		checkSumOfWeights(t, shares)
		compareShareWeights(t, test.shares, shares)
	}
}

func BenchmarkRecursion(b *testing.B) {
	fundData, _ := api.LoadFundData("testdata/example.json")
	investments := model.NewInvestmentsDag(fundData)
	for i := 0; i < benchmarkIterations; i++ {
		finder.GetSharesRecurse("Ethical Global Fund", investments)
	}
}

func BenchmarkMemoized(b *testing.B) {
	fundData, _ := api.LoadFundData("testdata/example.json")
	holdings := model.NewInvestmentsDag(fundData)
	memoir := make(finder.FinderMemoir)
	for i := 0; i < benchmarkIterations; i++ {
		finder.GetSharesMemoized("Ethical Global Fund", holdings, memoir)
	}
}
