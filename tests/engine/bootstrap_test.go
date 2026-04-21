package engine_test

import (
	"errors"
	"testing"

	eng "gofreq/internal/engine"
	"gofreq/internal/persistence"
	"gofreq/internal/recovery"
)

type bootstrapStoreMock struct {
	recoverErr error
	listErr    error
	orders     []persistence.OrderRecord
	recoverCnt int
	listCnt    int
}

func (m *bootstrapStoreMock) RecoverSequences() error {
	m.recoverCnt++
	return m.recoverErr
}

func (m *bootstrapStoreMock) ListOrders() ([]persistence.OrderRecord, error) {
	m.listCnt++
	return m.orders, m.listErr
}

type classifierMock struct {
	err error
	cnt int
}

func (m *classifierMock) Run() ([]recovery.ClassifiedOrder, error) {
	m.cnt++
	if m.err != nil {
		return nil, m.err
	}
	return []recovery.ClassifiedOrder{}, nil
}

type reconcilerMock struct {
	err error
	cnt int
}

func (m *reconcilerMock) Run(int64) error {
	m.cnt++
	return m.err
}

func TestBootstrapBootSuccess(t *testing.T) {
	store := &bootstrapStoreMock{orders: []persistence.OrderRecord{}}
	classifier := &classifierMock{}
	reconciler := &reconcilerMock{}

	b := eng.NewBootstrap(store, classifier, reconciler, eng.ModeLive)

	if err := b.Boot(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if b.State() != eng.StateReady {
		t.Fatalf("expected READY, got %s", b.State())
	}
}

func TestBootstrapBootRecoverSequencesFails(t *testing.T) {
	store := &bootstrapStoreMock{recoverErr: errors.New("recover failed")}
	classifier := &classifierMock{}
	reconciler := &reconcilerMock{}

	b := eng.NewBootstrap(store, classifier, reconciler, eng.ModeLive)

	if err := b.Boot(); err == nil {
		t.Fatalf("expected error")
	}

	if b.State() != eng.StateFailed {
		t.Fatalf("expected FAILED, got %s", b.State())
	}
}

func TestBootstrapBootClassifierFails(t *testing.T) {
	store := &bootstrapStoreMock{orders: []persistence.OrderRecord{}}
	classifier := &classifierMock{err: errors.New("classify failed")}
	reconciler := &reconcilerMock{}

	b := eng.NewBootstrap(store, classifier, reconciler, eng.ModeLive)

	if err := b.Boot(); err == nil {
		t.Fatalf("expected error")
	}

	if b.State() != eng.StateFailed {
		t.Fatalf("expected FAILED, got %s", b.State())
	}
}

func TestBootstrapBootReconcilerFails(t *testing.T) {
	store := &bootstrapStoreMock{orders: []persistence.OrderRecord{}}
	classifier := &classifierMock{}
	reconciler := &reconcilerMock{err: errors.New("reconcile failed")}

	b := eng.NewBootstrap(store, classifier, reconciler, eng.ModeLive)

	if err := b.Boot(); err == nil {
		t.Fatalf("expected error")
	}

	if b.State() != eng.StateFailed {
		t.Fatalf("expected FAILED, got %s", b.State())
	}
}

func TestBootstrapBootSanityCheckFails(t *testing.T) {
	store := &bootstrapStoreMock{orders: []persistence.OrderRecord{{
		EngineID:   "GF-test-1-0001",
		State:      persistence.OrderStatePending,
		UpdatedAt:  1,
		CreatedAt:  1,
		Amount:     1,
		Pair:       "BTC/USDT",
		ExchangeID: "",
	}}}
	classifier := &classifierMock{}
	reconciler := &reconcilerMock{}

	b := eng.NewBootstrap(store, classifier, reconciler, eng.ModeLive)

	if err := b.Boot(); err == nil {
		t.Fatalf("expected error")
	}

	if b.State() != eng.StateFailed {
		t.Fatalf("expected FAILED, got %s", b.State())
	}
}

func TestBootstrapBootIdempotent(t *testing.T) {
	store := &bootstrapStoreMock{orders: []persistence.OrderRecord{}}
	classifier := &classifierMock{}
	reconciler := &reconcilerMock{}

	b := eng.NewBootstrap(store, classifier, reconciler, eng.ModeLive)

	if err := b.Boot(); err != nil {
		t.Fatalf("first boot unexpected error: %v", err)
	}
	if err := b.Boot(); err != nil {
		t.Fatalf("second boot unexpected error: %v", err)
	}

	if b.State() != eng.StateReady {
		t.Fatalf("expected READY, got %s", b.State())
	}
	if store.recoverCnt != 2 {
		t.Fatalf("expected two sequence recovery calls, got %d", store.recoverCnt)
	}
	if classifier.cnt != 2 {
		t.Fatalf("expected two classifier calls, got %d", classifier.cnt)
	}
	if reconciler.cnt != 2 {
		t.Fatalf("expected two reconciler calls, got %d", reconciler.cnt)
	}
}
