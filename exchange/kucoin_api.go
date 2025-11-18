package exchange

import (
	"Paprika/models"
	"context"

	"github.com/go-resty/resty/v2"
)

type KucoinApi struct {
}

// TODO: do not forget to provide api-key
func NewKucoinApi() *KucoinApi {
	return &KucoinApi{}
}

func (this *KucoinApi) FetchContractTicker(ctx context.Context) (*models.KucoinFuturesTickerRes, error) {
	var res models.KucoinFuturesTickerRes
	_, err := resty.New().
		R().
		SetContext(ctx).
		SetResult(&res).
		Get("https://api-futures.kucoin.com/api/v1/allTickers")
	return &res, err
}

func (this *KucoinApi) FetchSpotTicker(ctx context.Context) (*models.KucoinSpotTickerRes, error) {
	var res models.KucoinSpotTickerRes
	_, err := resty.New().
		R().
		SetContext(ctx).
		SetResult(&res).
		Get("https://api.kucoin.com/api/v1/market/allTickers")
	return &res, err
}

func (this *KucoinApi) FetchCurrenciesInfo(ctx context.Context) (*models.KucoinCurrencyInfoRes, error) {
	var res models.KucoinCurrencyInfoRes
	_, err := resty.New().
		R().
		SetContext(ctx).
		SetResult(&res).
		Get("https://api.kucoin.com/api/v3/currencies")
	return &res, err
}
