package engine

import (
	"errors"

	"gofreq/internal/marketdata"
)

var (
	ErrEngineNotReady         = errors.New("engine not ready")
	ErrRuntimeStrategyMissing = errors.New("runtime strategy missing")
	ErrRuntimeExecutorMissing = errors.New("runtime executor missing")
	ErrRuntimeMarketMissing   = errors.New("runtime market missing")
)

type runtimeStrategy interface {
	OnRuntimeCandle(c marketdata.Candle) (*RuntimeAction, error)
}

type runtimeOrderExecutor interface {
	SubmitOrder(intent OrderIntent) (OrderAck, error)
}

type runtimeMarketData interface {
	SubscribeClosedCandles(pairs []string, timeframe string) (<-chan marketdata.Candle, error)
}

const StateRunning EngineState = "RUNNING"

type RuntimeAction struct {
	ClientOrderID string
	Pair          string
	Side          string
	Type          string
	Price         float64
	Amount        float64
}

type OrderIntent struct {
	ClientOrderID string
	Pair          string
	Side          string
	Type          string
	Price         float64
	Amount        float64
}

type OrderAck struct {
	ClientOrderID string
	ExchangeID    string
}

func (e *Engine) SetRuntimeStrategy(s runtimeStrategy) {
	e.runtimeStrategy = s
}

func (e *Engine) SetOrderExecutor(ex runtimeOrderExecutor) {
	e.runtimeExecutor = ex
}

func (e *Engine) SetMarketData(m runtimeMarketData) {
	e.runtimeMarket = m
}

func (e *Engine) Start(pairs []string, timeframe string) error {
	if e.ready == nil || e.ready.State() != StateReady {
		return ErrEngineNotReady
	}
	if e.runtimeStrategy == nil {
		return ErrRuntimeStrategyMissing
	}
	if e.runtimeExecutor == nil {
		return ErrRuntimeExecutorMissing
	}
	if e.runtimeMarket == nil {
		return ErrRuntimeMarketMissing
	}

	e.state = StateRunning

	ch, err := e.runtimeMarket.SubscribeClosedCandles(pairs, timeframe)
	if err != nil {
		e.state = StateFailed
		return err
	}

	for candle := range ch {
		if err := e.ProcessRuntimeTick(candle); err != nil {
			e.state = StateFailed
			return err
		}
	}

	return nil
}
