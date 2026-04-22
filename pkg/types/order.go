package types

import "github.com/shopspring/decimal"

type Side string
type OrderType string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"

	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
)

type Order struct {
	ID     OrderID
	Pair   string
	Side   Side
	Type   OrderType
	Price  decimal.Decimal
	Amount decimal.Decimal
	Tag    string
}
