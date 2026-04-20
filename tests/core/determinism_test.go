package core

import (
    "gofreq/pkg/actions"
    "gofreq/pkg/execution"
    "reflect"
    "testing"
)

func TestDeterministicValidation(t *testing.T) {
    acts := []actions.Action{
        {Pair: "BTC/USDT", Amount: 1, Tag: "a"},
        {Pair: "ETH/USDT", Amount: 2, Tag: "b"},
    }

    r1, _ := execution.ValidateActions(acts)
    r2, _ := execution.ValidateActions(acts)

    if !reflect.DeepEqual(r1, r2) {
        t.Fatalf("validation not deterministic")
    }
}
