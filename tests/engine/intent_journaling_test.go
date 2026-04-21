package engine_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	eng "gofreq/internal/engine"
	"gofreq/internal/execution"
	"gofreq/internal/persistence"
	"gofreq/pkg/actions"
	goctx "gofreq/pkg/context"
)

func newEngineTestStore(t *testing.T) *persistence.Store {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	store, err := persistence.Open(path)
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}

	t.Cleanup(func() {
		_ = store.Close()
		_ = os.Remove(path)
	})

	return store
}

type multiActionStrategy struct{}

func (s *multiActionStrategy) Name() string { return "multi" }

func (s *multiActionStrategy) OnCandle(ctx *goctx.CandleContext) ([]actions.Action, error) {
	return []actions.Action{
		{Pair: "BTC/USDT", Amount: 1, Tag: "a"},
		{Pair: "ETH/USDT", Amount: 2, Tag: "b"},
	}, nil
}

type failOnSecondExecutor struct {
	count int
}

func (e *failOnSecondExecutor) Execute(a []actions.Action) error {
	e.count++
	if e.count == 2 {
		return errors.New("simulated failure")
	}
	return nil
}

func TestIntentPersistedBeforeExecution(t *testing.T) {
	store := newEngineTestStore(t)

	strat := &multiActionStrategy{}
	risk := &execution.BasicRisk{MaxPerTrade: 10}
	alloc := &execution.DeterministicAllocator{}
	pipe := execution.NewPipeline(risk, alloc)

	exec := &failOnSecondExecutor{}

	engine := eng.NewEngine(strat, pipe, exec, store, 0)

	err := engine.ProcessTick(eng.Tick{})
	if err == nil {
		t.Fatalf("expected execution failure")
	}

	o1, err := store.GetOrder("a")
	if err != nil {
		t.Fatalf("missing order a")
	}
	if o1.State != persistence.OrderStateSubmitted {
		t.Fatalf("first order must be submitted")
	}

	o2, err := store.GetOrder("b")
	if err != nil {
		t.Fatalf("missing order b")
	}
	if o2.State != persistence.OrderStatePending {
		t.Fatalf("second order must remain pending")
	}
}

func TestNoExecutionWithoutPersistence(t *testing.T) {
	store := newEngineTestStore(t)

	strat := &multiActionStrategy{}
	risk := &execution.BasicRisk{MaxPerTrade: 10}
	alloc := &execution.DeterministicAllocator{}
	pipe := execution.NewPipeline(risk, alloc)

	exec := &failOnSecondExecutor{}

	engine := eng.NewEngine(strat, pipe, exec, store, 0)

	_ = engine.ProcessTick(eng.Tick{})

	if _, err := store.GetOrder("a"); err != nil {
		t.Fatalf("order a must exist")
	}
	if _, err := store.GetOrder("b"); err != nil {
		t.Fatalf("order b must exist")
	}
}
