package engine_test

import (
	"errors"
	"testing"

	eng "gofreq/internal/engine"
	"gofreq/internal/marketdata"
)

type runtimeBootstrapStub struct {
	state eng.EngineState
}

func (b *runtimeBootstrapStub) State() eng.EngineState {
	return b.state
}

type runtimeMarketStub struct {
	candles chan marketdata.Candle
	err     error
	calls   int
}

func (m *runtimeMarketStub) SubscribeClosedCandles([]string, string) (<-chan marketdata.Candle, error) {
	m.calls++
	if m.err != nil {
		return nil, m.err
	}
	return m.candles, nil
}

type runtimeStrategyStub struct {
	calls int
}

func (s *runtimeStrategyStub) OnRuntimeCandle(c marketdata.Candle) (*eng.RuntimeAction, error) {
	s.calls++
	return nil, nil
}

type runtimeExecutorStub struct {
	calls int
}

func (e *runtimeExecutorStub) SubmitOrder(eng.OrderIntent) (eng.OrderAck, error) {
	e.calls++
	return eng.OrderAck{}, nil
}

func TestStartFailsWhenBootstrapNotReady(t *testing.T) {
	engine := &eng.Engine{}
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateRecovering})
	engine.SetRuntimeStrategy(&runtimeStrategyStub{})
	engine.SetOrderExecutor(&runtimeExecutorStub{})
	engine.SetMarketData(&runtimeMarketStub{candles: make(chan marketdata.Candle)})

	err := engine.Start([]string{"BTC/USDT"}, "1m")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestStartFailsWhenRuntimeStrategyMissing(t *testing.T) {
	engine := &eng.Engine{}
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateReady})
	engine.SetOrderExecutor(&runtimeExecutorStub{})
	engine.SetMarketData(&runtimeMarketStub{candles: make(chan marketdata.Candle)})

	err := engine.Start([]string{"BTC/USDT"}, "1m")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestStartFailsWhenRuntimeExecutorMissing(t *testing.T) {
	engine := &eng.Engine{}
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateReady})
	engine.SetRuntimeStrategy(&runtimeStrategyStub{})
	engine.SetMarketData(&runtimeMarketStub{candles: make(chan marketdata.Candle)})

	err := engine.Start([]string{"BTC/USDT"}, "1m")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestStartFailsWhenMarketDependencyMissing(t *testing.T) {
	engine := &eng.Engine{}
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateReady})
	engine.SetRuntimeStrategy(&runtimeStrategyStub{})
	engine.SetOrderExecutor(&runtimeExecutorStub{})

	err := engine.Start([]string{"BTC/USDT"}, "1m")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestStartFailsWhenMarketSubscriptionFails(t *testing.T) {
	engine := &eng.Engine{}
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateReady})
	engine.SetRuntimeStrategy(&runtimeStrategyStub{})
	engine.SetOrderExecutor(&runtimeExecutorStub{})
	engine.SetMarketData(&runtimeMarketStub{err: errors.New("subscribe failed")})

	err := engine.Start([]string{"BTC/USDT"}, "1m")
	if err == nil {
		t.Fatalf("expected error")
	}
	if engine.State() != eng.StateFailed {
		t.Fatalf("expected failed state")
	}
}

func TestStartProcessesCandlesSequentially(t *testing.T) {
	ch := make(chan marketdata.Candle, 2)
	ch <- marketdata.Candle{Pair: "BTC/USDT", Timestamp: 1, Open: 1, High: 1, Low: 1, Close: 1, Closed: true}
	ch <- marketdata.Candle{Pair: "BTC/USDT", Timestamp: 2, Open: 1, High: 1, Low: 1, Close: 1, Closed: true}
	close(ch)

	market := &runtimeMarketStub{candles: ch}
	strategy := &runtimeStrategyStub{}
	executor := &runtimeExecutorStub{}

	engine := &eng.Engine{}
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateReady})
	engine.SetRuntimeStrategy(strategy)
	engine.SetOrderExecutor(executor)
	engine.SetMarketData(market)

	if err := engine.Start([]string{"BTC/USDT"}, "1m"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if engine.State() != eng.StateRunning {
		t.Fatalf("expected running state")
	}
	if strategy.calls != 2 {
		t.Fatalf("expected 2 strategy calls, got %d", strategy.calls)
	}
	if executor.calls != 0 {
		t.Fatalf("expected no executions")
	}
}
