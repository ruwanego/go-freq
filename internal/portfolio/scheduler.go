package portfolio

import (
	"errors"

	"gofreq/internal/marketdata"
)

var ErrSyncTimeout = errors.New("sync timeout")

type Scheduler struct {
	agg        *Aggregator
	timeoutMs  int64
	lastSeenTs map[int64]int64
}

func NewScheduler(pairs []string, timeoutMs int64) *Scheduler {
	return &Scheduler{
		agg:        NewAggregator(pairs),
		timeoutMs:  timeoutMs,
		lastSeenTs: make(map[int64]int64),
	}
}

func (s *Scheduler) OnCandle(c marketdata.Candle, now int64) (PortfolioTick, bool, error) {
	ts := c.Timestamp

	if _, ok := s.lastSeenTs[ts]; !ok {
		s.lastSeenTs[ts] = now
	}

	tick, ready := s.agg.Add(c)
	if ready {
		delete(s.lastSeenTs, ts)
		return tick, true, nil
	}

	if now-s.lastSeenTs[ts] > s.timeoutMs {
		delete(s.lastSeenTs, ts)
		s.agg.DropTimestamp(ts)
		return PortfolioTick{}, false, ErrSyncTimeout
	}

	return PortfolioTick{}, false, nil
}
