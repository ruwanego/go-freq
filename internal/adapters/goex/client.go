package goex

import "errors"

var ErrClientNotConfigured = errors.New("client_not_configured")

type RawOrder struct {
	ClientOrderID string
	ExchangeID    string
	Pair          string
}

type RawTrade struct {
	ClientOrderID string
	Amount        float64
	Price         float64
	Timestamp     int64
}

type Transport interface {
	GetOpenOrders() ([]RawOrder, error)
	GetTradesSince(ts int64) ([]RawTrade, error)
	SubmitOrder(req GoexOrderRequest) (GoexOrderResponse, error)
	CancelOrder(exchangeID string) error
}

type Client struct {
	transport Transport
}

func NewClient(transport Transport) *Client {
	return &Client{transport: transport}
}

func (c *Client) GetOpenOrders() ([]RawOrder, error) {
	if c == nil || c.transport == nil {
		return nil, ErrClientNotConfigured
	}

	return c.transport.GetOpenOrders()
}

func (c *Client) GetTradesSince(ts int64) ([]RawTrade, error) {
	if c == nil || c.transport == nil {
		return nil, ErrClientNotConfigured
	}

	return c.transport.GetTradesSince(ts)
}

func (c *Client) SubmitOrder(req GoexOrderRequest) (GoexOrderResponse, error) {
	if c == nil || c.transport == nil {
		return GoexOrderResponse{}, ErrClientNotConfigured
	}

	return c.transport.SubmitOrder(req)
}

func (c *Client) CancelOrder(exchangeID string) error {
	if c == nil || c.transport == nil {
		return ErrClientNotConfigured
	}

	return c.transport.CancelOrder(exchangeID)
}
