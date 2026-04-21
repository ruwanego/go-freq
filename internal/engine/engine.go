package engine

import (
	"errors"
	"fmt"

	internalexec "gofreq/internal/execution"
	"gofreq/internal/identity"
	"gofreq/internal/persistence"
	"gofreq/pkg/actions"
	goctx "gofreq/pkg/context"
	pkgexec "gofreq/pkg/execution"
)

type readinessChecker interface {
	State() EngineState
}

type Engine struct {
	strategy  Strategy
	pipeline  *internalexec.Pipeline
	executor  Executor
	store     *persistence.Store
	generator *identity.Generator
	ready     readinessChecker

	warmupRemaining int
	lastResult      pkgexec.ExecutionResult
}

func NewEngine(strategy Strategy, pipeline *internalexec.Pipeline, executor Executor, store *persistence.Store, warmupTicks int) *Engine {
	return &Engine{
		strategy:        strategy,
		pipeline:        pipeline,
		executor:        executor,
		store:           store,
		generator:       identity.NewGenerator("GF", map[string]int64{}),
		warmupRemaining: warmupTicks,
		lastResult: pkgexec.ExecutionResult{
			Accepted: []actions.Action{},
			Rejected: []pkgexec.RejectedAction{},
		},
	}
}

func (e *Engine) SetBootstrap(b readinessChecker) {
	e.ready = b
}

func (e *Engine) ProcessTick(tick Tick) error {
	if e.ready != nil && e.ready.State() != StateReady {
		return errors.New("engine_not_ready")
	}

	warmup := e.warmupRemaining > 0
	ctx := goctx.NewCandleContext(warmup, e.lastResult)

	proposed, err := e.strategy.OnCandle(ctx)
	if err != nil {
		return err
	}

	if warmup {
		e.lastResult = pkgexec.ExecutionResult{
			Accepted: []actions.Action{},
			Rejected: []pkgexec.RejectedAction{},
		}
		e.warmupRemaining--
		return nil
	}

	result, err := e.pipeline.Process(proposed)
	e.lastResult = result
	if err != nil {
		return err
	}

	if len(result.Accepted) == 0 {
		return nil
	}

	if e.store == nil {
		return fmt.Errorf("missing_store")
	}

	strategyName := e.strategy.Name()

	for _, action := range result.Accepted {
		rec := buildOrderRecord(e.generator, strategyName, action, tick.Timestamp)
		if err := e.store.CreateOrder(rec); err != nil {
			return err
		}

		if err := e.executor.Execute([]actions.Action{action}); err != nil {
			return err
		}

		if err := e.store.UpdateOrderState(rec.EngineID, persistence.OrderStateSubmitted, tick.Timestamp, ""); err != nil {
			return err
		}
	}

	return nil
}

func buildOrderRecord(generator *identity.Generator, strategyName string, action actions.Action, now int64) persistence.OrderRecord {
	id := generator.Next(strategyName, now)

	return persistence.OrderRecord{
		EngineID:      id,
		ClientOrderID: id,
		StrategyName:  strategyName,
		Pair:          action.Pair,
		Side:          string(action.Side),
		Price:         action.Price,
		Amount:        action.Amount,
		Tag:           action.Tag,
		State:         persistence.OrderStatePending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
