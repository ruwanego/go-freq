package goex_test

import (
	"github.com/shopspring/decimal"
	"testing"

	adapter "gofreq/internal/adapters/goex"
)

func TestMapIntentToGoex_EmptyClientOrderID(t *testing.T) {
	_, err := adapter.MapIntentToGoex(adapter.OrderIntent{ClientOrderID: "", Amount: decimal.NewFromInt(1), Side: "BUY", Type: "MARKET"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapIntentToGoex_InvalidSide(t *testing.T) {
	_, err := adapter.MapIntentToGoex(adapter.OrderIntent{ClientOrderID: "cid", Amount: decimal.NewFromInt(1), Side: "HOLD", Type: "MARKET"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapIntentToGoex_InvalidType(t *testing.T) {
	_, err := adapter.MapIntentToGoex(adapter.OrderIntent{ClientOrderID: "cid", Amount: decimal.NewFromInt(1), Side: "BUY", Type: "STOP"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapIntentToGoex_ValidLimit(t *testing.T) {
	got, err := adapter.MapIntentToGoex(adapter.OrderIntent{
		ClientOrderID: "cid-1",
		Pair:          "BTC/USDT",
		Side:          "BUY",
		Type:          "LIMIT",
		Price:         decimal.NewFromInt(60000),
		Amount:        decimal.NewFromInt(1),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ClientOrderID != "cid-1" || got.Side != "buy" || got.Type != "limit" {
		t.Fatalf("mapping mismatch")
	}
}

func TestMapIntentToGoex_ValidMarket(t *testing.T) {
	got, err := adapter.MapIntentToGoex(adapter.OrderIntent{
		ClientOrderID: "cid-2",
		Pair:          "BTC/USDT",
		Side:          "SELL",
		Type:          "MARKET",
		Amount:        decimal.NewFromInt(2),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ClientOrderID != "cid-2" || got.Side != "sell" || got.Type != "market" {
		t.Fatalf("mapping mismatch")
	}
}

func TestMapIntentToGoex_LimitWithoutPrice(t *testing.T) {
	_, err := adapter.MapIntentToGoex(adapter.OrderIntent{
		ClientOrderID: "cid-3",
		Pair:          "BTC/USDT",
		Side:          "BUY",
		Type:          "LIMIT",
		Price:         decimal.NewFromInt(0),
		Amount:        decimal.NewFromInt(1),
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMapIntentToGoex_InvalidAmount(t *testing.T) {
	_, err := adapter.MapIntentToGoex(adapter.OrderIntent{
		ClientOrderID: "cid-4",
		Pair:          "BTC/USDT",
		Side:          "BUY",
		Type:          "MARKET",
		Amount:        decimal.NewFromInt(0),
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}
