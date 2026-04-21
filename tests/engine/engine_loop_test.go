package engine_test

import (
	"testing"

	eng "gofreq/internal/engine"
	"gofreq/internal/execution"
	"gofreq/pkg/actions"
	goctx "gofreq/pkg/context"
)

type fixedStrategy struct{}

func (s *fixedStrategy) Name() string { return "fixed" }

func (s *fixedStrategy) OnCandle(ctx *goctx.CandleContext) ([]actions.Action, error) {
	return []actions.Action{
		{Pair: "BTC/USDT", Amount: 1, Tag: "a"},
	}, nil
}

func TestEngineExecutesAcceptedActionsAfterWarmup(t *testing.T) {
	strat := &fixedStrategy{}
	risk := &execution.BasicRisk{MaxPerTrade: 10}
	alloc := &execution.DeterministicAllocator{}
	pipe := execution.NewPipeline(risk, alloc)
	exec := &recordingExecutor{}

	store := newEngineTestStore(t)
	engine := eng.NewEngine(strat, pipe, exec, store, 0)

	err := engine.ProcessTick(eng.Tick{Pair: "BTC/USDT", Timestamp: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exec.calls != 1 {
		t.Fatalf("expected executor to be called once")
	}

	if len(exec.seen) != 1 || len(exec.seen[0]) != 1 {
		t.Fatalf("expected one accepted action batch")
	}
}
