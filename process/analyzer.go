package process

import (
	"Paprika/models"
	"Paprika/publisher"
	"context"
	"log"
	"math"
)

type Analyzer struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	ticks map[string]*models.Ticker

	ch chan *models.Ticker
}

func SpawnAnalyzerProcess() *Analyzer {
	return &Analyzer{
		ch: make(chan *models.Ticker, 5000),
	}
}

func (this *Analyzer) Do(globalCtx context.Context, pb *publisher.Publisher) error {
	ctx, cancel := context.WithCancel(globalCtx)
	this.ctx = ctx
	this.ctxCancel = cancel

	for {
		select {
		case <-this.ctx.Done():
			return nil

		case tick, ok := <-this.ch:
			if !ok {
				this.Stop()
				return nil
			}
			if tick.Volume < models.MIN_VOLUME {
				continue
			}
			
			// save
			this.ticks[this.key(tick.Exchange, tick.Market, tick.Symbol)] = tick

			// compare
			spread, ok := this.compare(tick)
			if !ok {
				continue
			}
			// send to the tg bot
			log.Printf("SPREAD: %+v", spread)
		}
	}
}

func (this *Analyzer) compare(t *models.Ticker) (*models.Spread, bool) {
	// compare with tick from the same exchange, but change market type
	var otherTick *models.Ticker
	var ok bool
	if t.Market == models.FUTURES {
		otherTick, ok = this.ticks[this.key(t.Exchange, models.SPOT, t.Symbol)]
	}
	if t.Market == models.SPOT {
		otherTick, ok = this.ticks[this.key(t.Exchange, models.FUTURES, t.Symbol)]
	}

	if !ok {
		return nil, false
	}
	spread := this.calcMidPrice(t.Ask, t.Bid) - this.calcMidPrice(otherTick.Ask, otherTick.Bid)
	if math.Abs(spread) <= models.MIN_SPREAD {
		return nil, false // too little spread, return false
	}

	return &models.Spread{From: t, To: otherTick, Value: spread}, true

	// TODO: make for other exchanges
}

func (this *Analyzer) key(exchange, market, symbol string) string {
	return exchange + market + symbol
}

// return a mid price between ask and bid
func (this *Analyzer) calcMidPrice(ask, bid float64) float64 {
	return (ask + bid) / 2
}

func (this *Analyzer) Stop() error {
	this.ctxCancel()
	return nil
}
