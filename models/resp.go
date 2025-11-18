package models

// MEXC
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

// KUCOIN
type KucoinCurrencyInfoRes struct {
	Code string `json:"code"`
	Data []struct {
		Currency        string                    `json:"currency"`
		Name            string                    `json:"name"`
		FullName        string                    `json:"fullName"`
		Precision       int                       `json:"precision"`
		Confirms        *int                      `json:"confirms"`        // nullable
		ContractAddress *string                   `json:"contractAddress"` // nullable
		IsMarginEnabled bool                      `json:"isMarginEnabled"`
		IsDebitEnabled  bool                      `json:"isDebitEnabled"`
		Chains          []KucoinCurrencyInfoChain `json:"chains"`
	} `json:"data"`
}

type KucoinCurrencyInfoChain struct {
	ChainName         string  `json:"chainName"`
	WithdrawalMinSize string  `json:"withdrawalMinSize"`
	DepositMinSize    *string `json:"depositMinSize"` // can be null
	WithdrawFeeRate   string  `json:"withdrawFeeRate"`
	WithdrawalMinFee  string  `json:"withdrawalMinFee"`
	IsWithdrawEnabled bool    `json:"isWithdrawEnabled"`
	IsDepositEnabled  bool    `json:"isDepositEnabled"`
	Confirms          int     `json:"confirms"`
	PreConfirms       int     `json:"preConfirms"`
	ContractAddress   string  `json:"contractAddress"`
	WithdrawPrecision int     `json:"withdrawPrecision"`
	MaxWithdraw       *string `json:"maxWithdraw"` // nullable
	MaxDeposit        *string `json:"maxDeposit"`  // nullable
	NeedTag           bool    `json:"needTag"`
	ChainId           string  `json:"chainId"`
}

type KucoinSpotTickerRes struct {
	Code string `json:"code"`
	Data struct {
		Time   int64    `json:"time"` // Unix millisecond timestamp
		Ticker []KucoinSpotTicker `json:"ticker"`
	} `json:"data"`
}

type KucoinSpotTicker struct {
	Symbol           string `json:"symbol"`     // e.g. "BTC-USDT"
	SymbolName       string `json:"symbolName"` // "BTC-USDT"
	Buy              string `json:"buy"`        // best bid price
	BestBidSize      string `json:"bestBidSize"`
	Sell             string `json:"sell"` // best ask price
	BestAskSize      string `json:"bestAskSize"`
	ChangeRate       string `json:"changeRate"`   // 24h change rate (e.g. "-0.0014")
	ChangePrice      string `json:"changePrice"`  // 24h price change
	High             string `json:"high"`         // 24h high
	Low              string `json:"low"`          // 24h low
	Vol              string `json:"vol"`          // 24h volume (in base currency)
	VolValue         string `json:"volValue"`     // 24h volume (in quote currency)
	Last             string `json:"last"`         // latest price
	AveragePrice     string `json:"averagePrice"` // 24h average price
	TakerFeeRate     string `json:"takerFeeRate"`
	MakerFeeRate     string `json:"makerFeeRate"`
	TakerCoefficient string `json:"takerCoefficient"`
	MakerCoefficient string `json:"makerCoefficient"`
}

type KucoinFuturesTickerRes struct {
	Code string                `json:"code"`
	Data []KucoinFuturesTicker `json:"data"`
}

type KucoinFuturesTicker struct {
	Sequence     int64  `json:"sequence"`     // Order book sequence number (int64!)
	Symbol       string `json:"symbol"`       // e.g. XBTUSDTM, ETHUSDTM
	Side         string `json:"side"`         // "buy" or "sell"
	Size         int64  `json:"size"`         // Trade size in contracts (integer)
	TradeID      string `json:"tradeId"`      // Trade ID
	Price        string `json:"price"`        // Last trade price → string for precision
	BestBidPrice string `json:"bestBidPrice"` // Best bid at the moment
	BestBidSize  int64  `json:"bestBidSize"`  // Size at best bid (contracts)
	BestAskPrice string `json:"bestAskPrice"` // Best ask
	BestAskSize  int64  `json:"bestAskSize"`  // Size at best ask
	Ts           int64  `json:"ts"`           // Timestamp in nanoseconds (int64)
}
