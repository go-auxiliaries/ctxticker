package ctxticker

import (
	"context"
	"time"
)

type Ticker struct {
	Cancel  func()
	Ctx     context.Context
	Timeout time.Duration
	Period  time.Duration
	ticker  *time.Ticker
	release chan bool
}

func New(ctx context.Context, timeout, period time.Duration, firstTickFast bool) *Ticker {
	var tickerCtx context.Context
	var cancel func()
	if timeout == time.Duration(0) {
		tickerCtx, cancel = context.WithCancel(ctx)
	} else {
		tickerCtx, cancel = context.WithTimeout(ctx, timeout)
	}
	releaseChannel := make(chan bool, 2)
	if firstTickFast {
		releaseChannel <- false // Trigger it to tick for the first time
	}
	return &Ticker{
		Cancel:  cancel,
		Ctx:     tickerCtx,
		Timeout: timeout,
		Period:  period,
		ticker:  time.NewTicker(period),
		release: releaseChannel,
	}
}

func (ct *Ticker) Tick() error {
	select {
	case <-ct.Ctx.Done():
		return context.Canceled
	default:
	}
	select {
	case <-ct.Ctx.Done():
		return context.Canceled
	case val := <-ct.release:
		if val {
			return context.Canceled
		}
		return nil
	case <-ct.ticker.C:
		return nil
	}
}

func (ct *Ticker) GetTickContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(ct.Ctx, ct.Period)
}

func (ct *Ticker) Stop() {
	ct.Cancel()
	ct.ticker.Stop()
}

func (ct *Ticker) Release() {
	select {
	case ct.release <- true:
	default:
	}
}

func (ct *Ticker) Trigger() {
	select {
	case ct.release <- false:
	default:
	}
}
