package engine

import "gofreq/internal/marketdata"

func (e *Engine) ProcessRuntimeTick(c marketdata.Candle) error {
	return e.ProcessCandle(c)
}
