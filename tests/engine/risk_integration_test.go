package engine_test

import (
	"testing"

	eng "gofreq/internal/engine"
	"gofreq/internal/execution"
	"gofreq/internal/risk"
	"gofreq/pkg/actions"
	goctx "gofreq/pkg/context"
)

type riskActionStrategy struct{}

func (s *riskActionStrategy) Name() string { return "risk-action" }

func (s *riskActionStrategy) OnCandle(ctx *goctx.CandleContext) ([]actions.Action, error) {
	return []actions.Action{{Pair: "BTC/USDT", Price: 100, Amount: 10, Tag: "r1"}}, nil
}

func TestProcessTickRiskRejectionPreventsExecution(t *testing.T) {
	strat := &riskActionStrategy{}
	riskEngine := &execution.BasicRisk{MaxPerTrade: 100}
	alloc := &execution.DeterministicAllocator{}
	pipe := execution.NewPipeline(riskEngine, alloc)
	exec := &recordingExecutor{}
	store := newEngineTestStore(t)

	engine := eng.NewEngine(strat, pipe, exec, store, 0)
	engine.SetRiskManager(risk.NewManager(5, 10000))

	err := engine.ProcessTick(eng.Tick{Pair: "BTC/USDT", Timestamp: 1})
	if err == nil {
		t.Fatalf("expected risk rejection")
	}
	if exec.calls != 0 {
		t.Fatalf("executor must not run when risk rejects")
	}
}

func TestProcessTickRiskApprovalAllowsExecution(t *testing.T) {
	strat := &riskActionStrategy{}
	riskEngine := &execution.BasicRisk{MaxPerTrade: 100}
	alloc := &execution.DeterministicAllocator{}
	pipe := execution.NewPipeline(riskEngine, alloc)
	exec := &recordingExecutor{}
	store := newEngineTestStore(t)

	engine := eng.NewEngine(strat, pipe, exec, store, 0)
	engine.SetRiskManager(risk.NewManager(20, 2000))

	err := engine.ProcessTick(eng.Tick{Pair: "BTC/USDT", Timestamp: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.calls != 1 {
		t.Fatalf("executor should run when risk approves")
	}
}
