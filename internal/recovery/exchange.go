package recovery

type OpenOrder struct {
	ClientOrderID string
	ExchangeID    string
	Pair          string
}

type Exchange interface {
	GetOpenOrders() ([]OpenOrder, error)
}
