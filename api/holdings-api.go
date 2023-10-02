package api

import (
	"encoding/json"
	"io"
	"os"
)

// HoldingRecord represents a share or fund and the ratio of the parent fund that it makes up
type HoldingRecord struct {
	Name   string  `json:"name"`
	Weight float32 `json:"weight"`
}

// FundRecord represent a fund definition including a list of holdings
type FundRecord struct {
	Name     string          `json:"name"`
	Holdings []HoldingRecord `json:"holdings"`
}

// FundData represents the document containing FundRecords
type FundData []FundRecord

// NewFundData creates a FundData object from a slice of bytes, which could be read from a file
func NewFundData(bytes []byte) FundData {
	fundData := make(FundData, 0)
	json.Unmarshal(bytes, &fundData)

	return fundData
}

// loadFundData loads api.FundData from a file
func LoadFundData(path string) (FundData, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	bytes, _ := io.ReadAll(jsonFile)
	return NewFundData(bytes), nil
}
