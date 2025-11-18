package process

import (
	"Paprika/exchange"
	"Paprika/models"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Mexc struct {
	// Process

	api *exchange.MexcApi
	// spotClient *exchange.Client    ??
	// futuresClient *exchange.Client ??

	ctx       context.Context
	ctxCancel context.CancelFunc

	wg sync.WaitGroup
	mu sync.Mutex
}

// client with mexc
func SpawnMexcProcess() *Mexc {
	api := exchange.NewMexcApi()
	return &Mexc{
		// // every second
		// spotClient: exchange.NewClient(60, 1, time.Second * 5),
		// // every second
		// futuresClient: exchange.NewClient(60, 1, time.Second * 5),

		api: api,
	}
}

func (this *Mexc) Do(globalCtx context.Context, ch chan *models.Ticker) error {
	// create private context
	ctx, cancel := context.WithCancel(globalCtx)
	this.ctx = ctx
	this.ctxCancel = cancel

	errCh := make(chan error, 2)
	this.wg.Add(2)
	go func() {
		defer this.wg.Done()
		errCh <- this.fetchSpotLoop(this.ctx, ch)
	}()
	go func() {
		defer this.wg.Done()
		errCh <- this.fetchFuturesLoop(this.ctx, ch)
	}()

	// Wait for both goroutines to finish
	go func() {
		this.wg.Wait()
		close(errCh)
	}()

	// return first non-nil error
	for err := range errCh {
		if err != nil {
			// stop all and wait till stop
			this.Stop()
			return err
		}
	}

	return nil
}

func (this *Mexc) fetchSpotLoop(ctx context.Context, ch chan *models.Ticker) error {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	

	for {
		select {
		case <-this.ctx.Done():
			return nil

		case <-ticker.C:
			// start to goroutines

			// without options to fetch all tickers.
			// SPOT: NOTE: refer to mexc docs,
			// the /api/v3/ticker/24hr endpoint costs 40 ip weight, mexc gives 500
			// for every 10 seconds, so 40 (weight per req) * 10 (sec) = 400 points / 500.
			// So it lefts only 100 points for other requests. Keep this in mind
			tickersSpot, err := this.api.Fetch24hTickerStats(ctx)
			if err != nil {
				// stop process and return the error
				log.Println("%+V", err)
				return fmt.Errorf("failed to make 'FetchTickers' request in MEXC SPOT, reason: %v", err)
			}

			// cast to the unified Ticker
			for _, tick := range tickersSpot {
				normalizedTick, err := this.normalizeSpotTicker(tick)
				if err != nil {
					return fmt.Errorf("failed to normalize (string -> float64) MEXC SPOT tick: %v", err)
				}

				if normalizedTick.Ask == 0 || normalizedTick.Bid == 0 || normalizedTick.Volume == 0 {
					continue
				}

				// send to the channels
				ch <- normalizedTick
			}
		}
	}
}

func (this *Mexc) normalizeSpotTicker(tick models.MexcTickerStats) (*models.Ticker, error) {
	ask, err := strconv.ParseFloat(tick.AskPrice, 64)
	if err != nil {
		return nil, err
	}
	bid, err := strconv.ParseFloat(tick.BidPrice, 64)
	if err != nil {
		return nil, err
	}
	volume, err := strconv.ParseFloat(tick.QuoteVolume, 64)
	if err != nil {
		return nil, err
	}

	// todo: cast volume to usdt

	return &models.Ticker{
		Exchange:  models.MEXC,
		Market:    models.SPOT,
		Symbol:    tick.Symbol,
		Ask:       ask,
		Bid:       bid,
		Volume:    volume,         // TODO: volume must always be in USDT, but QuoteVolume does not always in USDT, so i need to cast QuoteVolume to the usdt volume
		Timestamp: tick.CloseTime, // closeTime = timestamp
	}, nil
}

func (this *Mexc) fetchFuturesLoop(ctx context.Context, ch chan *models.Ticker) error {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case <-this.ctx.Done():
			return nil

		case <-ticker.C:
			// FUTURES: NOTE: refer to the mexc docs, the /api/v1/contract/ticker
			// endpoint has a limit 20req / 2 sec
			tickersFutures, err := this.api.FetchContractTicker(ctx)
			if err != nil {
				return fmt.Errorf("failed to make 'FetchTickers' request in MEXC FUTURES (SWAP), reason: %v", err)
			}

			for _, tick := range tickersFutures.Data {
				normalizedTick, err := this.normalizeFuturesTicker(tick)
				if err != nil {
					return fmt.Errorf("failed to normalize (string -> float64) MEXC FUTURES tick: %v", err)
				}

				if normalizedTick.Ask == 0 || normalizedTick.Bid == 0 || normalizedTick.Volume == 0 {
					continue
				}
				
				// send to the channel
				ch <- normalizedTick
			}

		}
	}
}

func (this *Mexc) normalizeFuturesTicker(tick models.MexcContractTick) (*models.Ticker, error) {
	return &models.Ticker{
		Exchange:  models.MEXC,
		Market:    models.FUTURES,
		Symbol:    strings.ReplaceAll(tick.Symbol, "_", ""), // BTC_USDT -> BTCUSDT
		Ask:       tick.Ask1,
		Bid:       tick.Bid1,
		Volume:    tick.Amount24H, // amount is a quote volume (not always in usdt)
		Timestamp: tick.Timesamp,
	}, nil
}

func (this *Mexc) Stop() error {
	this.ctxCancel()
	this.wg.Wait()

	return nil
}
