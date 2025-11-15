package models

type MexcContractTicksRes struct {
	Success bool               `json:"success"`
	Code    int                `json:"code"`
	Data    []MexcContractTick `json:"data"`
}

type MexcContractTick struct {
	Symbol    string  `json:"symbol"`
	Bid1      float64 `json:"bid1"`
	Ask1      float64 `json:"ask1"`
	Volume24H float64 `json:"volume24"`
	Amount24H float64 `json:"amount24"`
	Timesamp  int64   `json:"timestamp"`
}

// 24h ticker stats
type MexcTickerStats struct {
	Symbol      string `json:"symbol"`
	BidPrice    string `json:"bidPrice"`
	AskPrice    string `json:"askPrice"`
	Volume      string `json:"volume"`
	QuoteVolume string `json:"quoteVolume"`
	CloseTime   int64  `json:"closeTime"` // like timestamp
}
