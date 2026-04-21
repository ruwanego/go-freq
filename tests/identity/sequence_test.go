package identity_test

import (
	"gofreq/internal/identity"
	"gofreq/internal/persistence"
	"testing"
)

func TestRecoverSequences(t *testing.T) {
	records := []persistence.OrderRecord{
		{EngineID: "GF-MACD-1000-0001"},
		{EngineID: "GF-MACD-1000-0003"},
		{EngineID: "GF-RSI-1000-0002"},
	}

	seqs := identity.RecoverSequences(records)

	if seqs["MACD"] != 3 {
		t.Fatalf("expected MACD seq 3")
	}
	if seqs["RSI"] != 2 {
		t.Fatalf("expected RSI seq 2")
	}
}
