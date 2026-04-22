package goex_test

import (
	"errors"
	"github.com/shopspring/decimal"
	"testing"

	adapter "gofreq/internal/adapters/goex"
)

type fakeClient struct {
	orders []adapter.RawOrder
	trades []adapter.RawTrade
	err    error
}

func (f *fakeClient) GetOpenOrders() ([]adapter.RawOrder, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.orders, nil
}

func (f *fakeClient) GetTradesSince(int64) ([]adapter.RawTrade, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.trades, nil
}

func (f *fakeClient) SubmitOrder(adapter.GoexOrderRequest) (adapter.GoexOrderResponse, error) {
	if f.err != nil {
		return adapter.GoexOrderResponse{}, f.err
	}
	return adapter.GoexOrderResponse{}, nil
}

func (f *fakeClient) CancelOrder(string) error {
	return f.err
}

func (f *fakeClient) GetCandles(string, string, int) ([]adapter.RawCandle, error) {
	if f.err != nil {
		return nil, f.err
	}
	return nil, nil
}

func (f *fakeClient) SubscribeCandles([]string, string) (<-chan adapter.RawCandle, error) {
	if f.err != nil {
		return nil, f.err
	}
	ch := make(chan adapter.RawCandle)
	close(ch)
	return ch, nil
}

func TestRecoveryAdapterGetOpenOrders_NormalMapping(t *testing.T) {
	client := adapter.NewClient(&fakeClient{
		orders: []adapter.RawOrder{
			{ClientOrderID: "cid-1", ExchangeID: "ex-1", Pair: "BTC/USDT"},
		},
	})

	recoveryAdapter := adapter.NewRecoveryAdapter(client)
	got, err := recoveryAdapter.GetOpenOrders()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 order, got %d", len(got))
	}
	if got[0].ClientOrderID != "cid-1" {
		t.Fatalf("client order id mismatch")
	}
	if got[0].ExchangeID != "ex-1" {
		t.Fatalf("exchange id mismatch")
	}
	if got[0].Pair != "BTC/USDT" {
		t.Fatalf("pair mismatch")
	}
}

func TestRecoveryAdapterGetOpenOrders_SkipInvalid(t *testing.T) {
	client := adapter.NewClient(&fakeClient{
		orders: []adapter.RawOrder{
			{ClientOrderID: "cid-1", ExchangeID: "ex-1", Pair: "BTC/USDT"},
			{ClientOrderID: "", ExchangeID: "ex-2", Pair: "ETH/USDT"},
		},
	})

	recoveryAdapter := adapter.NewRecoveryAdapter(client)
	got, err := recoveryAdapter.GetOpenOrders()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 valid order, got %d", len(got))
	}
	if got[0].ClientOrderID != "cid-1" {
		t.Fatalf("unexpected order preserved")
	}
}

func TestRecoveryAdapterGetOpenOrders_EmptyList(t *testing.T) {
	client := adapter.NewClient(&fakeClient{orders: []adapter.RawOrder{}})

	recoveryAdapter := adapter.NewRecoveryAdapter(client)
	got, err := recoveryAdapter.GetOpenOrders()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty list")
	}
}

func TestRecoveryAdapterGetOpenOrders_ErrorPropagation(t *testing.T) {
	expectedErr := errors.New("boom")
	client := adapter.NewClient(&fakeClient{err: expectedErr})

	recoveryAdapter := adapter.NewRecoveryAdapter(client)
	_, err := recoveryAdapter.GetOpenOrders()
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRecoveryAdapterGetTradesSince_NormalMapping(t *testing.T) {
	client := adapter.NewClient(&fakeClient{
		trades: []adapter.RawTrade{
			{ClientOrderID: "cid-2", Amount: decimal.RequireFromString("1.5"), Price: decimal.NewFromInt(100), Timestamp: 123},
		},
	})

	recoveryAdapter := adapter.NewRecoveryAdapter(client)
	got, err := recoveryAdapter.GetTradesSince(100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(got))
	}
	if got[0].ClientOrderID != "cid-2" {
		t.Fatalf("client order id mismatch")
	}
	if !got[0].Amount.Equal(decimal.RequireFromString("1.5")) {
		t.Fatalf("amount mismatch")
	}
	if !got[0].Price.Equal(decimal.NewFromInt(100)) {
		t.Fatalf("price mismatch")
	}
	if got[0].Timestamp != 123 {
		t.Fatalf("timestamp mismatch")
	}
}

func TestRecoveryAdapterGetTradesSince_PartialSkip(t *testing.T) {
	client := adapter.NewClient(&fakeClient{
		trades: []adapter.RawTrade{
			{ClientOrderID: "cid-2", Amount: decimal.RequireFromString("1.5"), Price: decimal.NewFromInt(100), Timestamp: 123},
			{ClientOrderID: "", Amount: decimal.RequireFromString("2.0"), Price: decimal.NewFromInt(200), Timestamp: 124},
		},
	})

	recoveryAdapter := adapter.NewRecoveryAdapter(client)
	got, err := recoveryAdapter.GetTradesSince(100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 valid trade, got %d", len(got))
	}
	if got[0].ClientOrderID != "cid-2" {
		t.Fatalf("unexpected trade preserved")
	}
}

func TestRecoveryAdapterGetTradesSince_EmptyList(t *testing.T) {
	client := adapter.NewClient(&fakeClient{trades: []adapter.RawTrade{}})

	recoveryAdapter := adapter.NewRecoveryAdapter(client)
	got, err := recoveryAdapter.GetTradesSince(100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty list")
	}
}

func TestRecoveryAdapterGetTradesSince_ErrorPropagation(t *testing.T) {
	expectedErr := errors.New("boom")
	client := adapter.NewClient(&fakeClient{err: expectedErr})

	recoveryAdapter := adapter.NewRecoveryAdapter(client)
	_, err := recoveryAdapter.GetTradesSince(100)
	if err == nil {
		t.Fatalf("expected error")
	}
}
