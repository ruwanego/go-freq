package recovery

import "gofreq/internal/persistence"

type Classification string

const (
	Known Classification = "KNOWN"
	Alien Classification = "ALIEN"
)

type ClassifiedOrder struct {
	ClientOrderID string
	ExchangeID    string
	Class         Classification
}

type Classifier struct {
	exchange Exchange
	store    *persistence.Store
}

func NewClassifier(ex Exchange, store *persistence.Store) *Classifier {
	return &Classifier{
		exchange: ex,
		store:    store,
	}
}

func (c *Classifier) Run() ([]ClassifiedOrder, error) {
	openOrders, err := c.exchange.GetOpenOrders()
	if err != nil {
		return nil, err
	}

	result := make([]ClassifiedOrder, 0, len(openOrders))

	for _, o := range openOrders {
		rec, err := c.store.GetOrder(o.ClientOrderID)
		if err == nil {
			result = append(result, ClassifiedOrder{
				ClientOrderID: o.ClientOrderID,
				ExchangeID:    o.ExchangeID,
				Class:         Known,
			})

			if rec.State == persistence.OrderStatePending {
				if err := c.store.UpdateOrderState(
					o.ClientOrderID,
					persistence.OrderStateSubmitted,
					0,
					o.ExchangeID,
				); err != nil {
					return nil, err
				}
			}

			continue
		}

		if err != persistence.ErrOrderNotFound {
			return nil, err
		}

		result = append(result, ClassifiedOrder{
			ClientOrderID: o.ClientOrderID,
			ExchangeID:    o.ExchangeID,
			Class:         Alien,
		})
	}

	return result, nil
}
