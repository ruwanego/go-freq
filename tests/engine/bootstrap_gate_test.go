package engine_test

import (
	"testing"

	eng "gofreq/internal/engine"
	"gofreq/internal/execution"
	"gofreq/pkg/actions"
	goctx "gofreq/pkg/context"
)

type gatedStrategy struct {
	calls int
}

func (s *gatedStrategy) Name() string { return "gated" }

func (s *gatedStrategy) OnCandle(ctx *goctx.CandleContext) ([]actions.Action, error) {
	s.calls++
	return []actions.Action{{Pair: "BTC/USDT", Amount: 1, Tag: "a"}}, nil
}

type bootstrapStateStub struct {
	state eng.EngineState
}

func (b *bootstrapStateStub) State() eng.EngineState {
	return b.state
}

func TestProcessTickBlockedWhenBootstrapNotReady(t *testing.T) {
	strat := &gatedStrategy{}
	risk := &execution.BasicRisk{MaxPerTrade: 10}
	alloc := &execution.DeterministicAllocator{}
	pipe := execution.NewPipeline(risk, alloc)
	exec := &recordingExecutor{}
	store := newEngineTestStore(t)

	engine := eng.NewEngine(strat, pipe, exec, store, 0)
	engine.SetBootstrap(&bootstrapStateStub{state: eng.StateRecovering})

	err := engine.ProcessTick(eng.Tick{Pair: "BTC/USDT", Timestamp: 1})
	if err == nil {
		t.Fatalf("expected engine_not_ready error")
	}
	if strat.calls != 0 {
		t.Fatalf("strategy must not run before ready")
	}
	if exec.calls != 0 {
		t.Fatalf("executor must not run before ready")
	}
}

func TestProcessTickAllowedWhenBootstrapReady(t *testing.T) {
	strat := &gatedStrategy{}
	risk := &execution.BasicRisk{MaxPerTrade: 10}
	alloc := &execution.DeterministicAllocator{}
	pipe := execution.NewPipeline(risk, alloc)
	exec := &recordingExecutor{}
	store := newEngineTestStore(t)

	engine := eng.NewEngine(strat, pipe, exec, store, 0)
	engine.SetBootstrap(&bootstrapStateStub{state: eng.StateReady})

	err := engine.ProcessTick(eng.Tick{Pair: "BTC/USDT", Timestamp: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strat.calls != 1 {
		t.Fatalf("strategy should run after ready")
	}
}
