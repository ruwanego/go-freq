package goex_test

import (
	"errors"
	"testing"

	adapter "gofreq/internal/adapters/goex"
)

type fakeExecClient struct {
	submitResp adapter.GoexOrderResponse
	submitErr  error
	cancelErr  error
	orders     []adapter.RawOrder
	submitted  []adapter.GoexOrderRequest
	cancelled  []string
}

func (f *fakeExecClient) GetOpenOrders() ([]adapter.RawOrder, error) {
	if f.submitErr != nil && len(f.orders) == 0 {
		return nil, f.submitErr
	}
	return f.orders, nil
}

func (f *fakeExecClient) GetTradesSince(int64) ([]adapter.RawTrade, error) {
	return nil, nil
}

func (f *fakeExecClient) SubmitOrder(req adapter.GoexOrderRequest) (adapter.GoexOrderResponse, error) {
	f.submitted = append(f.submitted, req)
	if f.submitErr != nil {
		return adapter.GoexOrderResponse{}, f.submitErr
	}
	return f.submitResp, nil
}

func (f *fakeExecClient) CancelOrder(exchangeID string) error {
	f.cancelled = append(f.cancelled, exchangeID)
	if f.cancelErr != nil {
		return f.cancelErr
	}
	return nil
}

func TestExecutorSubmitOrder_Success(t *testing.T) {
	fake := &fakeExecClient{
		submitResp: adapter.GoexOrderResponse{ExchangeID: "ex-1"},
	}
	client := adapter.NewClient(fake)
	exec := adapter.NewExecutor(client)

	ack, err := exec.SubmitOrder(adapter.OrderIntent{
		ClientOrderID: "cid-1",
		Pair:          "BTC/USDT",
		Side:          "BUY",
		Type:          "LIMIT",
		Price:         60000,
		Amount:        1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ack.ClientOrderID != "cid-1" {
		t.Fatalf("client order id mismatch")
	}
	if ack.ExchangeID != "ex-1" {
		t.Fatalf("exchange id mismatch")
	}
	if len(fake.submitted) != 1 || fake.submitted[0].ClientOrderID != "cid-1" {
		t.Fatalf("submit did not preserve client order id")
	}
}

func TestExecutorSubmitOrder_ErrorPropagates(t *testing.T) {
	expectedErr := errors.New("submit failed")
	client := adapter.NewClient(&fakeExecClient{submitErr: expectedErr})
	exec := adapter.NewExecutor(client)

	_, err := exec.SubmitOrder(adapter.OrderIntent{
		ClientOrderID: "cid-2",
		Pair:          "BTC/USDT",
		Side:          "BUY",
		Type:          "MARKET",
		Amount:        1,
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestExecutorSubmitOrder_MissingExchangeIDFails(t *testing.T) {
	client := adapter.NewClient(&fakeExecClient{submitResp: adapter.GoexOrderResponse{ExchangeID: ""}})
	exec := adapter.NewExecutor(client)

	_, err := exec.SubmitOrder(adapter.OrderIntent{
		ClientOrderID: "cid-3",
		Pair:          "BTC/USDT",
		Side:          "BUY",
		Type:          "MARKET",
		Amount:        1,
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestExecutorCancelOrder_ResolvesExchangeID(t *testing.T) {
	fake := &fakeExecClient{
		orders: []adapter.RawOrder{{ClientOrderID: "cid-4", ExchangeID: "ex-4", Pair: "BTC/USDT"}},
	}
	client := adapter.NewClient(fake)
	exec := adapter.NewExecutor(client)

	if err := exec.CancelOrder("cid-4"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.cancelled) != 1 || fake.cancelled[0] != "ex-4" {
		t.Fatalf("cancel did not use resolved exchange id")
	}
}

func TestExecutorCancelOrder_OrderNotFound(t *testing.T) {
	client := adapter.NewClient(&fakeExecClient{orders: []adapter.RawOrder{}})
	exec := adapter.NewExecutor(client)

	err := exec.CancelOrder("missing")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestExecutorCancelOrder_EmptyIDFails(t *testing.T) {
	client := adapter.NewClient(&fakeExecClient{})
	exec := adapter.NewExecutor(client)

	err := exec.CancelOrder("")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestExecutorCancelOrder_ClientErrorPropagates(t *testing.T) {
	expectedErr := errors.New("cancel failed")
	fake := &fakeExecClient{
		orders:    []adapter.RawOrder{{ClientOrderID: "cid-5", ExchangeID: "ex-5", Pair: "BTC/USDT"}},
		cancelErr: expectedErr,
	}
	client := adapter.NewClient(fake)
	exec := adapter.NewExecutor(client)

	err := exec.CancelOrder("cid-5")
	if err == nil {
		t.Fatalf("expected error")
	}
}
