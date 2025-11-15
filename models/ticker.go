package models

type Ticker struct {
	Exchange     string  `json:"exchange"` // mex / binance / ...
	Market       string  `json:"market"`   // spot / futures
	Symbol       string  `json:"symbol"`
	Ask          float64 `json:"ask"`    // ask1 for futures and bestAsk for spot
	Bid          float64 `json:"bid"`    // bid1 for futures and bestBid for spot
	Volume       float64 `json:"volume"` // 24h
	IsVolumeUsdt bool    `json:"isVolumeUsdt"`
	Timestamp    int64   `json:"timestamp"` // when snapshot was received (closeTime for spot and timestamp for futures)
}
