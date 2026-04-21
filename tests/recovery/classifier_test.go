package recovery_test

import (
	"os"
	"path/filepath"
	"testing"

	"gofreq/internal/persistence"
	"gofreq/internal/recovery"
)

type mockExchange struct {
	orders []recovery.OpenOrder
}

func (m *mockExchange) GetOpenOrders() ([]recovery.OpenOrder, error) {
	return m.orders, nil
}

func (m *mockExchange) GetTradesSince(int64) ([]recovery.Trade, error) {
	return nil, nil
}

func newTestStore(t *testing.T) *persistence.Store {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	store, err := persistence.Open(path)
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}

	t.Cleanup(func() {
		_ = store.Close()
		_ = os.Remove(path)
	})

	return store
}

func TestClassification_KnownAndAlien(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:      "GF-MACD-1000-0001",
		ClientOrderID: "GF-MACD-1000-0001",
		Pair:          "BTC/USDT",
		Amount:        1,
		State:         persistence.OrderStatePending,
		CreatedAt:     1000,
		UpdatedAt:     1000,
	}

	if err := store.CreateOrder(rec); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	ex := &mockExchange{
		orders: []recovery.OpenOrder{
			{
				ClientOrderID: "GF-MACD-1000-0001",
				ExchangeID:    "ex-1",
			},
			{
				ClientOrderID: "UNKNOWN-ORDER",
				ExchangeID:    "ex-2",
			},
		},
	}

	classifier := recovery.NewClassifier(ex, store)

	out, err := classifier.Run()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if len(out) != 2 {
		t.Fatalf("expected 2 results")
	}

	if out[0].Class != recovery.Known {
		t.Fatalf("expected known")
	}
	if out[1].Class != recovery.Alien {
		t.Fatalf("expected alien")
	}
}

func TestPromotion_PendingToSubmitted(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:      "GF-MACD-1000-0001",
		ClientOrderID: "GF-MACD-1000-0001",
		Pair:          "BTC/USDT",
		Amount:        1,
		State:         persistence.OrderStatePending,
		CreatedAt:     1000,
		UpdatedAt:     1000,
	}

	if err := store.CreateOrder(rec); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	ex := &mockExchange{
		orders: []recovery.OpenOrder{
			{
				ClientOrderID: "GF-MACD-1000-0001",
				ExchangeID:    "ex-1",
			},
		},
	}

	classifier := recovery.NewClassifier(ex, store)

	if _, err := classifier.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := store.GetOrder("GF-MACD-1000-0001")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if updated.State != persistence.OrderStateSubmitted {
		t.Fatalf("expected SUBMITTED")
	}

	if updated.ExchangeID != "ex-1" {
		t.Fatalf("exchange id must be set")
	}
}

func TestNoRePromotionIfAlreadySubmitted(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:      "GF-MACD-1000-0001",
		ClientOrderID: "GF-MACD-1000-0001",
		Pair:          "BTC/USDT",
		Amount:        1,
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

	ex := &mockExchange{
		orders: []recovery.OpenOrder{
			{
				ClientOrderID: "GF-MACD-1000-0001",
				ExchangeID:    "ex-1",
			},
		},
	}

	classifier := recovery.NewClassifier(ex, store)

	if _, err := classifier.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := store.GetOrder("GF-MACD-1000-0001")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if updated.State != persistence.OrderStateSubmitted {
		t.Fatalf("must remain submitted")
	}
}
