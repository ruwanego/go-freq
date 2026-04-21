package persistence

import "fmt"

func ValidateTransition(from, to OrderState) error {
	valid := map[OrderState]map[OrderState]bool{
		OrderStatePending: {
			OrderStateSubmitted: true,
			OrderStateCancelled: true,
		},
		OrderStateSubmitted: {
			OrderStateFilled:    true,
			OrderStateCancelled: true,
		},
		OrderStateFilled:    {},
		OrderStateCancelled: {},
	}

	nexts, ok := valid[from]
	if !ok {
		return fmt.Errorf("unknown_from_state:%s", from)
	}
	if !nexts[to] {
		return fmt.Errorf("invalid_transition:%s->%s", from, to)
	}
	return nil
}
