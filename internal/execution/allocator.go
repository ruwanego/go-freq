package execution

import (
	"sort"

	"gofreq/pkg/actions"
	coreexec "gofreq/pkg/execution"
)

type Allocator interface {
	Allocate(input []actions.Action) ([]actions.Action, []coreexec.RejectedAction)
}

type DeterministicAllocator struct{}

func (a *DeterministicAllocator) Allocate(input []actions.Action) ([]actions.Action, []coreexec.RejectedAction) {
	allocated := append([]actions.Action(nil), input...)

	sort.Slice(allocated, func(i, j int) bool {
		if allocated[i].Tag != allocated[j].Tag {
			return allocated[i].Tag < allocated[j].Tag
		}
		return allocated[i].Pair < allocated[j].Pair
	})

	return allocated, []coreexec.RejectedAction{}
}
