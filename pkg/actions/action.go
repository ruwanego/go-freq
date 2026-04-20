package actions

import "gofreq/pkg/types"

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
    Price  float64
    Amount float64
    Tag    string
}
