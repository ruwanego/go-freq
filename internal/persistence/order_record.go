package persistence

import "github.com/shopspring/decimal"

type OrderState string

const (
	OrderStatePending         OrderState = "PENDING"
	OrderStateSubmitted       OrderState = "SUBMITTED"
	OrderStatePartiallyFilled OrderState = "PARTIALLY_FILLED"
	OrderStateFilled          OrderState = "FILLED"
	OrderStateCancelled       OrderState = "CANCELLED"
)

type OrderRecord struct {
	EngineID      string `json:"engine_id"`
	ClientOrderID string `json:"client_order_id"`
	ExchangeID    string `json:"exchange_id"`

	StrategyName string          `json:"strategy_name"`
	Pair         string          `json:"pair"`
	Side         string          `json:"side"`
	Price        decimal.Decimal `json:"price"`
	Amount       decimal.Decimal `json:"amount"`
	Tag          string          `json:"tag"`

	State OrderState `json:"state"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}
