package risk_test

import (
	"github.com/shopspring/decimal"
	"testing"

	"gofreq/internal/risk"
)

func TestManagerApproveValidOrder(t *testing.T) {
	m := risk.NewManager(decimal.NewFromInt(10), decimal.NewFromInt(10000))
	res := m.Evaluate(decimal.NewFromInt(100), decimal.NewFromInt(2))
	if res.Decision != risk.DecisionApprove {
		t.Fatalf("expected approve")
	}
}

func TestManagerRejectNonZeroRuleFirst(t *testing.T) {
	m := risk.NewManager(decimal.NewFromInt(10), decimal.NewFromInt(10000))
	res := m.Evaluate(decimal.NewFromInt(100), decimal.NewFromInt(0))
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason != "invalid amount" {
		t.Fatalf("unexpected reason: %s", res.Reason)
	}
}

func TestManagerRejectMaxOrderSize(t *testing.T) {
	m := risk.NewManager(decimal.NewFromInt(5), decimal.NewFromInt(10000))
	res := m.Evaluate(decimal.NewFromInt(100), decimal.NewFromInt(6))
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason != "amount exceeds max order size" {
		t.Fatalf("unexpected reason: %s", res.Reason)
	}
}

func TestManagerRejectMaxNotional(t *testing.T) {
	m := risk.NewManager(decimal.NewFromInt(10), decimal.NewFromInt(500))
	res := m.Evaluate(decimal.NewFromInt(200), decimal.NewFromInt(3))
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason != "notional exceeds limit" {
		t.Fatalf("unexpected reason: %s", res.Reason)
	}
}

func TestManagerRuleOrderDeterministic(t *testing.T) {
	m := risk.NewManager(decimal.NewFromInt(1), decimal.NewFromInt(1))
	res := m.Evaluate(decimal.NewFromInt(100), decimal.NewFromInt(0))
	if res.Reason != "invalid amount" {
		t.Fatalf("expected first rule rejection, got %s", res.Reason)
	}
}
