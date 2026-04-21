package persistence_test

import (
	"testing"

	"gofreq/internal/persistence"
)

func TestValidateTransition_ValidPaths(t *testing.T) {
	valid := []struct {
		from persistence.OrderState
		to   persistence.OrderState
	}{
		{persistence.OrderStatePending, persistence.OrderStateSubmitted},
		{persistence.OrderStatePending, persistence.OrderStateCancelled},
		{persistence.OrderStateSubmitted, persistence.OrderStateFilled},
		{persistence.OrderStateSubmitted, persistence.OrderStateCancelled},
	}

	for _, tt := range valid {
		if err := persistence.ValidateTransition(tt.from, tt.to); err != nil {
			t.Fatalf("expected valid transition %s -> %s, got error: %v", tt.from, tt.to, err)
		}
	}
}

func TestValidateTransition_InvalidPaths(t *testing.T) {
	invalid := []struct {
		from persistence.OrderState
		to   persistence.OrderState
	}{
		{persistence.OrderStatePending, persistence.OrderStateFilled},
		{persistence.OrderStateSubmitted, persistence.OrderStatePending},
		{persistence.OrderStateFilled, persistence.OrderStateSubmitted},
		{persistence.OrderStateCancelled, persistence.OrderStatePending},
	}

	for _, tt := range invalid {
		if err := persistence.ValidateTransition(tt.from, tt.to); err == nil {
			t.Fatalf("expected invalid transition %s -> %s", tt.from, tt.to)
		}
	}
}
