package engine

import (
	"fmt"

	internalexec "gofreq/internal/execution"
	"gofreq/internal/persistence"
	"gofreq/pkg/actions"
	goctx "gofreq/pkg/context"
	pkgexec "gofreq/pkg/execution"
)

type Engine struct {
	strategy Strategy
	pipeline *internalexec.Pipeline
	executor Executor
	store    *persistence.Store

	warmupRemaining int
	lastResult      pkgexec.ExecutionResult
}

func NewEngine(strategy Strategy, pipeline *internalexec.Pipeline, executor Executor, store *persistence.Store, warmupTicks int) *Engine {
	return &Engine{
		strategy:        strategy,
		pipeline:        pipeline,
		executor:        executor,
		store:           store,
		warmupRemaining: warmupTicks,
		lastResult: pkgexec.ExecutionResult{
			Accepted: []actions.Action{},
			Rejected: []pkgexec.RejectedAction{},
		},
	}
}

func (e *Engine) ProcessTick(tick Tick) error {
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

	for _, action := range result.Accepted {
		rec := buildOrderRecord(action, tick.Timestamp)
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

func buildOrderRecord(action actions.Action, now int64) persistence.OrderRecord {
	return persistence.OrderRecord{
		EngineID:      action.Tag,
		ClientOrderID: action.Tag,
		StrategyName:  "TODO",
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
