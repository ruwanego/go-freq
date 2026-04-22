package risk_test

import (
	"github.com/shopspring/decimal"
	"testing"

	"gofreq/internal/risk"
)

func TestNonZeroAmountRule_ApproveValid(t *testing.T) {
	rule := risk.NonZeroAmountRule()
	res := rule(decimal.NewFromInt(1))
	if res.Decision != risk.DecisionApprove {
		t.Fatalf("expected approve")
	}
}

func TestNonZeroAmountRule_RejectZero(t *testing.T) {
	rule := risk.NonZeroAmountRule()
	res := rule(decimal.NewFromInt(0))
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason == "" {
		t.Fatalf("expected reject reason")
	}
}

func TestMaxOrderSizeRule_ApproveWithinLimit(t *testing.T) {
	rule := risk.MaxOrderSizeRule(decimal.NewFromInt(5))
	res := rule(decimal.NewFromInt(3))
	if res.Decision != risk.DecisionApprove {
		t.Fatalf("expected approve")
	}
}

func TestMaxOrderSizeRule_RejectAboveLimit(t *testing.T) {
	rule := risk.MaxOrderSizeRule(decimal.NewFromInt(5))
	res := rule(decimal.NewFromInt(6))
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason != "amount exceeds max order size" {
		t.Fatalf("unexpected reason: %s", res.Reason)
	}
}

func TestMaxNotionalRule_ApproveWithinLimit(t *testing.T) {
	rule := risk.MaxNotionalRule(decimal.NewFromInt(1000))
	res := rule(decimal.NewFromInt(100), decimal.NewFromInt(5))
	if res.Decision != risk.DecisionApprove {
		t.Fatalf("expected approve")
	}
}

func TestMaxNotionalRule_RejectAboveLimit(t *testing.T) {
	rule := risk.MaxNotionalRule(decimal.NewFromInt(1000))
	res := rule(decimal.NewFromInt(300), decimal.NewFromInt(4))
	if res.Decision != risk.DecisionReject {
		t.Fatalf("expected reject")
	}
	if res.Reason != "notional exceeds limit" {
		t.Fatalf("unexpected reason: %s", res.Reason)
	}
}
