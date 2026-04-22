package recovery

import "github.com/shopspring/decimal"

type OpenOrder struct {
	ClientOrderID string
	ExchangeID    string
	Pair          string
}

type Trade struct {
	ClientOrderID string
	Amount        decimal.Decimal
	Price         decimal.Decimal
	Timestamp     int64
}

type Exchange interface {
	GetOpenOrders() ([]OpenOrder, error)
	GetTradesSince(ts int64) ([]Trade, error)
}
