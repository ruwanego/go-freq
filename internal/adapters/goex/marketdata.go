package goex

import (
	"sort"
	"sync"

	"gofreq/internal/marketdata"
)

type MarketDataAdapter struct {
	client *Client
	stop   chan struct{}
	once   sync.Once
}

func NewMarketDataAdapter(c *Client) *MarketDataAdapter {
	return &MarketDataAdapter{
		client: c,
		stop:   make(chan struct{}),
	}
}

func (m *MarketDataAdapter) GetHistoricalCandles(pair string, timeframe string, limit int) ([]marketdata.Candle, error) {
	raw, err := m.client.GetCandles(pair, timeframe, limit)
	if err != nil {
		return nil, err
	}

	out := make([]marketdata.Candle, 0, len(raw))
	for _, r := range raw {
		c, err := MapGoexCandle(r.Pair, r.Timestamp, r.Open, r.High, r.Low, r.Close, r.Volume, r.Closed)
		if err != nil {
			continue
		}
		out = append(out, c)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Timestamp < out[j].Timestamp
	})

	return out, nil
}

func (m *MarketDataAdapter) SubscribeClosedCandles(pairs []string, timeframe string) (<-chan marketdata.Candle, error) {
	rawCh, err := m.client.SubscribeCandles(pairs, timeframe)
	if err != nil {
		return nil, err
	}

	out := make(chan marketdata.Candle)
	go func() {
		defer close(out)

		for {
			select {
			case <-m.stop:
				return
			case r, ok := <-rawCh:
				if !ok {
					return
				}

				c, err := MapGoexCandle(r.Pair, r.Timestamp, r.Open, r.High, r.Low, r.Close, r.Volume, r.Closed)
				if err != nil {
					continue
				}

				out <- c
			}
		}
	}()

	return out, nil
}

func (m *MarketDataAdapter) Close() error {
	m.once.Do(func() {
		close(m.stop)
	})
	return nil
}
