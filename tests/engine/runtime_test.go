package engine_test

import (
	"errors"
	"github.com/shopspring/decimal"
	"testing"

	eng "gofreq/internal/engine"
	"gofreq/internal/execution"
	"gofreq/internal/marketdata"
	"gofreq/pkg/actions"
	goctx "gofreq/pkg/context"
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

type runtimeEngineStrategyStub struct {
	calls int
}

func (s *runtimeEngineStrategyStub) Name() string { return "runtime-engine" }

func (s *runtimeEngineStrategyStub) OnCandle(ctx *goctx.CandleContext) ([]actions.Action, error) {
	s.calls++
	return nil, nil
}

type runtimeEngineExecutorStub struct {
	calls int
	seen  [][]actions.Action
}

func (e *runtimeEngineExecutorStub) Execute(a []actions.Action) error {
	e.calls++
	e.seen = append(e.seen, a)
	return nil
}

func newRuntimeEngine() *eng.Engine {
	pipe := execution.NewPipeline(&execution.BasicRisk{MaxPerTrade: decimal.NewFromInt(100)}, &execution.DeterministicAllocator{})
	return eng.NewEngine(nil, pipe, nil, nil, 0)
}

func TestStartFailsWhenBootstrapNotReady(t *testing.T) {
	engine := newRuntimeEngine()
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateRecovering})
	engine.SetStrategy(&runtimeEngineStrategyStub{})
	engine.SetExecutor(&runtimeEngineExecutorStub{})
	engine.SetMarketData(&runtimeMarketStub{candles: make(chan marketdata.Candle)})

	err := engine.Start([]string{"BTC/USDT"}, "1m")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestStartFailsWhenStrategyMissing(t *testing.T) {
	engine := newRuntimeEngine()
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateReady})
	engine.SetExecutor(&runtimeEngineExecutorStub{})
	engine.SetMarketData(&runtimeMarketStub{candles: make(chan marketdata.Candle)})

	err := engine.Start([]string{"BTC/USDT"}, "1m")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestStartFailsWhenExecutorMissing(t *testing.T) {
	engine := newRuntimeEngine()
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateReady})
	engine.SetStrategy(&runtimeEngineStrategyStub{})
	engine.SetMarketData(&runtimeMarketStub{candles: make(chan marketdata.Candle)})

	err := engine.Start([]string{"BTC/USDT"}, "1m")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestStartFailsWhenMarketDependencyMissing(t *testing.T) {
	engine := newRuntimeEngine()
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateReady})
	engine.SetStrategy(&runtimeEngineStrategyStub{})
	engine.SetExecutor(&runtimeEngineExecutorStub{})

	err := engine.Start([]string{"BTC/USDT"}, "1m")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestStartFailsWhenMarketSubscriptionFails(t *testing.T) {
	engine := newRuntimeEngine()
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateReady})
	engine.SetStrategy(&runtimeEngineStrategyStub{})
	engine.SetExecutor(&runtimeEngineExecutorStub{})
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
	ch <- marketdata.Candle{Pair: "BTC/USDT", Timestamp: 1, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true}
	ch <- marketdata.Candle{Pair: "BTC/USDT", Timestamp: 2, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true}
	close(ch)

	market := &runtimeMarketStub{candles: ch}
	strategy := &runtimeEngineStrategyStub{}
	executor := &runtimeEngineExecutorStub{}

	engine := newRuntimeEngine()
	engine.SetBootstrap(&runtimeBootstrapStub{state: eng.StateReady})
	engine.SetStrategy(strategy)
	engine.SetExecutor(executor)
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
