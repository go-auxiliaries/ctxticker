## Description ##

This is simple wrapper on `time.Ticker` that takes into consideration `context.Context` state

## Examples ##

### 1. Service example ###

```go
package main

import (
	"context"
	"fmt"
	"github.com/go-auxiliaries/ctxticker"
	"time"
	"sync"
)

type Service struct {
	ticker *ctxticker.Ticker
	period time.Duration
	body func(ctx context.Context)
}

// Service that is running body every period
func NewService(period time.Duration, body func(ctx context.Context)) *Service {
    return &Service{
		period: period,
		body: body,
    }
}

func (s *Service) Start(ctx context.Context) {
	s.ticker = ctxticker.New(ctx, time.Duration(0), time.Second, true)
	go s.serviceBodyCycle(ctx)
}

func (s *Service) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
    }
}

func (s *Service) serviceBodyCycle(ctx context.Context) {
	defer s.ticker.Stop()
    for {
        if err := s.ticker.Tick(); err != nil {
            break
        }
        s.body(ctx)
    }
}
```


### 2. Service example, body execution time is limited ###

```go
package main

import (
	"context"
	"fmt"
	"github.com/go-auxiliaries/ctxticker"
	"time"
	"sync"
)

type Service struct {
	ticker *ctxticker.Ticker
	period time.Duration
	body func(ctx context.Context)
}

// Service that is running body every period, 
// but body execution time is limited by context timeout
func NewService(period time.Duration, body func(ctx context.Context)) *Service {
	return &Service{
		period: period,
		body: body,
	}
}

func (s *Service) Start(ctx context.Context) {
	s.ticker = ctxticker.New(ctx, time.Duration(0), time.Second, true)
	go s.serviceBodyCycle(ctx)
}

func (s *Service) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
}

func (s *Service) serviceBodyCycle(ctx context.Context) {
	defer s.ticker.Stop()
	for {
		if err := s.ticker.Tick(); err != nil {
			break
		}
		s.body(ctx)
	}
}

func (s *Service) serviceBody(ticker *ctxticker.Ticker) {
	ctx, cancel := ticker.GetTickContext()
    defer cancel()
	s.body(ctx)
}
```