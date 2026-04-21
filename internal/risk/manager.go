package risk

type Manager struct {
	maxOrderSizeRule func(amount float64) Result
	maxNotionalRule  func(price, amount float64) Result
	nonZeroRule      func(amount float64) Result
}

func NewManager(maxOrderSize float64, maxNotional float64) *Manager {
	return &Manager{
		maxOrderSizeRule: MaxOrderSizeRule(maxOrderSize),
		maxNotionalRule:  MaxNotionalRule(maxNotional),
		nonZeroRule:      NonZeroAmountRule(),
	}
}

func (m *Manager) Evaluate(price float64, amount float64) Result {
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
