package finder

import (
	"fmt"

	"sharefinder/model"
)

// finderFunction is a definition for a function that can find shares from the model.Holdings DAG
type finderFunction func(holding *model.Holding, shares *[]string)

// getShares calls the supplied find strategy
func getShares(primaryHoldingName string, holdings map[string]*model.Holding, find finderFunction) ([]string, error) {
	primaryHolding, found := holdings[primaryHoldingName]
	if !found {
		return nil, fmt.Errorf("holding %s was not found", primaryHoldingName)
	}

	shares := []string{}
	find(primaryHolding, &shares)
	return shares, nil
}

// findSharesRecurse finds all shares with a basic recursion function
func findSharesRecurse(holding *model.Holding, shares *[]string) {
	for _, h := range holding.Holdings {
		if h.IsFund {
			findSharesRecurse(h, shares)
		} else {
			*shares = append(*shares, h.Name)
		}
	}
}

// GetSharesRecurse demonstrates the basic recursion strategy
func GetSharesRecurse(primaryHoldingName string, holdings map[string]*model.Holding) ([]string, error) {
	return getShares(primaryHoldingName, holdings, findSharesRecurse)
}
