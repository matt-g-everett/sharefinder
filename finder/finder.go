package finder

import (
	"fmt"
	"maps"

	"sharefinder/model"
)

// ShareSet is a simple Set analog that contains share names
type ShareSet map[string]struct{}

// FinderMemento is used to store ShareSets that have previously been enumerated for a fund,
// where the key is the fund name
type FinderMemento map[string]ShareSet

// findSharesRecurse finds all shares with a basic recursion function
func findSharesRecurse(holding *model.Holding, shares ShareSet) {
	for _, h := range holding.Holdings {
		if h.IsFund {
			findSharesRecurse(h, shares)
		} else {
			shares[h.Name] = struct{}{}
		}
	}
}

// GetSharesRecurse demonstrates the basic recursion strategy
func GetSharesRecurse(primaryHoldingName string, holdings map[string]*model.Holding) (ShareSet, error) {
	primaryHolding, found := holdings[primaryHoldingName]
	if !found {
		return nil, fmt.Errorf("holding %s was not found", primaryHoldingName)
	}

	shares := make(ShareSet)
	findSharesRecurse(primaryHolding, shares)
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
func findSharesMemento(holding *model.Holding, shares ShareSet, memento FinderMemento) {
	for _, h := range holding.Holdings {
		if h.IsFund {
			// If this is a fund, then we may have already found all of it's shares
			storedShares, ok := memento[h.Name]
			if ok {
				// The shares were already found, just append them
				maps.Copy(shares, storedShares)
			} else {
				// We haven't found the shares for this fund yet, so do it now and store it in the memento
				sharesForHolding := make(ShareSet)
				findSharesMemento(h, sharesForHolding, memento)
				memento[h.Name] = sharesForHolding

				// Now we can use the shares we just found for the current operation
				maps.Copy(shares, sharesForHolding)
			}
		} else {
			// This is a share, just add it
			shares[h.Name] = struct{}{}
		}
	}
}

// GetSharesMemento demonstrates a dynamic programming approach to the problem
func GetSharesMemento(primaryHoldingName string, holdings map[string]*model.Holding, memento FinderMemento) (ShareSet, error) {
	primaryHolding, found := holdings[primaryHoldingName]
	if !found {
		return nil, fmt.Errorf("holding %s was not found", primaryHoldingName)
	}

	// Wrap our primary holding in a root object; makes it easier to store the result in the memento
	rootHolding := model.Holding{Name: "root", Holdings: []*model.Holding{primaryHolding}, IsFund: true}

	shares := make(ShareSet)
	findSharesMemento(&rootHolding, shares, memento)
	return shares, nil
}
