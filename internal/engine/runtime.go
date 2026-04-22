package engine

import (
	"errors"

	"gofreq/internal/marketdata"
)

var (
	ErrEngineNotReady = errors.New("engine not ready")
	ErrStrategyMissing = errors.New("strategy missing")
	ErrExecutorMissing = errors.New("executor missing")
	ErrMarketMissing   = errors.New("market missing")
)

type marketDataSource interface {
	SubscribeClosedCandles(pairs []string, timeframe string) (<-chan marketdata.Candle, error)
}

func (e *Engine) SetStrategy(s Strategy) {
	e.strategy = s
}

func (e *Engine) SetExecutor(ex Executor) {
	e.executor = ex
}

func (e *Engine) SetMarketData(m marketDataSource) {
	e.market = m
}

func (e *Engine) Start(pairs []string, timeframe string) error {
	if e.ready == nil || e.ready.State() != StateReady {
		return ErrEngineNotReady
	}
	if e.strategy == nil {
		return ErrStrategyMissing
	}
	if e.executor == nil {
		return ErrExecutorMissing
	}
	if e.market == nil {
		return ErrMarketMissing
	}

	e.state = StateRunning

	ch, err := e.market.SubscribeClosedCandles(pairs, timeframe)
	if err != nil {
		e.state = StateFailed
		return err
	}

	for candle := range ch {
		if err := e.ProcessCandle(candle); err != nil {
			e.state = StateFailed
			return err
		}
	}

	return nil
}
