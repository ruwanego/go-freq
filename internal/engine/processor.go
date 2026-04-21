package engine

import "gofreq/internal/marketdata"

func (e *Engine) ProcessRuntimeTick(c marketdata.Candle) error {
	if e.runtimeStrategy == nil {
		return ErrRuntimeStrategyMissing
	}
	if e.runtimeExecutor == nil {
		return ErrRuntimeExecutorMissing
	}

	action, err := e.runtimeStrategy.OnRuntimeCandle(c)
	if err != nil {
		return err
	}

	if action == nil {
		return nil
	}

	intent := OrderIntent{
		ClientOrderID: action.ClientOrderID,
		Pair:          action.Pair,
		Side:          action.Side,
		Type:          action.Type,
		Price:         action.Price,
		Amount:        action.Amount,
	}

	_, err = e.runtimeExecutor.SubmitOrder(intent)
	if err != nil {
		return err
	}

	return nil
}
