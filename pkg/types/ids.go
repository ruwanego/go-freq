package types

type EngineID string
type ClientOrderID string
type ExchangeID string

type OrderID struct {
    EngineID      EngineID
    ClientOrderID ClientOrderID
    ExchangeID    ExchangeID
}
