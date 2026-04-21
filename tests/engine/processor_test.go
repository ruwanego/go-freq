package engine_test

import (
	"errors"
	"testing"

	eng "gofreq/internal/engine"
	"gofreq/internal/marketdata"
)

type processorStrategyStub struct {
	action *eng.RuntimeAction
	err    error
	calls  int
	seen   []marketdata.Candle
}

func (s *processorStrategyStub) OnRuntimeCandle(c marketdata.Candle) (*eng.RuntimeAction, error) {
	s.calls++
	s.seen = append(s.seen, c)
	if s.err != nil {
		return nil, s.err
	}
	return s.action, nil
}

type processorExecutorStub struct {
	ack   eng.OrderAck
	err   error
	calls int
	seen  []eng.OrderIntent
}

func (e *processorExecutorStub) SubmitOrder(intent eng.OrderIntent) (eng.OrderAck, error) {
	e.calls++
	e.seen = append(e.seen, intent)
	if e.err != nil {
		return eng.OrderAck{}, e.err
	}
	return e.ack, nil
}

func TestProcessRuntimeTickNoActionNoExecution(t *testing.T) {
	strategy := &processorStrategyStub{}
	executor := &processorExecutorStub{}
	engine := &eng.Engine{}
	engine.SetRuntimeStrategy(strategy)
	engine.SetOrderExecutor(executor)

	err := engine.ProcessRuntimeTick(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 1, Open: 1, High: 1, Low: 1, Close: 1, Closed: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if executor.calls != 0 {
		t.Fatalf("expected no execution")
	}
}

func TestProcessRuntimeTickValidActionSubmitsOrder(t *testing.T) {
	strategy := &processorStrategyStub{
		action: &eng.RuntimeAction{
			ClientOrderID: "cid-1",
			Pair:          "BTC/USDT",
			Side:          "BUY",
			Type:          "LIMIT",
			Price:         60000,
			Amount:        1,
		},
	}
	executor := &processorExecutorStub{}
	engine := &eng.Engine{}
	engine.SetRuntimeStrategy(strategy)
	engine.SetOrderExecutor(executor)

	err := engine.ProcessRuntimeTick(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 1, Open: 1, High: 1, Low: 1, Close: 1, Closed: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if executor.calls != 1 {
		t.Fatalf("expected one execution")
	}
	if len(executor.seen) != 1 || executor.seen[0].ClientOrderID != "cid-1" {
		t.Fatalf("intent not mapped correctly")
	}
}

func TestProcessRuntimeTickStrategyErrorFails(t *testing.T) {
	strategy := &processorStrategyStub{err: errors.New("strategy failed")}
	executor := &processorExecutorStub{}
	engine := &eng.Engine{}
	engine.SetRuntimeStrategy(strategy)
	engine.SetOrderExecutor(executor)

	err := engine.ProcessRuntimeTick(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 1, Open: 1, High: 1, Low: 1, Close: 1, Closed: true})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestProcessRuntimeTickExecutorErrorFails(t *testing.T) {
	strategy := &processorStrategyStub{
		action: &eng.RuntimeAction{
			ClientOrderID: "cid-2",
			Pair:          "BTC/USDT",
			Side:          "SELL",
			Type:          "MARKET",
			Amount:        2,
		},
	}
	executor := &processorExecutorStub{err: errors.New("submit failed")}
	engine := &eng.Engine{}
	engine.SetRuntimeStrategy(strategy)
	engine.SetOrderExecutor(executor)

	err := engine.ProcessRuntimeTick(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 1, Open: 1, High: 1, Low: 1, Close: 1, Closed: true})
	if err == nil {
		t.Fatalf("expected error")
	}
}
