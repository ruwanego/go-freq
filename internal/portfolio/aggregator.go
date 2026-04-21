package portfolio

import "gofreq/internal/marketdata"

type Aggregator struct {
	expectedPairs map[string]struct{}
	buffer        map[int64]map[string]marketdata.Candle
}

func NewAggregator(pairs []string) *Aggregator {
	set := make(map[string]struct{}, len(pairs))
	for _, p := range pairs {
		set[p] = struct{}{}
	}

	return &Aggregator{
		expectedPairs: set,
		buffer:        make(map[int64]map[string]marketdata.Candle),
	}
}

func (a *Aggregator) Add(c marketdata.Candle) (PortfolioTick, bool) {
	ts := c.Timestamp

	if _, ok := a.expectedPairs[c.Pair]; !ok {
		return PortfolioTick{}, false
	}

	if _, ok := a.buffer[ts]; !ok {
		a.buffer[ts] = make(map[string]marketdata.Candle)
	}

	a.buffer[ts][c.Pair] = c

	if len(a.buffer[ts]) == len(a.expectedPairs) {
		candles := make(map[string]marketdata.Candle, len(a.buffer[ts]))
		for pair, candle := range a.buffer[ts] {
			candles[pair] = candle
		}

		tick := PortfolioTick{
			Timestamp: ts,
			Candles:   candles,
		}

		delete(a.buffer, ts)
		return tick, true
	}

	return PortfolioTick{}, false
}

func (a *Aggregator) DropTimestamp(ts int64) {
	delete(a.buffer, ts)
}
