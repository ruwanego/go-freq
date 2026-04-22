package recovery_test

import (
	"github.com/shopspring/decimal"
	"testing"

	"gofreq/internal/persistence"
	"gofreq/internal/recovery"
)

type tradeExchange struct {
	trades []recovery.Trade
}

func (t *tradeExchange) GetOpenOrders() ([]recovery.OpenOrder, error) {
	return nil, nil
}

func (t *tradeExchange) GetTradesSince(int64) ([]recovery.Trade, error) {
	return t.trades, nil
}

func TestReconcile_FullFill(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:      "GF-MACD-1000-0001",
		ClientOrderID: "GF-MACD-1000-0001",
		Pair:          "BTC/USDT",
		Amount:        decimal.NewFromInt(10),
		State:         persistence.OrderStatePending,
		CreatedAt:     1000,
		UpdatedAt:     1000,
	}

	if err := store.CreateOrder(rec); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := store.UpdateOrderState("GF-MACD-1000-0001", persistence.OrderStateSubmitted, 1001, "ex-1"); err != nil {
		t.Fatalf("submit failed: %v", err)
	}

	ex := &tradeExchange{
		trades: []recovery.Trade{
			{ClientOrderID: "GF-MACD-1000-0001", Amount: decimal.NewFromInt(10)},
		},
	}

	r := recovery.NewReconciler(ex, store)

	if err := r.Run(0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	o, err := store.GetOrder("GF-MACD-1000-0001")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if o.State != persistence.OrderStateFilled {
		t.Fatalf("expected FILLED")
	}
}

func TestReconcile_PartialFill(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:      "GF-MACD-1000-0002",
		ClientOrderID: "GF-MACD-1000-0002",
		Pair:          "BTC/USDT",
		Amount:        decimal.NewFromInt(10),
		State:         persistence.OrderStatePending,
		CreatedAt:     1000,
		UpdatedAt:     1000,
	}

	if err := store.CreateOrder(rec); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := store.UpdateOrderState("GF-MACD-1000-0002", persistence.OrderStateSubmitted, 1001, "ex-2"); err != nil {
		t.Fatalf("submit failed: %v", err)
	}

	ex := &tradeExchange{
		trades: []recovery.Trade{
			{ClientOrderID: "GF-MACD-1000-0002", Amount: decimal.NewFromInt(4)},
		},
	}

	r := recovery.NewReconciler(ex, store)

	if err := r.Run(0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	o, err := store.GetOrder("GF-MACD-1000-0002")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if o.State != persistence.OrderStatePartiallyFilled {
		t.Fatalf("expected PARTIALLY_FILLED")
	}
}

func TestReconcile_Cancelled(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:      "GF-MACD-1000-0003",
		ClientOrderID: "GF-MACD-1000-0003",
		Pair:          "BTC/USDT",
		Amount:        decimal.NewFromInt(10),
		State:         persistence.OrderStatePending,
		CreatedAt:     1000,
		UpdatedAt:     1000,
	}

	if err := store.CreateOrder(rec); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := store.UpdateOrderState("GF-MACD-1000-0003", persistence.OrderStateSubmitted, 1001, "ex-3"); err != nil {
		t.Fatalf("submit failed: %v", err)
	}

	ex := &tradeExchange{
		trades: []recovery.Trade{},
	}

	r := recovery.NewReconciler(ex, store)

	if err := r.Run(0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	o, err := store.GetOrder("GF-MACD-1000-0003")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if o.State != persistence.OrderStateCancelled {
		t.Fatalf("expected CANCELLED")
	}
}
