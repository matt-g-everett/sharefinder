package model

import (
	"sharefinder/api"
)

// Investment represents a fund or share in which we can invest
//
// When the holding is first discovered, we don't know whether its a fund or a share,
// so this is represented in the IsFund field and can be updated on the fly
//
// NOTE
// We could have chosen to lock down the fund/share definition by looking at the Name to see whether it has
// "Fund" in it. However, this constrains the fund name and requires text parsing which is more error prone
type Investment struct {
	Name     string
	IsFund   bool
	Holdings map[string]*Holding
}

// Holding represents the holding of an investment (fund or a share) and it's weight
//
// NOTE
// floats can struggle to represent values exactly, we'll assume for now that this is ok
// and that rounding corrections would happen at the point the ratio was converted back into currency
type Holding struct {
	Investment *Investment
	Weight     float64
}

// InvestmentDag represents a map of Investments that are connected together in a DAG (Directed Acyclid Graph)
type InvestmentsDag map[string]*Investment

// ensureInvestment returns an existing holding or creates a new one as required
func ensureInvestment(name string, isFund bool, investments InvestmentsDag) *Investment {
	investment, found := investments[name]
	if !found {
		investment = &Investment{Name: name, IsFund: isFund, Holdings: make(map[string]*Holding)}
		investments[name] = investment
	}

	// Latch to being a fund
	investment.IsFund = investment.IsFund || isFund

	return investment
}

// NewInvestmentsDag takes a flat list of fund records and converts them into a DAG (Directed Acyclic Graph)
// of Investments and Holdings presented in a map with an entry for each Investment (fund or share) by name
func NewInvestmentsDag(fundData api.FundData) InvestmentsDag {
	// Create a map of all the investments by name
	investments := make(InvestmentsDag)

	// Find the appropriate attachments for each fundRecord
	for _, fundRecord := range fundData {
		fund := ensureInvestment(fundRecord.Name, true, investments)
		for _, holdingRecord := range fundRecord.Holdings {
			investment := ensureInvestment(holdingRecord.Name, false, investments)
			fund.Holdings[investment.Name] = &Holding{Investment: investment, Weight: holdingRecord.Weight}
		}
	}

	// Technically, we don't need share names in this map, but it's assumed that the value in removing them isn't high
	return investments
}
