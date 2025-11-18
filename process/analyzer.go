package process

import (
	"Paprika/models"
	"Paprika/publisher"
	"context"
	"math"
)

type Analyzer struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	pub       *publisher.Publisher

	ticks map[string]*models.Ticker
}

func SpawnAnalyzerProcess(pub *publisher.Publisher) *Analyzer {
	return &Analyzer{
		ticks: make(map[string]*models.Ticker),
		pub:   pub,
	}
}

func (this *Analyzer) Do(globalCtx context.Context, ch chan *models.Ticker) error {
	ctx, cancel := context.WithCancel(globalCtx)
	this.ctx = ctx
	this.ctxCancel = cancel

	for {
		select {
		case <-this.ctx.Done():
			return nil

		case tick, ok := <-ch:
			if !ok {
				this.Stop()
				return nil
			}
			if tick.Volume < models.MIN_VOLUME {
				continue
			}

			// save
			this.ticks[this.key(tick.Market, tick.Exchange, tick.Symbol)] = tick

			// compare
			spread, ok := this.compare(tick)
			if !ok {
				continue
			}
			if spread.To.Symbol == "XUSDT" {
				continue
			}

			// send to the tg bot
			// log.Printf("SPREAD: FROM %+v\n  TO: %+v\n  Value: %f", *spread.From, *spread.To, spread.Value)
			this.pub.PublishSpread(spread.GetTopic(), spread)
		}
	}
}

func (this *Analyzer) compare(t *models.Ticker) (*models.Spread, bool) {
	// compare with tick from the same exchange, but change market type
	var otherTick *models.Ticker
	var ok bool
	if t.Market == models.FUTURES {
		otherTick, ok = this.ticks[this.key(models.SPOT, t.Exchange, t.Symbol)]
	}
	if t.Market == models.SPOT {
		otherTick, ok = this.ticks[this.key(models.FUTURES, t.Exchange, t.Symbol)]
	}
	if !ok {
		return nil, false
	}

	spread := this.calcSpread(this.calcMidPrice(t.Ask, t.Bid), this.calcMidPrice(otherTick.Ask, otherTick.Bid))
	if math.Abs(spread) <= models.MIN_SPREAD {
		return nil, false // too little spread, return false
	}

	return &models.Spread{From: t, To: otherTick, Value: spread}, true

	// TODO: make for other exchanges
}

func (this *Analyzer) key(exchange, market, symbol string) string {
	return exchange + market + symbol
}

func (this *Analyzer) calcSpread(x, y float64) float64 {
	return math.Abs(x-y) / this.calcMidPrice(x, y)
}

// return a mid price between ask and bid
func (this *Analyzer) calcMidPrice(ask, bid float64) float64 {
	return (ask + bid) / 2
}

func (this *Analyzer) Stop() error {
	this.ctxCancel()
	return nil
}
