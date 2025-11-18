package process

import (
	"Paprika/exchange"
	"Paprika/models"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrFetchSpotTicker = errors.New("failed to fetch kucoin_spot_tickers")
	ErrFetchFuturesTicker = errors.New("failed to fetch kucoin_futures_tickers")
	ErrNormalizeSpotTickWrongType = errors.New("failed to normalized kucoin spot tick (some field string -> float64)")
	ErrNormalizeFuturesTickWrongType = errors.New("failed to normalized kucoin futures tick (some field string -> float64)")
)

func errWrap(err error, reason error) error {
	return fmt.Errorf("%v, reason: %v", err, reason)
}

type Kucoin struct {
	api exchange.ExchangeApi[
		models.KucoinFuturesTickerRes,
		models.KucoinSpotTickerRes,
		models.KucoinCurrencyInfoRes,
	]

	ctx       context.Context
	ctxCancel context.CancelFunc

	wg sync.WaitGroup
	mu sync.Mutex
}

func SpawnKucoinProcess() *Kucoin {
	return &Kucoin{api: exchange.NewKucoinApi()}
}

func (this *Kucoin) Do(globalCtx context.Context, ch chan *models.Ticker) error {
	this.ctx, this.ctxCancel = context.WithCancel(globalCtx)

	errCh := make(chan error, 2)

	this.wg.Add(2)
	go func() {
		defer this.wg.Done()
		errCh <- this.fetchSpotLoop(ch)
	}()
	go func() {
		defer this.wg.Done()
		errCh <- this.fetchFuturesLoop(ch)
	}()

	go func() {
		this.wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			this.Stop()
			return err
		}
	}

	return nil
}

func (this *Kucoin) fetchSpotLoop(ch chan *models.Ticker) error {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case <-this.ctx.Done():
			return nil

		case <-ticker.C:
			tickers, err := this.api.FetchSpotTicker(this.ctx)
			if err != nil {
				return errWrap(ErrFetchSpotTicker, err)
			}

			for _, tick := range tickers.Data.Ticker {
				normalizedTick, err := this.normalizeSpotTicker(&tick, &tickers.Data.Time)
				if err != nil {
					return errWrap(ErrNormalizeSpotTickWrongType, err)
				}
				
				ch <- normalizedTick
			}
		}
	}
}

func (this *Kucoin) normalizeSpotTicker(tick *models.KucoinSpotTicker, ts *int64) (*models.Ticker, error) {
	askFloat, err := strconv.ParseFloat(tick.Sell, 64) // ask = sell
	if err != nil {
		return nil, err
	}
	bidFloat, err := strconv.ParseFloat(tick.Buy, 64) // bid = buy
	if err != nil {
		return nil, err
	}
	volFloat, err := strconv.ParseFloat(tick.VolValue, 64)
	if err != nil {
		return nil, err
	}

	return &models.Ticker{
		Exchange:     models.KUCOIN,
		Market:       models.SPOT,
		Symbol:       strings.ReplaceAll(tick.Symbol, "-", ""),
		Ask:          askFloat,
		Bid:          bidFloat,
		Volume:       volFloat,
		IsVolumeUsdt: false,
		Timestamp:    *ts,
	}, nil
}

func (this *Kucoin) fetchFuturesLoop(ch chan *models.Ticker) error {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case <-this.ctx.Done():
			return nil

		case <-ticker.C:
			tickers, err := this.api.FetchContractTicker(this.ctx)
			if err != nil {
				return errWrap(ErrFetchFuturesTicker, err)
			}

			for _, tick := range tickers.Data {
				normalizedTick, err := this.normalizeFuturesTicker(&tick)
				if err != nil {
					return errWrap(ErrNormalizeFuturesTickWrongType, err)
				}
				
				ch <- normalizedTick
			}
		}
	}
}

func (this *Kucoin) normalizeFuturesTicker(tick *models.KucoinFuturesTicker) (*models.Ticker, error) {
	askFloat, err := strconv.ParseFloat(tick.BestAskPrice, 64)
	if err != nil {
		return nil, err
	}
	bidFloat, err := strconv.ParseFloat(tick.BestBidPrice, 64)
	if err != nil {
		return nil, err
	}

	return &models.Ticker{
		Exchange:     models.KUCOIN,
		Market:       models.SPOT,
		Symbol:       tick.Symbol,
		Ask:          askFloat,
		Bid:          bidFloat,
		Volume:       0, // kucoin api does not provide volume for futures tickers
		IsVolumeUsdt: false,
		Timestamp:    tick.Ts,
	}, nil
}

func (this *Kucoin) Stop() error {
	this.ctxCancel()
	this.wg.Wait()

	return nil
}
