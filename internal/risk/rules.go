package risk

import "github.com/shopspring/decimal"

func MaxOrderSizeRule(max decimal.Decimal) func(amount decimal.Decimal) Result {
	return func(amount decimal.Decimal) Result {
		if amount.GreaterThan(max) {
			return Result{
				Decision: DecisionReject,
				Reason:   "amount exceeds max order size",
			}
		}
		return Result{Decision: DecisionApprove}
	}
}

func NonZeroAmountRule() func(amount decimal.Decimal) Result {
	return func(amount decimal.Decimal) Result {
		if amount.LessThanOrEqual(decimal.Zero) {
			return Result{
				Decision: DecisionReject,
				Reason:   "invalid amount",
			}
		}
		return Result{Decision: DecisionApprove}
	}
}

func MaxNotionalRule(max decimal.Decimal) func(price, amount decimal.Decimal) Result {
	return func(price, amount decimal.Decimal) Result {
		if price.Mul(amount).GreaterThan(max) {
			return Result{
				Decision: DecisionReject,
				Reason:   "notional exceeds limit",
			}
		}
		return Result{Decision: DecisionApprove}
	}
}
