package execution

import (
	"errors"

	"gofreq/pkg/actions"
	coreexec "gofreq/pkg/execution"
)

type Pipeline struct {
	risk      RiskEngine
	allocator Allocator
}

func NewPipeline(r RiskEngine, a Allocator) *Pipeline {
	return &Pipeline{
		risk:      r,
		allocator: a,
	}
}

func (p *Pipeline) Process(input []actions.Action) (coreexec.ExecutionResult, error) {
	validated, err := coreexec.ValidateActions(input)
	if err != nil {
		return rejectEntireBatch(input, firstRejectedReason(validated.Rejected)), err
	}
	if len(validated.Rejected) > 0 {
		return rejectEntireBatch(input, firstRejectedReason(validated.Rejected)), errors.New("validation_failed")
	}

	safe, riskRejected := p.risk.Apply(validated.Accepted)
	if len(riskRejected) > 0 {
		return coreexec.ExecutionResult{
			Accepted: []actions.Action{},
			Rejected: rejectBatchFromReasons(input, riskRejected),
		}, nil
	}

	allocated, allocRejected := p.allocator.Allocate(safe)
	if len(allocRejected) > 0 {
		return coreexec.ExecutionResult{
			Accepted: []actions.Action{},
			Rejected: rejectBatchFromReasons(input, allocRejected),
		}, nil
	}

	return coreexec.ExecutionResult{
		Accepted: allocated,
		Rejected: []coreexec.RejectedAction{},
	}, nil
}

func rejectEntireBatch(input []actions.Action, reason string) coreexec.ExecutionResult {
	if reason == "" {
		reason = "validation_failed"
	}

	rejected := make([]coreexec.RejectedAction, 0, len(input))
	for _, a := range input {
		rejected = append(rejected, coreexec.RejectedAction{
			Action: a,
			Reason: reason,
		})
	}

	return coreexec.ExecutionResult{
		Accepted: []actions.Action{},
		Rejected: rejected,
	}
}

func rejectBatchFromReasons(input []actions.Action, stageRejected []coreexec.RejectedAction) []coreexec.RejectedAction {
	reason := firstRejectedReason(stageRejected)
	if reason == "" {
		reason = "stage_rejected"
	}

	rejected := make([]coreexec.RejectedAction, 0, len(input))
	for _, a := range input {
		rejected = append(rejected, coreexec.RejectedAction{
			Action: a,
			Reason: reason,
		})
	}

	return rejected
}

func firstRejectedReason(rejected []coreexec.RejectedAction) string {
	if len(rejected) == 0 {
		return ""
	}

	return rejected[0].Reason
}

var _ = errors.New
