package engine_test

import (
	"testing"

	eng "gofreq/internal/engine"
	"gofreq/internal/execution"
	"gofreq/pkg/actions"
	pkgexec "gofreq/pkg/execution"
	goctx "gofreq/pkg/context"
)

type resultAwareStrategy struct {
	seen []pkgexec.ExecutionResult
}

func (s *resultAwareStrategy) Name() string { return "result-aware" }

func (s *resultAwareStrategy) OnCandle(ctx *goctx.CandleContext) ([]actions.Action, error) {
	s.seen = append(s.seen, ctx.LastExecutionResult())
	return []actions.Action{
		{Pair: "BTC/USDT", Amount: 1, Tag: "a"},
	}, nil
}

func TestLastExecutionResultVisibleOnNextTick(t *testing.T) {
	strat := &resultAwareStrategy{}
	risk := &execution.BasicRisk{MaxPerTrade: 10}
	alloc := &execution.DeterministicAllocator{}
	pipe := execution.NewPipeline(risk, alloc)
	exec := &recordingExecutor{}

	store := newEngineTestStore(t)
	engine := eng.NewEngine(strat, pipe, exec, store, 0)

	_ = engine.ProcessTick(eng.Tick{Pair: "BTC/USDT", Timestamp: 1})
	_ = engine.ProcessTick(eng.Tick{Pair: "BTC/USDT", Timestamp: 2})

	if len(strat.seen) != 2 {
		t.Fatalf("expected two strategy observations")
	}

	if len(strat.seen[1].Accepted) != 1 {
		t.Fatalf("expected previous tick execution result on second tick")
	}
}
