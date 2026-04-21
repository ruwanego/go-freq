package goex_test

import (
	"testing"

	adapter "gofreq/internal/adapters/goex"
)

func TestMapGoexCandle_Valid(t *testing.T) {
	got, err := adapter.MapGoexCandle("BTC/USDT", 1000, 10, 12, 9, 11, 100, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Pair != "BTC/USDT" || got.Timestamp != 1000 || got.Open != 10 || got.High != 12 || got.Low != 9 || got.Close != 11 || got.Volume != 100 || !got.Closed {
		t.Fatalf("mapped candle mismatch")
	}
}

func TestMapGoexCandle_InvalidTimestamp(t *testing.T) {
	_, err := adapter.MapGoexCandle("BTC/USDT", 0, 10, 12, 9, 11, 100, true)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapGoexCandle_HighLessThanLow(t *testing.T) {
	_, err := adapter.MapGoexCandle("BTC/USDT", 1000, 10, 8, 9, 11, 100, true)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapGoexCandle_InvalidOpenOrClose(t *testing.T) {
	_, err := adapter.MapGoexCandle("BTC/USDT", 1000, 0, 12, 9, 11, 100, true)
	if err == nil {
		t.Fatalf("expected error")
	}

	_, err = adapter.MapGoexCandle("BTC/USDT", 1000, 10, 12, 9, 0, 100, true)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapGoexCandle_NotClosed(t *testing.T) {
	_, err := adapter.MapGoexCandle("BTC/USDT", 1000, 10, 12, 9, 11, 100, false)
	if err == nil {
		t.Fatalf("expected error")
	}
}
