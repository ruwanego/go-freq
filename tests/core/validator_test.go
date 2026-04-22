package core

import (
	"github.com/shopspring/decimal"
	"gofreq/pkg/actions"
	"gofreq/pkg/execution"
	"testing"
)

func TestRejectDuplicateTags(t *testing.T) {
	acts := []actions.Action{
		{Pair: "BTC/USDT", Amount: decimal.NewFromInt(1), Tag: "x"},
		{Pair: "BTC/USDT", Amount: decimal.NewFromInt(2), Tag: "x"},
	}

	res, _ := execution.ValidateActions(acts)

	if len(res.Rejected) != 1 {
		t.Fatalf("expected 1 rejection")
	}
}

func TestRejectInvalidAmount(t *testing.T) {
	acts := []actions.Action{
		{Pair: "BTC/USDT", Amount: decimal.NewFromInt(0), Tag: "x"},
	}

	res, err := execution.ValidateActions(acts)

	if err == nil {
		t.Fatalf("expected error")
	}

	if len(res.Rejected) != 1 {
		t.Fatalf("expected rejection")
	}
}
