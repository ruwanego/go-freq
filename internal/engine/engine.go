package engine

import (
	internalexec "gofreq/internal/execution"
	"gofreq/pkg/actions"
	pkgexec "gofreq/pkg/execution"
	goctx "gofreq/pkg/context"
)

type Engine struct {
	strategy Strategy
	pipeline *internalexec.Pipeline
	executor Executor

	warmupRemaining int
	lastResult      pkgexec.ExecutionResult
}

func NewEngine(strategy Strategy, pipeline *internalexec.Pipeline, executor Executor, warmupTicks int) *Engine {
	return &Engine{
		strategy:        strategy,
		pipeline:        pipeline,
		executor:        executor,
		warmupRemaining: warmupTicks,
		lastResult: pkgexec.ExecutionResult{
			Accepted: []actions.Action{},
			Rejected: []pkgexec.RejectedAction{},
		},
	}
}

func (e *Engine) ProcessTick(_ Tick) error {
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

	return e.executor.Execute(result.Accepted)
}
