package engine_test

import (
	"testing"

	eng "gofreq/internal/engine"
	"gofreq/internal/execution"
	"gofreq/pkg/actions"
	goctx "gofreq/pkg/context"
)

type warmupStrategy struct {
	calls int
}

func (s *warmupStrategy) Name() string { return "warmup" }

func (s *warmupStrategy) OnCandle(ctx *goctx.CandleContext) ([]actions.Action, error) {
	s.calls++
	return []actions.Action{
		{Pair: "BTC/USDT", Amount: 1, Tag: "a"},
	}, nil
}

type recordingExecutor struct {
	calls int
	seen  [][]actions.Action
}

func (e *recordingExecutor) Execute(a []actions.Action) error {
	e.calls++
	e.seen = append(e.seen, a)
	return nil
}

func TestWarmupBlocksExecutionButStillRunsStrategy(t *testing.T) {
	strat := &warmupStrategy{}
	risk := &execution.BasicRisk{MaxPerTrade: 10}
	alloc := &execution.DeterministicAllocator{}
	pipe := execution.NewPipeline(risk, alloc)
	exec := &recordingExecutor{}

	engine := eng.NewEngine(strat, pipe, exec, 2)

	_ = engine.ProcessTick(eng.Tick{Pair: "BTC/USDT", Timestamp: 1})
	_ = engine.ProcessTick(eng.Tick{Pair: "BTC/USDT", Timestamp: 2})

	if strat.calls != 2 {
		t.Fatalf("expected strategy to run during warmup")
	}
	if exec.calls != 0 {
		t.Fatalf("executor must not run during warmup")
	}
}
