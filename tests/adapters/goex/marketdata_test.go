package goex_test

import (
	"errors"
	"github.com/shopspring/decimal"
	"testing"
	"time"

	adapter "gofreq/internal/adapters/goex"
)

type fakeMarketClient struct {
	candles []adapter.RawCandle
	ch      chan adapter.RawCandle
	err     error
}

func (f *fakeMarketClient) GetOpenOrders() ([]adapter.RawOrder, error) {
	return nil, nil
}

func (f *fakeMarketClient) GetTradesSince(int64) ([]adapter.RawTrade, error) {
	return nil, nil
}

func (f *fakeMarketClient) SubmitOrder(adapter.GoexOrderRequest) (adapter.GoexOrderResponse, error) {
	return adapter.GoexOrderResponse{}, nil
}

func (f *fakeMarketClient) CancelOrder(string) error {
	return nil
}

func (f *fakeMarketClient) GetCandles(string, string, int) ([]adapter.RawCandle, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.candles, nil
}

func (f *fakeMarketClient) SubscribeCandles([]string, string) (<-chan adapter.RawCandle, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.ch, nil
}

func TestGetHistoricalCandles_SortedAscending(t *testing.T) {
	client := adapter.NewClient(&fakeMarketClient{
		candles: []adapter.RawCandle{
			{Pair: "BTC/USDT", Timestamp: 3, Open: decimal.NewFromInt(10), High: decimal.NewFromInt(12), Low: decimal.NewFromInt(9), Close: decimal.NewFromInt(11), Volume: decimal.NewFromInt(100), Closed: true},
			{Pair: "BTC/USDT", Timestamp: 1, Open: decimal.NewFromInt(8), High: decimal.NewFromInt(9), Low: decimal.NewFromInt(7), Close: decimal.RequireFromString("8.5"), Volume: decimal.NewFromInt(90), Closed: true},
			{Pair: "BTC/USDT", Timestamp: 2, Open: decimal.NewFromInt(9), High: decimal.NewFromInt(10), Low: decimal.NewFromInt(8), Close: decimal.RequireFromString("9.5"), Volume: decimal.NewFromInt(95), Closed: true},
		},
	})

	md := adapter.NewMarketDataAdapter(client)
	got, err := md.GetHistoricalCandles("BTC/USDT", "1m", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 candles, got %d", len(got))
	}
	if got[0].Timestamp != 1 || got[1].Timestamp != 2 || got[2].Timestamp != 3 {
		t.Fatalf("candles not sorted ascending")
	}
}

func TestGetHistoricalCandles_SkipsInvalid(t *testing.T) {
	client := adapter.NewClient(&fakeMarketClient{
		candles: []adapter.RawCandle{
			{Pair: "BTC/USDT", Timestamp: 1, Open: decimal.NewFromInt(8), High: decimal.NewFromInt(9), Low: decimal.NewFromInt(7), Close: decimal.RequireFromString("8.5"), Volume: decimal.NewFromInt(90), Closed: true},
			{Pair: "BTC/USDT", Timestamp: 2, Open: decimal.NewFromInt(0), High: decimal.NewFromInt(10), Low: decimal.NewFromInt(8), Close: decimal.RequireFromString("9.5"), Volume: decimal.NewFromInt(95), Closed: true},
		},
	})

	md := adapter.NewMarketDataAdapter(client)
	got, err := md.GetHistoricalCandles("BTC/USDT", "1m", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 valid candle, got %d", len(got))
	}
}

func TestGetHistoricalCandles_EmptyInput(t *testing.T) {
	client := adapter.NewClient(&fakeMarketClient{candles: []adapter.RawCandle{}})

	md := adapter.NewMarketDataAdapter(client)
	got, err := md.GetHistoricalCandles("BTC/USDT", "1m", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty result")
	}
}

func TestGetHistoricalCandles_ErrorPropagation(t *testing.T) {
	expectedErr := errors.New("candles failed")
	client := adapter.NewClient(&fakeMarketClient{err: expectedErr})

	md := adapter.NewMarketDataAdapter(client)
	_, err := md.GetHistoricalCandles("BTC/USDT", "1m", 1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestSubscribeClosedCandles_OnlyClosedEmitted(t *testing.T) {
	ch := make(chan adapter.RawCandle, 3)
	ch <- adapter.RawCandle{Pair: "BTC/USDT", Timestamp: 1, Open: decimal.NewFromInt(10), High: decimal.NewFromInt(12), Low: decimal.NewFromInt(9), Close: decimal.NewFromInt(11), Volume: decimal.NewFromInt(100), Closed: true}
	ch <- adapter.RawCandle{Pair: "BTC/USDT", Timestamp: 2, Open: decimal.NewFromInt(11), High: decimal.NewFromInt(13), Low: decimal.NewFromInt(10), Close: decimal.NewFromInt(12), Volume: decimal.NewFromInt(110), Closed: false}
	close(ch)

	client := adapter.NewClient(&fakeMarketClient{ch: ch})
	md := adapter.NewMarketDataAdapter(client)

	out, err := md.SubscribeClosedCandles([]string{"BTC/USDT"}, "1m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got []int64
	for c := range out {
		got = append(got, c.Timestamp)
	}

	if len(got) != 1 || got[0] != 1 {
		t.Fatalf("expected only closed valid candle")
	}
}

func TestSubscribeClosedCandles_InvalidSkipped(t *testing.T) {
	ch := make(chan adapter.RawCandle, 2)
	ch <- adapter.RawCandle{Pair: "BTC/USDT", Timestamp: 0, Open: decimal.NewFromInt(10), High: decimal.NewFromInt(12), Low: decimal.NewFromInt(9), Close: decimal.NewFromInt(11), Volume: decimal.NewFromInt(100), Closed: true}
	ch <- adapter.RawCandle{Pair: "BTC/USDT", Timestamp: 3, Open: decimal.NewFromInt(12), High: decimal.NewFromInt(14), Low: decimal.NewFromInt(11), Close: decimal.NewFromInt(13), Volume: decimal.NewFromInt(120), Closed: true}
	close(ch)

	client := adapter.NewClient(&fakeMarketClient{ch: ch})
	md := adapter.NewMarketDataAdapter(client)

	out, err := md.SubscribeClosedCandles([]string{"BTC/USDT"}, "1m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count := 0
	for range out {
		count++
	}
	if count != 1 {
		t.Fatalf("expected 1 valid candle, got %d", count)
	}
}

func TestSubscribeClosedCandles_ChannelClosesOnStop(t *testing.T) {
	ch := make(chan adapter.RawCandle)
	client := adapter.NewClient(&fakeMarketClient{ch: ch})
	md := adapter.NewMarketDataAdapter(client)

	out, err := md.SubscribeClosedCandles([]string{"BTC/USDT"}, "1m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := md.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}

	select {
	case _, ok := <-out:
		if ok {
			t.Fatalf("expected channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for channel close")
	}
}

func TestSubscribeClosedCandles_ErrorPropagation(t *testing.T) {
	expectedErr := errors.New("subscribe failed")
	client := adapter.NewClient(&fakeMarketClient{err: expectedErr})
	md := adapter.NewMarketDataAdapter(client)

	_, err := md.SubscribeClosedCandles([]string{"BTC/USDT"}, "1m")
	if err == nil {
		t.Fatalf("expected error")
	}
}
