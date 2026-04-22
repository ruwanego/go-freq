package execution

import (
	"github.com/shopspring/decimal"
	"gofreq/pkg/actions"
	coreexec "gofreq/pkg/execution"
)

type RiskEngine interface {
	Apply(input []actions.Action) ([]actions.Action, []coreexec.RejectedAction)
}

type BasicRisk struct {
	MaxPerTrade decimal.Decimal
}

func (r *BasicRisk) Apply(input []actions.Action) ([]actions.Action, []coreexec.RejectedAction) {
	accepted := make([]actions.Action, 0, len(input))
	rejected := make([]coreexec.RejectedAction, 0)

	for _, a := range input {
		if r.MaxPerTrade.GreaterThan(decimal.Zero) && a.Amount.GreaterThan(r.MaxPerTrade) {
			rejected = append(rejected, coreexec.RejectedAction{
				Action: a,
				Reason: "risk_max_per_trade_exceeded",
			})
			continue
		}

		accepted = append(accepted, a)
	}

	return accepted, rejected
}
