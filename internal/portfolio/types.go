package portfolio

import "gofreq/internal/marketdata"

type PortfolioTick struct {
	Timestamp int64
	Candles   map[string]marketdata.Candle
}
