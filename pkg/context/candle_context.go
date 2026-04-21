package context

import execution "gofreq/pkg/execution"

type CandleContext struct {
	warmup              bool
	lastExecutionResult execution.ExecutionResult
}

func NewCandleContext(warmup bool, last execution.ExecutionResult) *CandleContext {
	return &CandleContext{
		warmup:              warmup,
		lastExecutionResult: last,
	}
}

func (c *CandleContext) IsWarmup() bool {
	return c.warmup
}

func (c *CandleContext) LastExecutionResult() execution.ExecutionResult {
	return c.lastExecutionResult
}
