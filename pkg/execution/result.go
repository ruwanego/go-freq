package execution

import "gofreq/pkg/actions"

type RejectedAction struct {
	Action actions.Action
	Reason string
}

type ExecutionResult struct {
	Accepted []actions.Action
	Rejected []RejectedAction
}
