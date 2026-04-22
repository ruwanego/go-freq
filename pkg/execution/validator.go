package execution

import (
	"errors"
	"github.com/shopspring/decimal"
	"gofreq/pkg/actions"
)

func ValidateActions(input []actions.Action) (ExecutionResult, error) {
	result := ExecutionResult{
		Accepted: []actions.Action{},
		Rejected: []RejectedAction{},
	}

	seenTags := make(map[string]bool)

	for _, a := range input {
		if a.Pair == "" {
			result.Rejected = append(result.Rejected, RejectedAction{a, "missing_pair"})
			continue
		}

		if a.Amount.LessThanOrEqual(decimal.Zero) {
			result.Rejected = append(result.Rejected, RejectedAction{a, "invalid_amount"})
			continue
		}

		if a.Tag == "" {
			result.Rejected = append(result.Rejected, RejectedAction{a, "missing_tag"})
			continue
		}

		if seenTags[a.Tag] {
			result.Rejected = append(result.Rejected, RejectedAction{a, "duplicate_tag"})
			continue
		}

		seenTags[a.Tag] = true
		result.Accepted = append(result.Accepted, a)
	}

	if len(result.Accepted) == 0 && len(result.Rejected) > 0 {
		return result, errors.New("all_actions_rejected")
	}

	return result, nil
}
