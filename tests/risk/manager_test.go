package risk_test

import (
	"testing"

	"gofreq/internal/risk"
)

func TestManagerApproveValidOrder(t *testing.T) {
	m := risk.NewManager(10, 10000)
	res := m.Evaluate(100, 2)
	if res.Decision != risk.DecisionApprove {
		t.Fatalf("expected approve")
	}
}

func TestManagerRejectNonZeroRuleFirst(t *testing.T) {
	m := risk.NewManager(10, 10000)
	res := m.Evaluate(100, 0)
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason != "invalid amount" {
		t.Fatalf("unexpected reason: %s", res.Reason)
	}
}

func TestManagerRejectMaxOrderSize(t *testing.T) {
	m := risk.NewManager(5, 10000)
	res := m.Evaluate(100, 6)
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason != "amount exceeds max order size" {
		t.Fatalf("unexpected reason: %s", res.Reason)
	}
}

func TestManagerRejectMaxNotional(t *testing.T) {
	m := risk.NewManager(10, 500)
	res := m.Evaluate(200, 3)
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason != "notional exceeds limit" {
		t.Fatalf("unexpected reason: %s", res.Reason)
	}
}

func TestManagerRuleOrderDeterministic(t *testing.T) {
	m := risk.NewManager(1, 1)
	res := m.Evaluate(100, 0)
	if res.Reason != "invalid amount" {
		t.Fatalf("expected first rule rejection, got %s", res.Reason)
	}
}
