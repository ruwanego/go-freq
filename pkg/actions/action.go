package actions

import (
	"github.com/shopspring/decimal"
	"gofreq/pkg/types"
)

type ActionType string

const (
	ActionBuy    ActionType = "BUY"
	ActionSell   ActionType = "SELL"
	ActionCancel ActionType = "CANCEL"
)

type Action struct {
	Type   ActionType
	Pair   string
	Side   types.Side
	Price  decimal.Decimal
	Amount decimal.Decimal
	Tag    string
}
