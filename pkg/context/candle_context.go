package context

import (
	"gofreq/internal/marketdata"
	execution "gofreq/pkg/execution"
)

type CandleContext struct {
	warmup              bool
	lastExecutionResult execution.ExecutionResult
	candle              marketdata.Candle
}

func NewCandleContext(warmup bool, last execution.ExecutionResult, candle marketdata.Candle) *CandleContext {
	return &CandleContext{
		warmup:              warmup,
		lastExecutionResult: last,
		candle:              candle,
	}
}

func (c *CandleContext) IsWarmup() bool {
	return c.warmup
}

func (c *CandleContext) LastExecutionResult() execution.ExecutionResult {
	return c.lastExecutionResult
}

func (c *CandleContext) Candle() marketdata.Candle {
	return c.candle
}
