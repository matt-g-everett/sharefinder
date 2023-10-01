package finder

import (
	"fmt"

	"sharefinder/model"
)

type FinderMemento map[string][]string

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
	primaryHolding, found := holdings[primaryHoldingName]
	if !found {
		return nil, fmt.Errorf("holding %s was not found", primaryHoldingName)
	}

	shares := []string{}
	findSharesRecurse(primaryHolding, &shares)
	return shares, nil
}

// findSharesRecurse finds all shares with a basic recursion function but adds a memento for performance of repeat lookups
//
// NOTE
// This is worth exploring because the finding of shares is a set of potentially repeated sub-problems,
// like a dynamic programming problem
//
// NOTE
// It may even be practical to parallelise the sub problems, but let's not get too carried away at this stage :)
func findSharesMemento(holding *model.Holding, shares *[]string, memento FinderMemento) {
	for _, h := range holding.Holdings {
		if h.IsFund {
			// If this is a fund, then we may have already found all of it's shares
			storedShares, ok := memento[h.Name]
			if ok {
				// The shares were already found, just append them
				*shares = append(*shares, storedShares...)
			} else {
				// We haven't found the shares for this fund yet, so do it now and store it in the memento
				sharesForHolding := []string{}
				findSharesMemento(h, &sharesForHolding, memento)
				memento[h.Name] = sharesForHolding

				// Now we can use the shares we just found for the current operation
				*shares = append(*shares, sharesForHolding...)
			}
		} else {
			// This is a share, just add it
			*shares = append(*shares, h.Name)
		}
	}
}

// GetSharesMemento demonstrates a dynamic programming approach to the problem
func GetSharesMemento(primaryHoldingName string, holdings map[string]*model.Holding, memento FinderMemento) ([]string, error) {
	primaryHolding, found := holdings[primaryHoldingName]
	if !found {
		return nil, fmt.Errorf("holding %s was not found", primaryHoldingName)
	}

	// Wrap our primary holding in a root object; makes it easier to store the result in the memento
	rootHolding := model.Holding{Name: "root", Holdings: []*model.Holding{primaryHolding}, IsFund: true}

	shares := []string{}
	findSharesMemento(&rootHolding, &shares, memento)
	return shares, nil
}
