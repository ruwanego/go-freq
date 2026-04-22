package goex

import (
	"errors"
	"github.com/shopspring/decimal"

	"gofreq/internal/recovery"
)

var ErrMissingClientOrderID = errors.New("missing client order id")

func MapOrderToOpenOrder(clientOrderID string, exchangeID string, pair string) (recovery.OpenOrder, error) {
	if clientOrderID == "" {
		return recovery.OpenOrder{}, ErrMissingClientOrderID
	}

	return recovery.OpenOrder{
		ClientOrderID: clientOrderID,
		ExchangeID:    exchangeID,
		Pair:          pair,
	}, nil
}

func MapTradeToTrade(clientOrderID string, amount decimal.Decimal, price decimal.Decimal, ts int64) (recovery.Trade, error) {
	if clientOrderID == "" {
		return recovery.Trade{}, ErrMissingClientOrderID
	}

	return recovery.Trade{
		ClientOrderID: clientOrderID,
		Amount:        amount,
		Price:         price,
		Timestamp:     ts,
	}, nil
}
