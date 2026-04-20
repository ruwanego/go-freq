package pipeline

import (
	"testing"

	"gofreq/internal/execution"
	"gofreq/pkg/actions"
)

func TestPipelineRejectsEntireBatchOnValidationFailure(t *testing.T) {
	risk := &execution.BasicRisk{MaxPerTrade: 5}
	alloc := &execution.DeterministicAllocator{}

	p := execution.NewPipeline(risk, alloc)

	acts := []actions.Action{
		{Pair: "BTC/USDT", Amount: 1, Tag: "a"},
		{Pair: "", Amount: 1, Tag: "b"},
	}

	res, err := p.Process(acts)

	if err == nil {
		t.Fatalf("expected validation error")
	}

	if len(res.Accepted) != 0 {
		t.Fatalf("no actions should be accepted")
	}

	if len(res.Rejected) != 2 {
		t.Fatalf("entire batch must be rejected")
	}
}
