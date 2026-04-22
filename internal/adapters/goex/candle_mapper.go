package goex

import (
	"errors"
	"github.com/shopspring/decimal"

	"gofreq/internal/marketdata"
)

var ErrInvalidCandle = errors.New("invalid candle")

func MapGoexCandle(
	pair string,
	ts int64,
	open, high, low, close, volume decimal.Decimal,
	closed bool,
) (marketdata.Candle, error) {
	if ts <= 0 {
		return marketdata.Candle{}, ErrInvalidCandle
	}
	if high.LessThan(low) {
		return marketdata.Candle{}, ErrInvalidCandle
	}
	if open.LessThanOrEqual(decimal.Zero) || close.LessThanOrEqual(decimal.Zero) {
		return marketdata.Candle{}, ErrInvalidCandle
	}
	if !closed {
		return marketdata.Candle{}, ErrInvalidCandle
	}

	return marketdata.Candle{
		Pair:      pair,
		Timestamp: ts,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    volume,
		Closed:    true,
	}, nil
}
