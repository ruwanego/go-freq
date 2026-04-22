package goex

import (
	"errors"
	"github.com/shopspring/decimal"
)

var ErrClientNotConfigured = errors.New("client_not_configured")

type RawOrder struct {
	ClientOrderID string
	ExchangeID    string
	Pair          string
}

type RawTrade struct {
	ClientOrderID string
	Amount        decimal.Decimal
	Price         decimal.Decimal
	Timestamp     int64
}

type RawCandle struct {
	Pair      string
	Timestamp int64
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Volume    decimal.Decimal
	Closed    bool
}

type Transport interface {
	GetOpenOrders() ([]RawOrder, error)
	GetTradesSince(ts int64) ([]RawTrade, error)
	SubmitOrder(req GoexOrderRequest) (GoexOrderResponse, error)
	CancelOrder(exchangeID string) error
	GetCandles(pair, timeframe string, limit int) ([]RawCandle, error)
	SubscribeCandles(pairs []string, timeframe string) (<-chan RawCandle, error)
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

func (c *Client) GetCandles(pair, timeframe string, limit int) ([]RawCandle, error) {
	if c == nil || c.transport == nil {
		return nil, ErrClientNotConfigured
	}

	return c.transport.GetCandles(pair, timeframe, limit)
}

func (c *Client) SubscribeCandles(pairs []string, timeframe string) (<-chan RawCandle, error) {
	if c == nil || c.transport == nil {
		return nil, ErrClientNotConfigured
	}

	return c.transport.SubscribeCandles(pairs, timeframe)
}
