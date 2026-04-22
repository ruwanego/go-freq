package recovery

import (
	"github.com/shopspring/decimal"
	"gofreq/internal/persistence"
)

type Reconciler struct {
	exchange Exchange
	store    *persistence.Store
}

func NewReconciler(ex Exchange, store *persistence.Store) *Reconciler {
	return &Reconciler{
		exchange: ex,
		store:    store,
	}
}

func (r *Reconciler) Run(since int64) error {
	trades, err := r.exchange.GetTradesSince(since)
	if err != nil {
		return err
	}

	tradeMap := map[string]decimal.Decimal{}
	for _, t := range trades {
		tradeMap[t.ClientOrderID] = tradeMap[t.ClientOrderID].Add(t.Amount)
	}

	orders, err := r.store.ListOrders()
	if err != nil {
		return err
	}

	for _, o := range orders {
		if o.State != persistence.OrderStateSubmitted && o.State != persistence.OrderStatePending {
			continue
		}

		filled := tradeMap[o.ClientOrderID]
		if filled.GreaterThan(decimal.Zero) {
			if filled.LessThan(o.Amount) {
				if err := r.store.UpdateOrderState(
					o.EngineID,
					persistence.OrderStatePartiallyFilled,
					o.UpdatedAt+1,
					o.ExchangeID,
				); err != nil {
					return err
				}
			} else {
				if err := r.store.UpdateOrderState(
					o.EngineID,
					persistence.OrderStateFilled,
					o.UpdatedAt+1,
					o.ExchangeID,
				); err != nil {
					return err
				}
			}
		} else {
			if err := r.store.UpdateOrderState(
				o.EngineID,
				persistence.OrderStateCancelled,
				o.UpdatedAt+1,
				o.ExchangeID,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
