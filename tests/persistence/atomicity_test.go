package persistence_test

import (
	"testing"

	"gofreq/internal/persistence"
)

func TestCreateOrder_InvalidRecordDoesNotPersistPartialState(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:  "",
		Pair:      "BTC/USDT",
		Amount:    1.0,
		State:     persistence.OrderStatePending,
		CreatedAt: 1000,
		UpdatedAt: 1000,
	}

	err := store.CreateOrder(rec)
	if err == nil {
		t.Fatalf("expected invalid record to fail")
	}

	_, err = store.GetOrder("")
	if err == nil {
		t.Fatalf("invalid record must not be persisted")
	}
}

func TestFailedStateTransitionIsAtomic(t *testing.T) {
	store := newTestStore(t)

	rec := persistence.OrderRecord{
		EngineID:      "eng-4",
		ClientOrderID: "cli-4",
		StrategyName:  "macd",
		Pair:          "BTC/USDT",
		Side:          "BUY",
		Price:         61000,
		Amount:        1.0,
		Tag:           "entry",
		State:         persistence.OrderStatePending,
		CreatedAt:     1000,
		UpdatedAt:     1000,
	}

	if err := store.CreateOrder(rec); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err := store.UpdateOrderState("eng-4", persistence.OrderStateFilled, 2000, "bad-jump")
	if err == nil {
		t.Fatalf("expected invalid transition to fail")
	}

	got, err := store.GetOrder("eng-4")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got.State != persistence.OrderStatePending {
		t.Fatalf("record must remain unchanged after failed transaction")
	}
	if got.ExchangeID != "" {
		t.Fatalf("exchange id must not be partially written")
	}
	if got.UpdatedAt != 1000 {
		t.Fatalf("updated_at must not be partially written")
	}
}
