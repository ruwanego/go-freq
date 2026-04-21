package risk_test

import (
	"testing"

	"gofreq/internal/risk"
)

func TestNonZeroAmountRule_ApproveValid(t *testing.T) {
	rule := risk.NonZeroAmountRule()
	res := rule(1)
	if res.Decision != risk.DecisionApprove {
		t.Fatalf("expected approve")
	}
}

func TestNonZeroAmountRule_RejectZero(t *testing.T) {
	rule := risk.NonZeroAmountRule()
	res := rule(0)
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason == "" {
		t.Fatalf("expected reject reason")
	}
}

func TestMaxOrderSizeRule_ApproveWithinLimit(t *testing.T) {
	rule := risk.MaxOrderSizeRule(5)
	res := rule(3)
	if res.Decision != risk.DecisionApprove {
		t.Fatalf("expected approve")
	}
}

func TestMaxOrderSizeRule_RejectAboveLimit(t *testing.T) {
	rule := risk.MaxOrderSizeRule(5)
	res := rule(6)
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason != "amount exceeds max order size" {
		t.Fatalf("unexpected reason: %s", res.Reason)
	}
}

func TestMaxNotionalRule_ApproveWithinLimit(t *testing.T) {
	rule := risk.MaxNotionalRule(1000)
	res := rule(100, 5)
	if res.Decision != risk.DecisionApprove {
		t.Fatalf("expected approve")
	}
}

func TestMaxNotionalRule_RejectAboveLimit(t *testing.T) {
	rule := risk.MaxNotionalRule(1000)
	res := rule(300, 4)
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason != "notional exceeds limit" {
		t.Fatalf("unexpected reason: %s", res.Reason)
	}
}
