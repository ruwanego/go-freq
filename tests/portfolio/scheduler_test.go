package portfolio_test

import (
	"github.com/shopspring/decimal"
	"testing"

	"gofreq/internal/marketdata"
	"gofreq/internal/portfolio"
)

func TestSchedulerAlignedArrivalEmitsTick(t *testing.T) {
	s := portfolio.NewScheduler([]string{"BTC/USDT", "ETH/USDT"}, 1000)

	_, ready, err := s.OnCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 100, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true}, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ready {
		t.Fatalf("expected not ready")
	}

	tick, ready, err := s.OnCandle(marketdata.Candle{Pair: "ETH/USDT", Timestamp: 100, Open: decimal.NewFromInt(2), High: decimal.NewFromInt(2), Low: decimal.NewFromInt(2), Close: decimal.NewFromInt(2), Closed: true}, 1001)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ready {
		t.Fatalf("expected ready")
	}
	if tick.Timestamp != 100 {
		t.Fatalf("unexpected timestamp: %d", tick.Timestamp)
	}
}

func TestSchedulerDelayedArrivalBeforeTimeoutEmitsTick(t *testing.T) {
	s := portfolio.NewScheduler([]string{"BTC/USDT", "ETH/USDT"}, 1000)

	_, ready, err := s.OnCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 100, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true}, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ready {
		t.Fatalf("expected not ready")
	}

	_, ready, err = s.OnCandle(marketdata.Candle{Pair: "ETH/USDT", Timestamp: 100, Open: decimal.NewFromInt(2), High: decimal.NewFromInt(2), Low: decimal.NewFromInt(2), Close: decimal.NewFromInt(2), Closed: true}, 1999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ready {
		t.Fatalf("expected ready before timeout")
	}
}

func TestSchedulerTimeoutReturnsError(t *testing.T) {
	s := portfolio.NewScheduler([]string{"BTC/USDT", "ETH/USDT"}, 1000)

	_, ready, err := s.OnCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 100, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true}, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ready {
		t.Fatalf("expected not ready")
	}

	_, ready, err = s.OnCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 100, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true}, 2001)
	if err == nil {
		t.Fatalf("expected timeout error")
	}
	if ready {
		t.Fatalf("expected not ready on timeout")
	}
}

func TestSchedulerBufferClearedAfterEmit(t *testing.T) {
	s := portfolio.NewScheduler([]string{"BTC/USDT", "ETH/USDT"}, 1000)

	_, _, _ = s.OnCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 100, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true}, 1000)
	_, ready, err := s.OnCandle(marketdata.Candle{Pair: "ETH/USDT", Timestamp: 100, Open: decimal.NewFromInt(2), High: decimal.NewFromInt(2), Low: decimal.NewFromInt(2), Close: decimal.NewFromInt(2), Closed: true}, 1001)
	if err != nil || !ready {
		t.Fatalf("expected first tick ready")
	}

	_, ready, err = s.OnCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 100, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true}, 1002)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ready {
		t.Fatalf("expected no duplicate emission")
	}
}

func TestSchedulerNoDuplicateEmissions(t *testing.T) {
	s := portfolio.NewScheduler([]string{"BTC/USDT", "ETH/USDT"}, 1000)

	_, _, _ = s.OnCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 100, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true}, 1000)
	_, ready, err := s.OnCandle(marketdata.Candle{Pair: "ETH/USDT", Timestamp: 100, Open: decimal.NewFromInt(2), High: decimal.NewFromInt(2), Low: decimal.NewFromInt(2), Close: decimal.NewFromInt(2), Closed: true}, 1001)
	if err != nil || !ready {
		t.Fatalf("expected emission")
	}

	_, ready, err = s.OnCandle(marketdata.Candle{Pair: "ETH/USDT", Timestamp: 100, Open: decimal.NewFromInt(2), High: decimal.NewFromInt(2), Low: decimal.NewFromInt(2), Close: decimal.NewFromInt(2), Closed: true}, 1002)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ready {
		t.Fatalf("expected no duplicate emission")
	}
}
