package exchange

import (
	"Paprika/models"
	"context"

	"github.com/go-resty/resty/v2"
)

type MexcApi struct{}

func NewMexcApi() *MexcApi {
	return &MexcApi{}
}

// FetchContractTicker - futures ticker. Maximum limit - 20 req / 2 sec
func (api *MexcApi) FetchContractTicker(ctx context.Context) (*models.MexcContractTicksRes, error) {
	var res models.MexcContractTicksRes
	_, err := resty.New().
		R().
		SetContext(ctx).
		SetResult(&res).
		Get("https://contract.mexc.com/api/v1/contract/ticker") // api.cfg.Mexc.API.CONTRACTS_TICKERS
	return &res, err
}

// 40 weight. 500 weight per 10 sec. 40 * 10(sec) = 400 weight. Every second
func (api *MexcApi) Fetch24hTickerStats(ctx context.Context) ([]models.MexcTickerStats, error) {
	var res []models.MexcTickerStats
	_, err := resty.New().
		R().
		SetContext(ctx).
		SetResult(&res).
		Get("https://api.mexc.com/api/v3/ticker/24hr") // api.cfg.Mexc.API.SPOT_TICKERS_24HR
	return res, err
}

// MapperMexcToSpot accepts base and quote tokens and returns a specific mexc spot pair string.
//
// BTCUSDT
func MapperMexcSpot(base, quote string) string {
	return base + quote
}

func MapperMexcFutures(base, quote string) string {
	return base + "_" + quote
}
