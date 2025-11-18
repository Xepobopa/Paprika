package exchange

import "context"

type ExchangeApi[ContractTickerRes, SpotTickerRes, MarketInfoRes any] interface {
	FetchContractTicker(ctx context.Context) (*ContractTickerRes, error)
	FetchSpotTicker(ctx context.Context) (*SpotTickerRes, error)
	FetchCurrenciesInfo(ctx context.Context) (*MarketInfoRes, error)
}