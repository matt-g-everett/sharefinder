package model

import (
	"sharefinder/api"
)

// Holding represents a fund or a share holding, when the holding is first discovered, we don't know whether
// its a fund or a share, so this is represented in the IsFund field and can be updated on the fly
//
// NOTE
// We could have chosen to lock down the fund/share definition by looking at the Name to see whether it has
// "Fund" in it. However, this constrains the fund name more and requires text parsing which is more error prone
type Holding struct {
	Name     string
	Holdings []*Holding
	IsFund   bool
}

// ensureHolding returns an existing holding or creates a new one as required
func ensureHolding(name string, isFund bool, holdings map[string]*Holding) *Holding {
	holding, found := holdings[name]
	if !found {
		holding = &Holding{Name: name, IsFund: isFund, Holdings: []*Holding{}}
		holdings[name] = holding
	}

	// Latch to being a fund
	holding.IsFund = holding.IsFund || isFund

	return holding
}

// NewHoldingsDag takes a flat list of fund records and converts them into a DAG (Directed Acyclic Graph)
// of Holdings presented in a map with an entry for each Holding (fund or share) by name
func NewHoldingsDag(fundData api.FundData) map[string]*Holding {
	// Create a map of all the holdings by name
	holdings := make(map[string]*Holding)

	// Find the appropriate attachments for each fundRecord
	for _, fundRecord := range fundData {
		fundHolding := ensureHolding(fundRecord.Name, true, holdings)
		for _, holdingRecord := range fundRecord.Holdings {
			holding := ensureHolding(holdingRecord.Name, false, holdings)
			fundHolding.Holdings = append(fundHolding.Holdings, holding)
		}
	}

	// Technically, we don't need share names in this map, but it's assumed that the value in removing them isn't high
	return holdings
}
