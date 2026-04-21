package risk

func MaxOrderSizeRule(max float64) func(amount float64) Result {
	return func(amount float64) Result {
		if amount > max {
			return Result{
				Decision: DecisionReject,
				Reason:   "amount exceeds max order size",
			}
		}
		return Result{Decision: DecisionApprove}
	}
}

func NonZeroAmountRule() func(amount float64) Result {
	return func(amount float64) Result {
		if amount <= 0 {
			return Result{
				Decision: DecisionReject,
				Reason:   "invalid amount",
			}
		}
		return Result{Decision: DecisionApprove}
	}
}

func MaxNotionalRule(max float64) func(price, amount float64) Result {
	return func(price, amount float64) Result {
		if price*amount > max {
			return Result{
				Decision: DecisionReject,
				Reason:   "notional exceeds limit",
			}
		}
		return Result{Decision: DecisionApprove}
	}
}
