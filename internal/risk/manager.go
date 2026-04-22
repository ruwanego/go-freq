package risk

import "github.com/shopspring/decimal"

type Manager struct {
	maxOrderSizeRule func(amount decimal.Decimal) Result
	maxNotionalRule  func(price, amount decimal.Decimal) Result
	nonZeroRule      func(amount decimal.Decimal) Result
}

func NewManager(maxOrderSize decimal.Decimal, maxNotional decimal.Decimal) *Manager {
	return &Manager{
		maxOrderSizeRule: MaxOrderSizeRule(maxOrderSize),
		maxNotionalRule:  MaxNotionalRule(maxNotional),
		nonZeroRule:      NonZeroAmountRule(),
	}
}

func (m *Manager) Evaluate(price decimal.Decimal, amount decimal.Decimal) Result {
	if r := m.nonZeroRule(amount); r.Decision == DecisionReject {
		return r
	}

	if r := m.maxOrderSizeRule(amount); r.Decision == DecisionReject {
		return r
	}

	if r := m.maxNotionalRule(price, amount); r.Decision == DecisionReject {
		return r
	}

	return Result{Decision: DecisionApprove}
}
