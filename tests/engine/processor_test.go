package engine_test

import (
	"errors"
	"github.com/shopspring/decimal"
	"testing"

	eng "gofreq/internal/engine"
	"gofreq/internal/execution"
	"gofreq/internal/marketdata"
	"gofreq/internal/persistence"
	"gofreq/pkg/actions"
	goctx "gofreq/pkg/context"
	"gofreq/pkg/types"
)

type candleStrategyStub struct {
	actions []actions.Action
	err     error
	calls   int
	seen    []marketdata.Candle
}

func (s *candleStrategyStub) Name() string { return "candle-stub" }

func (s *candleStrategyStub) OnCandle(ctx *goctx.CandleContext) ([]actions.Action, error) {
	s.calls++
	s.seen = append(s.seen, ctx.Candle())
	if s.err != nil {
		return nil, s.err
	}
	return s.actions, nil
}

type trackingExecutorStub struct {
	err   error
	calls int
	seen  [][]actions.Action
}

func (e *trackingExecutorStub) Execute(a []actions.Action) error {
	e.calls++
	e.seen = append(e.seen, a)
	if e.err != nil {
		return e.err
	}
	return nil
}

func newProcessorEngine(t *testing.T, strategy *candleStrategyStub, executor *trackingExecutorStub) *eng.Engine {
	t.Helper()

	riskEngine := &execution.BasicRisk{MaxPerTrade: decimal.NewFromInt(100)}
	alloc := &execution.DeterministicAllocator{}
	pipe := execution.NewPipeline(riskEngine, alloc)
	store := newEngineTestStore(t)

	return eng.NewEngine(strategy, pipe, executor, store, 0)
}

func TestProcessCandleNoActionNoExecution(t *testing.T) {
	strategy := &candleStrategyStub{}
	executor := &trackingExecutorStub{}
	engine := newProcessorEngine(t, strategy, executor)

	err := engine.ProcessCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 1, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if executor.calls != 0 {
		t.Fatalf("expected no execution")
	}
}

func TestProcessCandleValidActionExecutes(t *testing.T) {
	strategy := &candleStrategyStub{
		actions: []actions.Action{{
			Type:   actions.ActionBuy,
			Pair:   "BTC/USDT",
			Side:   types.SideBuy,
			Price:  decimal.NewFromInt(60000),
			Amount: decimal.NewFromInt(1),
			Tag:    "cid-1",
		}},
	}
	executor := &trackingExecutorStub{}
	engine := newProcessorEngine(t, strategy, executor)

	err := engine.ProcessCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 1, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if executor.calls != 1 {
		t.Fatalf("expected one execution")
	}
	if len(executor.seen) != 1 || len(executor.seen[0]) != 1 || executor.seen[0][0].Tag != "cid-1" {
		t.Fatalf("action not executed correctly")
	}
}

func TestProcessCandleStrategyErrorFails(t *testing.T) {
	strategy := &candleStrategyStub{err: errors.New("strategy failed")}
	executor := &trackingExecutorStub{}
	engine := newProcessorEngine(t, strategy, executor)

	err := engine.ProcessCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 1, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestProcessCandleExecutorErrorFails(t *testing.T) {
	strategy := &candleStrategyStub{
		actions: []actions.Action{{
			Type:   actions.ActionSell,
			Pair:   "BTC/USDT",
			Side:   types.SideSell,
			Price:  decimal.NewFromInt(1),
			Amount: decimal.NewFromInt(2),
			Tag:    "cid-2",
		}},
	}
	executor := &trackingExecutorStub{err: errors.New("execute failed")}
	engine := newProcessorEngine(t, strategy, executor)

	err := engine.ProcessCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 1, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestProcessCandlePersistsBeforeExecution(t *testing.T) {
	strategy := &candleStrategyStub{
		actions: []actions.Action{{
			Type:   actions.ActionBuy,
			Pair:   "BTC/USDT",
			Side:   types.SideBuy,
			Price:  decimal.NewFromInt(10),
			Amount: decimal.NewFromInt(1),
			Tag:    "persist",
		}},
	}
	executor := &trackingExecutorStub{}

	riskEngine := &execution.BasicRisk{MaxPerTrade: decimal.NewFromInt(100)}
	alloc := &execution.DeterministicAllocator{}
	pipe := execution.NewPipeline(riskEngine, alloc)
	store := newEngineTestStore(t)
	engine := eng.NewEngine(strategy, pipe, executor, store, 0)

	err := engine.ProcessCandle(marketdata.Candle{Pair: "BTC/USDT", Timestamp: 7, Open: decimal.NewFromInt(1), High: decimal.NewFromInt(1), Low: decimal.NewFromInt(1), Close: decimal.NewFromInt(1), Closed: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	orders, err := store.ListOrders()
	if err != nil {
		t.Fatalf("list orders failed: %v", err)
	}
	if len(orders) != 1 {
		t.Fatalf("expected 1 persisted order, got %d", len(orders))
	}
	if orders[0].State != persistence.OrderStateSubmitted {
		t.Fatalf("expected submitted order state")
	}
}
