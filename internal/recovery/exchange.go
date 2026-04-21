package recovery

type OpenOrder struct {
	ClientOrderID string
	ExchangeID    string
	Pair          string
}

type Trade struct {
	ClientOrderID string
	Amount        float64
	Price         float64
	Timestamp     int64
}

type Exchange interface {
	GetOpenOrders() ([]OpenOrder, error)
	GetTradesSince(ts int64) ([]Trade, error)
}
