package model

type Holding struct {
	Name     string
	Holdings []*Holding
	IsFund   bool
}
