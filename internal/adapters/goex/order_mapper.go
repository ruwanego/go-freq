package goex

import (
	"errors"
	"github.com/shopspring/decimal"
)

type OrderIntent struct {
	ClientOrderID string
	Pair          string
	Side          string
	Type          string
	Price         decimal.Decimal
	Amount        decimal.Decimal
	TimeInForce   string
}

type OrderAck struct {
	ClientOrderID string
	ExchangeID    string
}

type GoexOrderRequest struct {
	ClientOrderID string
	Pair          string
	Side          string
	Type          string
	Price         decimal.Decimal
	Amount        decimal.Decimal
}

type GoexOrderResponse struct {
	ExchangeID string
}

var (
	ErrEmptyClientOrderID = errors.New("empty client order id")
	ErrUnsupportedSide    = errors.New("unsupported side")
	ErrUnsupportedType    = errors.New("unsupported order type")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInvalidPrice       = errors.New("invalid price")
	ErrMissingExchangeID  = errors.New("missing exchange id")
	ErrOrderNotFound      = errors.New("order not found for cancellation")
)

func MapIntentToGoex(intent OrderIntent) (GoexOrderRequest, error) {
	if intent.ClientOrderID == "" {
		return GoexOrderRequest{}, ErrEmptyClientOrderID
	}

	if intent.Amount.LessThanOrEqual(decimal.Zero) {
		return GoexOrderRequest{}, ErrInvalidAmount
	}

	side, err := mapSide(intent.Side)
	if err != nil {
		return GoexOrderRequest{}, err
	}

	otype, err := mapType(intent.Type)
	if err != nil {
		return GoexOrderRequest{}, err
	}

	if intent.Type == "LIMIT" && intent.Price.LessThanOrEqual(decimal.Zero) {
		return GoexOrderRequest{}, ErrInvalidPrice
	}

	return GoexOrderRequest{
		ClientOrderID: intent.ClientOrderID,
		Pair:          intent.Pair,
		Side:          side,
		Type:          otype,
		Price:         intent.Price,
		Amount:        intent.Amount,
	}, nil
}

func mapSide(side string) (string, error) {
	switch side {
	case "BUY":
		return "buy", nil
	case "SELL":
		return "sell", nil
	default:
		return "", ErrUnsupportedSide
	}
}

func mapType(t string) (string, error) {
	switch t {
	case "LIMIT":
		return "limit", nil
	case "MARKET":
		return "market", nil
	default:
		return "", ErrUnsupportedType
	}
}
