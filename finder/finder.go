package finder

import (
	"fmt"

	"sharefinder/model"
)

// ShareWeights is a simple Set analog that contains share names against holding weight
type ShareWeights map[string]float64

// FinderMemento is used to store ShareWeights that have previously been enumerated for a fund,
// where the key is the fund name
type FinderMemento map[string]ShareWeights

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
			findSharesRecurse(h.Investment, sharesForHolding)

			// Apply the weight of this holding to the shares we found
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

// findSharesRecurse finds all shares with a basic recursion function but adds a memento for performance of repeat lookups
//
// NOTE
// This is worth exploring because the finding of shares is a set of potentially repeated sub-problems,
// like a dynamic programming problem
//
// NOTE
// It may even be practical to parallelise the sub problems, but let's not get too carried away at this stage :)
// func findSharesMemento(holding *model.Holding, shares ShareWeights, memento FinderMemento) {
// 	for _, h := range holding.Holdings {
// 		if h.IsFund {
// 			// If this is a fund, then we may have already found all of it's shares
// 			storedShares, ok := memento[h.Name]
// 			if ok {
// 				// The shares were already found, just append them
// 				maps.Copy(shares, storedShares)
// 			} else {
// 				// We haven't found the shares for this fund yet, so do it now and store it in the memento
// 				sharesForHolding := make(ShareWeights)
// 				findSharesMemento(h, sharesForHolding, memento)
// 				memento[h.Name] = sharesForHolding

// 				// Now we can use the shares we just found for the current operation
// 				maps.Copy(shares, sharesForHolding)
// 			}
// 		} else {
// 			// This is a share, just add it
// 			shares[h.Name] = struct{}{}
// 		}
// 	}
// }

// // GetSharesMemento demonstrates a dynamic programming approach to the problem
// func GetSharesMemento(primaryHoldingName string, holdings map[string]*model.Holding, memento FinderMemento) (ShareWeights, error) {
// 	primaryHolding, found := holdings[primaryHoldingName]
// 	if !found {
// 		return nil, fmt.Errorf("holding %s was not found", primaryHoldingName)
// 	}

// 	// Wrap our primary holding in a root object; makes it easier to store the result in the memento
// 	rootHolding := model.Holding{Name: "root", Holdings: []*model.Holding{primaryHolding}, IsFund: true}

// 	shares := make(ShareWeights)
// 	findSharesMemento(&rootHolding, shares, memento)
// 	return shares, nil
// }
