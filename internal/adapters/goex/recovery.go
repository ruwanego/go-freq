package goex

import "gofreq/internal/recovery"

type RecoveryAdapter struct {
	client *Client
}

func NewRecoveryAdapter(c *Client) *RecoveryAdapter {
	return &RecoveryAdapter{client: c}
}

func (r *RecoveryAdapter) GetOpenOrders() ([]recovery.OpenOrder, error) {
	rawOrders, err := r.client.GetOpenOrders()
	if err != nil {
		return nil, err
	}

	out := make([]recovery.OpenOrder, 0, len(rawOrders))
	for _, o := range rawOrders {
		mapped, err := MapOrderToOpenOrder(o.ClientOrderID, o.ExchangeID, o.Pair)
		if err != nil {
			continue
		}

		out = append(out, mapped)
	}

	return out, nil
}

func (r *RecoveryAdapter) GetTradesSince(ts int64) ([]recovery.Trade, error) {
	rawTrades, err := r.client.GetTradesSince(ts)
	if err != nil {
		return nil, err
	}

	out := make([]recovery.Trade, 0, len(rawTrades))
	for _, t := range rawTrades {
		mapped, err := MapTradeToTrade(t.ClientOrderID, t.Amount, t.Price, t.Timestamp)
		if err != nil {
			continue
		}

		out = append(out, mapped)
	}

	return out, nil
}
