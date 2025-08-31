package containers

import (
	"cmp"
	"slices"
)

// TODO: get container removals

func GetContainerChanges(requests []CardRequest, allocations []ContainerAllocation) ([]ContainerChanges, error) {
	additions := []CardRequest{}
	for _, request := range requests {
		if request.Delta > 0 {
			additions = append(additions, request)
		}
	}

	addChanges, err := getContainerAdditions(additions, allocations)
	if err != nil {
		return nil, err
	}

	return addChanges, nil
}

func getContainerAdditions(additions []CardRequest, allocations []ContainerAllocation) ([]ContainerChanges, error) {
	totalAdds := 0
	for _, add := range additions {
		totalAdds += add.Delta
	}

	fitAllAdds := []ContainerAllocation{}
	for _, allocation := range allocations {
		remaining := allocation.MaxCapacity - allocation.Used
		if remaining >= totalAdds {
			fitAllAdds = append(fitAllAdds, allocation)
		}
	}

	if len(fitAllAdds) > 0 {
		targetContainer := slices.MinFunc(fitAllAdds, compareRemainingAllocations)
		fitAllChanges := ContainerChanges{ContainerId: targetContainer.ContainerId, Requests: additions}

		return []ContainerChanges{fitAllChanges}, nil
	}

	return nil, nil
}

func compareRemainingAllocations(a, b ContainerAllocation) int {
	return cmp.Compare(a.MaxCapacity-a.Used, b.MaxCapacity-b.Used)
}
