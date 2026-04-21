package engine

import "gofreq/pkg/actions"

type Executor interface {
	Execute([]actions.Action) error
}
