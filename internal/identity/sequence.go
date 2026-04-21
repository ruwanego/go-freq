package identity

import "gofreq/internal/persistence"

func RecoverSequences(records []persistence.OrderRecord) map[string]int64 {
	out := map[string]int64{}

	for _, r := range records {
		parsed, err := Parse(r.EngineID)
		if err != nil {
			continue
		}

		if parsed.Sequence > out[parsed.Strategy] {
			out[parsed.Strategy] = parsed.Sequence
		}
	}

	return out
}
