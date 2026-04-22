package portfolio_test

import (
	"github.com/shopspring/decimal"
	"testing"

	"gofreq/internal/marketdata"
	"gofreq/internal/portfolio"
)

func TestAggregatorPartialFillNotReady(t *testing.T) {
	agg := portfolio.NewAggregator([]string{"BTC/USDT", "ETH/USDT"})

	tick, ready := agg.Add(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 100, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true})
	if ready {
		t.Fatalf("expected not ready")
	}
	if tick.Timestamp != 0 {
		t.Fatalf("expected zero tick when not ready")
	}
}

func TestAggregatorFullFillReady(t *testing.T) {
	agg := portfolio.NewAggregator([]string{"BTC/USDT", "ETH/USDT"})

	_, ready := agg.Add(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 100, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true})
	if ready {
		t.Fatalf("expected first add not ready")
	}

	tick, ready := agg.Add(marketdata.Candle{Pair: "ETH/USDT", Timestamp: 100, Open: decimal.NewFromInt(2), High: decimal.NewFromInt(2), Low: decimal.NewFromInt(2), Close: decimal.NewFromInt(2), Closed: true})
	if !ready {
		t.Fatalf("expected ready")
	}
	if tick.Timestamp != 100 {
		t.Fatalf("unexpected timestamp: %d", tick.Timestamp)
	}
	if len(tick.Candles) != 2 {
		t.Fatalf("expected 2 candles")
	}
	if _, ok := tick.Candles["BTC/USDT"]; !ok {
		t.Fatalf("missing BTC candle")
	}
	if _, ok := tick.Candles["ETH/USDT"]; !ok {
		t.Fatalf("missing ETH candle")
	}
}

func TestAggregatorMultipleTimestampsHandledCorrectly(t *testing.T) {
	agg := portfolio.NewAggregator([]string{"BTC/USDT", "ETH/USDT"})

	_, ready := agg.Add(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 100, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true})
	if ready {
		t.Fatalf("expected ts=100 not ready")
	}

	_, ready = agg.Add(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 200, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true})
	if ready {
		t.Fatalf("expected ts=200 not ready")
	}

	tick, ready := agg.Add(marketdata.Candle{Pair: "ETH/USDT", Timestamp: 100, Open: decimal.NewFromInt(2), High: decimal.NewFromInt(2), Low: decimal.NewFromInt(2), Close: decimal.NewFromInt(2), Closed: true})
	if !ready {
		t.Fatalf("expected ts=100 ready")
	}
	if tick.Timestamp != 100 {
		t.Fatalf("unexpected timestamp: %d", tick.Timestamp)
	}

	tick, ready = agg.Add(marketdata.Candle{Pair: "ETH/USDT", Timestamp: 200, Open: decimal.NewFromInt(2), High: decimal.NewFromInt(2), Low: decimal.NewFromInt(2), Close: decimal.NewFromInt(2), Closed: true})
	if !ready {
		t.Fatalf("expected ts=200 ready")
	}
	if tick.Timestamp != 200 {
		t.Fatalf("unexpected timestamp: %d", tick.Timestamp)
	}
}
