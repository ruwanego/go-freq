package goex_test

import (
	"github.com/shopspring/decimal"
	"testing"

	adapter "gofreq/internal/adapters/goex"
)

func TestMapGoexCandle_Valid(t *testing.T) {
	got, err := adapter.MapGoexCandle("BTC/USDT", 1000, decimal.NewFromInt(10), decimal.NewFromInt(12), decimal.NewFromInt(9), decimal.NewFromInt(11), decimal.NewFromInt(100), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Pair != "BTC/USDT" || got.Timestamp != 1000 || !got.Open.Equal(decimal.NewFromInt(10)) || !got.High.Equal(decimal.NewFromInt(12)) || !got.Low.Equal(decimal.NewFromInt(9)) || !got.Close.Equal(decimal.NewFromInt(11)) || !got.Volume.Equal(decimal.NewFromInt(100)) || !got.Closed {
		t.Fatalf("mapped candle mismatch")
	}
}

func TestMapGoexCandle_InvalidTimestamp(t *testing.T) {
	_, err := adapter.MapGoexCandle("BTC/USDT", 0, decimal.NewFromInt(10), decimal.NewFromInt(12), decimal.NewFromInt(9), decimal.NewFromInt(11), decimal.NewFromInt(100), true)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapGoexCandle_HighLessThanLow(t *testing.T) {
	_, err := adapter.MapGoexCandle("BTC/USDT", 1000, decimal.NewFromInt(10), decimal.NewFromInt(8), decimal.NewFromInt(9), decimal.NewFromInt(11), decimal.NewFromInt(100), true)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapGoexCandle_InvalidOpenOrClose(t *testing.T) {
	_, err := adapter.MapGoexCandle("BTC/USDT", 1000, decimal.NewFromInt(0), decimal.NewFromInt(12), decimal.NewFromInt(9), decimal.NewFromInt(11), decimal.NewFromInt(100), true)
	if err == nil {
		t.Fatalf("expected error")
	}

	_, err = adapter.MapGoexCandle("BTC/USDT", 1000, decimal.NewFromInt(10), decimal.NewFromInt(12), decimal.NewFromInt(9), decimal.NewFromInt(0), decimal.NewFromInt(100), true)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapGoexCandle_NotClosed(t *testing.T) {
	_, err := adapter.MapGoexCandle("BTC/USDT", 1000, decimal.NewFromInt(10), decimal.NewFromInt(12), decimal.NewFromInt(9), decimal.NewFromInt(11), decimal.NewFromInt(100), false)
	if err == nil {
		t.Fatalf("expected error")
	}
}
