package persistence_test

import (
	"os"
	"path/filepath"
	"testing"

	"gofreq/internal/persistence"
)

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

func TestCreateAndGetOrder(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:      "eng-1",
		ClientOrderID: "cli-1",
		StrategyName:  "macd",
		Pair:          "BTC/USDT",
		Side:          "BUY",
		Price:         60000,
		Amount:        1.5,
		Tag:           "entry",
		State:         persistence.OrderStatePending,
		CreatedAt:     1000,
		UpdatedAt:     1000,
	}

	if err := store.CreateOrder(rec); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := store.GetOrder("eng-1")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got.EngineID != rec.EngineID {
		t.Fatalf("engine id mismatch")
	}
	if got.State != persistence.OrderStatePending {
		t.Fatalf("state mismatch")
	}
}

func TestCreateOrder_RejectsDuplicateEngineID(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:      "eng-1",
		ClientOrderID: "cli-1",
		StrategyName:  "macd",
		Pair:          "BTC/USDT",
		Side:          "BUY",
		Price:         60000,
		Amount:        1.0,
		Tag:           "entry",
		State:         persistence.OrderStatePending,
		CreatedAt:     1000,
		UpdatedAt:     1000,
	}

	if err := store.CreateOrder(rec); err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	err := store.CreateOrder(rec)
	if err == nil {
		t.Fatalf("expected duplicate create to fail")
	}
}

func TestGetOrder_NotFound(t *testing.T) {
	store := newTestStore(t)

	_, err := store.GetOrder("missing")
	if err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestUpdateOrderState_ValidTransition(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:      "eng-2",
		ClientOrderID: "cli-2",
		StrategyName:  "macd",
		Pair:          "ETH/USDT",
		Side:          "SELL",
		Price:         3000,
		Amount:        2.0,
		Tag:           "exit",
		State:         persistence.OrderStatePending,
		CreatedAt:     1000,
		UpdatedAt:     1000,
	}

	if err := store.CreateOrder(rec); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if err := store.UpdateOrderState("eng-2", persistence.OrderStateSubmitted, 2000, "ex-22"); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	got, err := store.GetOrder("eng-2")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got.State != persistence.OrderStateSubmitted {
		t.Fatalf("expected submitted state")
	}
	if got.ExchangeID != "ex-22" {
		t.Fatalf("expected exchange id to be updated")
	}
	if got.UpdatedAt != 2000 {
		t.Fatalf("expected updated_at to change")
	}
}

func TestUpdateOrderState_InvalidTransitionRejected(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:      "eng-3",
		ClientOrderID: "cli-3",
		StrategyName:  "macd",
		Pair:          "SOL/USDT",
		Side:          "BUY",
		Price:         150,
		Amount:        3.0,
		Tag:           "entry",
		State:         persistence.OrderStatePending,
		CreatedAt:     1000,
		UpdatedAt:     1000,
	}

	if err := store.CreateOrder(rec); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err := store.UpdateOrderState("eng-3", persistence.OrderStateFilled, 2000, "")
	if err == nil {
		t.Fatalf("expected invalid transition to fail")
	}

	got, err := store.GetOrder("eng-3")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.State != persistence.OrderStatePending {
		t.Fatalf("state must remain unchanged after failed transition")
	}
}
