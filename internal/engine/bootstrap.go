package engine

import (
	"errors"

	"gofreq/internal/persistence"
	"gofreq/internal/recovery"
)

type sequenceRecoverer interface {
	RecoverSequences() error
}

type orderLister interface {
	ListOrders() ([]persistence.OrderRecord, error)
}

type classifierRunner interface {
	Run() ([]recovery.ClassifiedOrder, error)
}

type reconcilerRunner interface {
	Run(int64) error
}

type Bootstrap struct {
	store interface {
		sequenceRecoverer
		orderLister
	}
	classifier classifierRunner
	reconciler reconcilerRunner
	mode       RuntimeMode

	state EngineState
}

func NewBootstrap(
	store interface {
		sequenceRecoverer
		orderLister
	},
	classifier classifierRunner,
	reconciler reconcilerRunner,
	mode RuntimeMode,
) *Bootstrap {
	return &Bootstrap{
		store:      store,
		classifier: classifier,
		reconciler: reconciler,
		mode:       mode,
		state:      StateBooting,
	}
}

func (b *Bootstrap) Boot() error {
	b.state = StateRecovering

	if err := b.store.RecoverSequences(); err != nil {
		b.state = StateFailed
		return err
	}

	if _, err := b.classifier.Run(); err != nil {
		b.state = StateFailed
		return err
	}

	if err := b.reconciler.Run(0); err != nil {
		b.state = StateFailed
		return err
	}

	if err := b.sanityCheck(); err != nil {
		b.state = StateFailed
		return err
	}

	b.state = StateReady
	return nil
}

func (b *Bootstrap) State() EngineState {
	return b.state
}

func (b *Bootstrap) sanityCheck() error {
	orders, err := b.store.ListOrders()
	if err != nil {
		return err
	}

	for _, o := range orders {
		switch o.State {
		case persistence.OrderStatePending:
			return errors.New("pending order after recovery")
		case persistence.OrderStateSubmitted:
			if b.mode == ModeLive && o.ExchangeID == "" {
				return errors.New("submitted order missing exchange id")
			}
		}
	}

	return nil
}
