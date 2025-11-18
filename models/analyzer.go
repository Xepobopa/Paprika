package models

type Spread struct {
	From  *Ticker
	To    *Ticker
	Value float64 // a value of spread
}

func (s *Spread) GetTopic() string {
	return "cex." + s.From.Market + "." + s.From.Exchange + "." + s.From.Symbol
}
