package identity_test

import (
	"gofreq/internal/identity"
	"testing"
)

func TestParseValidID(t *testing.T) {
	id := "GF-MACD-1000-0001"

	parsed, err := identity.Parse(id)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if parsed.Strategy != "MACD" {
		t.Fatalf("wrong strategy")
	}
	if parsed.Sequence != 1 {
		t.Fatalf("wrong sequence")
	}
}
