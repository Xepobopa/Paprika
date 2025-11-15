package models

type Spread struct {
	From *Ticker
	To *Ticker
	Value float64 // a value of spread
}