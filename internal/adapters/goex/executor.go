package goex

type Executor struct {
	client *Client
}

func NewExecutor(c *Client) *Executor {
	return &Executor{client: c}
}

func (e *Executor) SubmitOrder(intent OrderIntent) (OrderAck, error) {
	req, err := MapIntentToGoex(intent)
	if err != nil {
		return OrderAck{}, err
	}

	resp, err := e.client.SubmitOrder(req)
	if err != nil {
		return OrderAck{}, err
	}

	if resp.ExchangeID == "" {
		return OrderAck{}, ErrMissingExchangeID
	}

	return OrderAck{
		ClientOrderID: intent.ClientOrderID,
		ExchangeID:    resp.ExchangeID,
	}, nil
}

func (e *Executor) CancelOrder(clientOrderID string) error {
	if clientOrderID == "" {
		return ErrEmptyClientOrderID
	}

	orders, err := e.client.GetOpenOrders()
	if err != nil {
		return err
	}

	var exchangeID string
	for _, o := range orders {
		if o.ClientOrderID == clientOrderID {
			exchangeID = o.ExchangeID
			break
		}
	}

	if exchangeID == "" {
		return ErrOrderNotFound
	}

	return e.client.CancelOrder(exchangeID)
}
