package finder

import (
	"fmt"

	"sharefinder/model"
)

// ShareWeights is a simple Set analog that contains share names against holding weight
type ShareWeights map[string]float64

// FinderMemoir is used to store ShareWeights that have previously been enumerated for a fund,
// where the key is the fund name
type FinderMemoir map[string]ShareWeights

func addExposure(shareName string, exposure float64, shares ShareWeights) {
	// If we found a share, make sure the share map is initialised for it
	if _, found := shares[shareName]; !found {
		shares[shareName] = 0.0
	}

	// Add the exposure to the shares map
	shares[shareName] += exposure
}

// findSharesRecurse finds all shares with a basic recursion function
func findSharesRecurse(investment *model.Investment, shares ShareWeights) {
	for _, h := range investment.Holdings {
		if h.Investment.IsFund {
			sharesForHolding := make(ShareWeights)

			// NOTE
			// Recursing in the middle of a function makes tail call optimisation (probably) impossible.
			// However, golang doesn't really have comprehensive tail call optimisation anyway, so if
			// we did want to protect ourselves from stack overflows, we'd probably have to do some
			// trampolining anyway ... we'll call that a future optimisation for now!
			findSharesRecurse(h.Investment, sharesForHolding)

			// Apply the weight of this holding to the shares we've already found
			for shareName, weight := range sharesForHolding {
				shareExposure := h.Weight * weight
				addExposure(shareName, shareExposure, shares)
			}
		} else {
			addExposure(h.Investment.Name, h.Weight, shares)
		}
	}
}

// GetSharesRecurse demonstrates the basic recursion strategy
func GetSharesRecurse(primaryHoldingName string, investments model.InvestmentsDag) (ShareWeights, error) {
	primaryHolding, found := investments[primaryHoldingName]
	if !found {
		return nil, fmt.Errorf("holding %s was not found", primaryHoldingName)
	}

	shares := make(ShareWeights)
	findSharesRecurse(primaryHolding, shares)
	return shares, nil
}

// findSharesRecurse finds all shares with a basic recursion function but adds memoization for performance of repeat lookups
//
// NOTE
// This is worth exploring because the finding of shares is a set of potentially repeated sub-problems.
// There's a passing resemblance to a Fibonacci sequence here and we can probably treat it as a dynamic programming problem.
// Memoization is a simple tool in dynamic programming solutions and works nicely if the DAG won't always be fully interrogated.
// Since we're enumerating top-down from a single fund name, memoization seems like a good choice.
//
// NOTE
// It may even be practical to parallelise the sub problems, but let's not get too carried away at this stage :)
// If we did parallelise the subproblems, we may have to protect ourselves from goroutine overuse too,
// maybe with a worker pool implementation
func findSharesMemoized(investment *model.Investment, shares ShareWeights, memoir FinderMemoir) {
	for _, h := range investment.Holdings {
		if h.Investment.IsFund {
			// If this is a fund, then we may have already found all of it's shares
			childShares, ok := memoir[h.Investment.Name]
			if !ok {
				// We haven't found the shares for this fund yet, so do it now and store it in the memoir
				childShares = make(ShareWeights)
				findSharesMemoized(h.Investment, childShares, memoir)
				memoir[h.Investment.Name] = childShares
			}

			// Add the weight of the stored shares to the ones we've already found
			for shareName, weight := range childShares {
				shareExposure := h.Weight * weight
				addExposure(shareName, shareExposure, shares)
			}
		} else {
			addExposure(h.Investment.Name, h.Weight, shares)
		}
	}
}

// GetSharesMemoized demonstrates an approach using memoization to show how an optimisation could be applied
func GetSharesMemoized(primaryHoldingName string, investments model.InvestmentsDag, memoir FinderMemoir) (ShareWeights, error) {
	primaryInvestment, found := investments[primaryHoldingName]
	if !found {
		return nil, fmt.Errorf("holding %s was not found", primaryHoldingName)
	}

	// Wrap our primary holding in a root object; makes it easier to store the result in the memoir
	primaryHolding := model.Holding{Investment: primaryInvestment, Weight: 1.0}
	rootInvestment := model.Investment{Name: "root", Holdings: map[string]*model.Holding{primaryHoldingName: &primaryHolding}, IsFund: true}

	shares := make(ShareWeights)
	findSharesMemoized(&rootInvestment, shares, memoir)
	return shares, nil
}
