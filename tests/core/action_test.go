package core

import (
    "gofreq/pkg/actions"
    "testing"
)

func TestBuilderCreatesActions(t *testing.T) {
    b := actions.NewBuilder()

    b.BuyLimit("BTC/USDT", 60000, 1.0, "test1")

    acts := b.Build()

    if len(acts) != 1 {
        t.Fatalf("expected 1 action, got %d", len(acts))
    }

    if acts[0].Tag != "test1" {
        t.Fatalf("tag mismatch")
    }
}
