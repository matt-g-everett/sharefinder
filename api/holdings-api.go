package api

type HoldingRecord struct {
	Name   string  `json:"name"`
	Weight float32 `json:"weight"`
}

type FundRecord struct {
	Name     string          `json:"name"`
	Holdings []HoldingRecord `json:"holdings"`
}
