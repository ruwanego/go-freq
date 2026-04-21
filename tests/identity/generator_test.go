package identity_test

import (
	"gofreq/internal/identity"
	"testing"
)

func TestGeneratorSequential(t *testing.T) {
	gen := identity.NewGenerator("GF", map[string]int64{
		"MACD": 0,
	})

	id1 := gen.Next("MACD", 1000)
	id2 := gen.Next("MACD", 1000)

	if id1 == id2 {
		t.Fatalf("ids must be unique")
	}
}
