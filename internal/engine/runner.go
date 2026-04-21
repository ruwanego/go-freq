package engine

import (
	"gofreq/pkg/actions"
	goctx "gofreq/pkg/context"
)

type Strategy interface {
	Name() string
	OnCandle(*goctx.CandleContext) ([]actions.Action, error)
}
