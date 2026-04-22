package marketdata

import "github.com/shopspring/decimal"

type Candle struct {
	Pair      string
	Timestamp int64
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Volume    decimal.Decimal
	Closed    bool
}
