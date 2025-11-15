package process

import (
	"Paprika/exchange"
	"Paprika/models"
	"Paprika/publisher"
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Mexc struct {
	// Process

	api *exchange.MexcApi
	// spotClient *exchange.Client    ??
	// futuresClient *exchange.Client ??

	ctx context.Context
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

func (this *Mexc) Do(globalCtx context.Context, pb *publisher.Publisher) error {
	// create private context
	ctx, cancel := context.WithCancel(globalCtx)
	this.ctx = ctx
	this.ctxCancel = cancel

	errCh := make(chan error, 2)
	this.wg.Add(2)
	go func() {
		defer this.wg.Done()
		errCh <- this.fetchSpotLoop(ctx, pb)
	}()
	go func() {
		defer this.wg.Done()
		errCh <- this.fetchFuturesLoop(ctx, pb)
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

func (this *Mexc) fetchSpotLoop(ctx context.Context, pb *publisher.Publisher) error {
	ticker := time.NewTicker(time.Second)
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
				return fmt.Errorf("failed to make 'FetchTickers' request in MEXC SPOT, reason: %v", err)
			}

			// cast to the unified Ticker
			payloadSpot := make([]models.Ticker, 0, len(tickersSpot))
			for _, tick := range tickersSpot {
				normalizedTick, err := this.normalizeSpotTicker(tick)
				if err != nil {
					return fmt.Errorf("failed to normalize (string -> float64) MEXC SPOT tick: ", err)
				}

				payloadSpot = append(payloadSpot, normalizedTick)
			}

			// publish
			if err := pb.PublishTickers("cex.mexc.spot", &payloadSpot); err != nil {
				return fmt.Errorf("failed to publish mexc spot shapshot: %v", err)
			}
		}
	}
}

func (this *Mexc) normalizeSpotTicker(tick models.MexcTickerStats) (models.Ticker, error) {
	ask, err := strconv.ParseFloat(tick.AskPrice, 64)
	if err != nil {
		return models.Ticker{}, err
	}
	bid, err := strconv.ParseFloat(tick.BidPrice, 64)
	if err != nil {
		return models.Ticker{}, err
	}
	volume, err := strconv.ParseFloat(tick.QuoteVolume, 64)
	if err != nil {
		return models.Ticker{}, err
	}
	
	// todo: cast volume to usdt

	return models.Ticker{
		Exchange:  models.MEXC,
		Market:    models.SPOT,
		Symbol:    tick.Symbol,
		Ask:       ask,
		Bid:       bid,
		Volume:    volume,         // TODO: volume must always be in USDT, but QuoteVolume does not always in USDT, so i need to cast QuoteVolume to the usdt volume
		Timestamp: tick.CloseTime, // closeTime = timestamp
	}, nil
}

func (this *Mexc) fetchFuturesLoop(ctx context.Context, pb *publisher.Publisher) error {
	ticker := time.NewTicker(time.Second)
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

			payloadFutures := make([]models.Ticker, 0, len(tickersFutures.Data))
			for _, tick := range tickersFutures.Data {

				normalizedTick, err := this.normalizeFuturesTicker(tick)
				if err != nil {
					return fmt.Errorf("failed to normalize (string -> float64) MEXC FUTURES tick: ", err)
				}

				payloadFutures = append(payloadFutures, normalizedTick)
			}

			if err := pb.PublishTickers("cex.mexc.futures", &payloadFutures); err != nil {
				return fmt.Errorf("failed to publish mexc futures shapshot: %v", err)
			}
		}
	}
}

func (this *Mexc) normalizeFuturesTicker(tick models.MexcContractTick) (models.Ticker, error) {
	return models.Ticker{
		Exchange:  models.MEXC,
		Market:    models.SPOT,
		Symbol:    tick.Symbol,
		Ask:       tick.Ask1,
		Bid:       tick.Bid1,
		Volume:    tick.Amount24H, // amount is a quote volume (and pray that quote is a usdt)
		Timestamp: tick.Timesamp,
	}, nil
}

func (this *Mexc) Stop() error {
	this.ctxCancel()
	this.wg.Wait()
	return nil
}
