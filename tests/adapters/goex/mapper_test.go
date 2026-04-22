package goex_test

import (
	"github.com/shopspring/decimal"
	"testing"

	adapter "gofreq/internal/adapters/goex"
)

func TestMapOrderToOpenOrder_Valid(t *testing.T) {
	got, err := adapter.MapOrderToOpenOrder("cid-1", "ex-1", "BTC/USDT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ClientOrderID != "cid-1" {
		t.Fatalf("client order id mismatch")
	}
	if got.ExchangeID != "ex-1" {
		t.Fatalf("exchange id mismatch")
	}
	if got.Pair != "BTC/USDT" {
		t.Fatalf("pair mismatch")
	}
}

func TestMapOrderToOpenOrder_MissingClientOrderID(t *testing.T) {
	_, err := adapter.MapOrderToOpenOrder("", "ex-1", "BTC/USDT")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapTradeToTrade_Valid(t *testing.T) {
	got, err := adapter.MapTradeToTrade("cid-2", decimal.RequireFromString("1.25"), decimal.NewFromInt(65000), 1234)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ClientOrderID != "cid-2" {
		t.Fatalf("client order id mismatch")
	}
	if !got.Amount.Equal(decimal.RequireFromString("1.25")) {
		t.Fatalf("amount mismatch")
	}
	if !got.Price.Equal(decimal.NewFromInt(65000)) {
		t.Fatalf("price mismatch")
	}
	if got.Timestamp != 1234 {
		t.Fatalf("timestamp mismatch")
	}
}

func TestMapTradeToTrade_MissingClientOrderID(t *testing.T) {
	_, err := adapter.MapTradeToTrade("", decimal.RequireFromString("1"), decimal.NewFromInt(100), 1)
	if err == nil {
		t.Fatalf("expected error")
	}
}
