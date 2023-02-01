package ctxticker_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-auxiliaries/ctxticker"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func TestContextTickerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &Suite{})
}

func (s *Suite) Test_Basic() {
	signal := getTicks(context.Background(), 0, time.Hour, false, &cbInfo{
		rightAfter: func(ticker *ctxticker.Ticker) {
			ticker.Release()
		},
	})
	s.Equal(1, len(signal))
}

func (s *Suite) Test_Release_FirstTickFast() {
	interval := time.Millisecond * 300
	signal := getTicks(context.Background(), time.Second/2, interval, true, &cbInfo{
		rightAfter: func(ticker *ctxticker.Ticker) {
			ticker.Release()
		},
	})
	s.Equal(2, len(signal))
	s.Less(signal[1].Sub(signal[0]), interval)
}

func (s *Suite) Test_Release() {
	interval := time.Millisecond * 300
	signal := getTicks(context.Background(), time.Second/2, interval, false, &cbInfo{
		rightAfter: func(ticker *ctxticker.Ticker) {
			ticker.Release()
		},
	})
	s.Equal(1, len(signal))
}

func (s *Suite) Test_Trigger() {
	interval := time.Millisecond * 300
	signal := getTicks(context.Background(), time.Second/2, interval, false, &cbInfo{
		rightAfter: func(ticker *ctxticker.Ticker) {
			ticker.Trigger()
		},
	})
	s.Equal(3, len(signal))
	s.Less(signal[1].Sub(signal[0]), interval)
	s.LessOrEqual(interval, signal[2].Sub(signal[1]))
}

func (s *Suite) Test_FirstTickFast() {
	interval := time.Millisecond * 300
	signal := getTicks(context.Background(), time.Second/2, interval, true, nil)
	s.Equal(3, len(signal))
	s.Less(signal[1].Sub(signal[0]), interval)
	s.LessOrEqual(interval, signal[2].Sub(signal[1]))
}

func (s *Suite) TestCancelContext_FirstTickFast1() {
	ctx, cancel := context.WithCancel(context.Background())
	interval := time.Millisecond * 300
	cancel()
	time.Sleep(time.Millisecond * 10)
	signal := getTicks(ctx, time.Second/2, interval, true, nil)
	s.Equal(1, len(signal))
}

func (s *Suite) TestCancelContext_FirstTickFast2() {
	ctx, cancel := context.WithCancel(context.Background())
	interval := time.Millisecond * 300
	signal := getTicks(ctx, time.Second/2, interval, true, &cbInfo{
		rightAfter: func(ticker *ctxticker.Ticker) {
			cancel()
		},
	})
	s.Equal(2, len(signal))
	s.Less(signal[1].Sub(signal[0]), interval)
}

func (s *Suite) TestCancelContext() {
	ctx, cancel := context.WithCancel(context.Background())
	interval := time.Millisecond * 300
	cancel()
	signal := getTicks(ctx, time.Second/2, interval, false, nil)
	s.Equal(1, len(signal))
}

func getTicks(ctx context.Context, timeout, period time.Duration, firstTickFast bool, cb *cbInfo) []time.Time {
	ticker := ctxticker.New(ctx, timeout, period, firstTickFast)
	wg := sync.WaitGroup{}
	wg.Add(1)
	var signal []time.Time
	signal = append(signal, time.Now().UTC())
	started := make(chan bool)
	go func() {
		started <- true
		for {
			if err := ticker.Tick(); err != nil {
				break
			}
			signal = append(signal, time.Now().UTC())
		}
		wg.Done()
	}()
	<-started
	if cb != nil && cb.rightAfter != nil {
		cb.rightAfter(ticker)
	}
	wg.Wait()
	return signal
}

type cbInfo struct {
	rightAfter func(ticker *ctxticker.Ticker)
}
