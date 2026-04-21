package identity

import (
	"fmt"
	"strconv"
	"strings"
)

type ParsedID struct {
	EnginePrefix string
	Strategy     string
	Timestamp    int64
	Sequence     int64
}

func Parse(id string) (ParsedID, error) {
	parts := strings.Split(id, "-")
	if len(parts) != 4 {
		return ParsedID{}, fmt.Errorf("invalid_id_format")
	}
	if parts[0] == "" || parts[1] == "" {
		return ParsedID{}, fmt.Errorf("invalid_id_format")
	}

	ts, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return ParsedID{}, fmt.Errorf("invalid_timestamp")
	}

	seq, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return ParsedID{}, fmt.Errorf("invalid_sequence")
	}

	return ParsedID{
		EnginePrefix: parts[0],
		Strategy:     parts[1],
		Timestamp:    ts,
		Sequence:     seq,
	}, nil
}
