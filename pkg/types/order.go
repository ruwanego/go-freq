package types

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
    Price  float64
    Amount float64
    Tag    string
}
