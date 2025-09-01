package containers

import (
	"cmp"
	"errors"
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
	if len(additions) == 0 {
		return []ContainerChanges{}, nil
	}

	totalAdds := 0
	for _, add := range additions {
		totalAdds += add.Delta
	}

	fitAllAdds := []ContainerAllocation{}
	for _, allocation := range allocations {
		if allocation.Remaining() >= totalAdds {
			fitAllAdds = append(fitAllAdds, allocation)
		}
	}

	if len(fitAllAdds) > 0 {
		targetContainer := slices.MinFunc(fitAllAdds, compareRemainingAllocations)
		fitAllChanges := ContainerChanges{targetContainer.ContainerId, additions}
		return []ContainerChanges{fitAllChanges}, nil
	}

	totalRemaining := 0
	for _, allocation := range allocations {
		totalRemaining += allocation.Remaining()
	}

	if totalRemaining < totalAdds {
		return nil, errors.New("not enough space to fit the new additions")
	}

	additionAssignments := findBestFitAssignments(totalAdds, allocations)
	allChanges := assignContainerChanges(additions, allocations, additionAssignments)

	return allChanges, nil
}

func compareRemainingAllocations(a, b ContainerAllocation) int {
	return cmp.Compare(a.Remaining(), b.Remaining())
}

func findBestFitAssignments(totalAdds int, allocations []ContainerAllocation) map[int]int {
	assignments := map[int]int{}
	return assignments
}

func assignContainerChanges(additions []CardRequest, allocations []ContainerAllocation, assignments map[int]int) []ContainerChanges {
	allChanges := []ContainerChanges{}

	requestIndex := 0
	allocationIndex := 0

	currentRequest := additions[requestIndex]
	currentContainerId := allocations[allocationIndex].ContainerId

	remainingAssignments := assignments[currentContainerId]
	containerRequests := []CardRequest{}

	for requestIndex < len(additions) && allocationIndex < len(allocations) {
		if currentRequest.Delta < remainingAssignments {
			remainingAssignments -= currentRequest.Delta
			containerRequests = append(containerRequests, currentRequest)

			requestIndex += 1
			if requestIndex < len(additions) {
				currentRequest = additions[requestIndex]
			} else {
				currentRequest = CardRequest{}
			}
		} else if currentRequest.Delta > remainingAssignments {
			containerRequests = append(containerRequests, CardRequest{currentRequest.ScryfallId, remainingAssignments})
			allChanges = append(allChanges, ContainerChanges{currentContainerId, containerRequests})

			leftover := currentRequest.Delta - remainingAssignments
			currentRequest = CardRequest{currentRequest.ScryfallId, leftover}

			allocationIndex += 1
			if allocationIndex < len(allocations) {
				currentContainerId = allocations[allocationIndex].ContainerId
			} else {
				currentContainerId = 0
			}

			remainingAssignments = assignments[currentContainerId]
			containerRequests = []CardRequest{}
		} else {
			containerRequests = append(containerRequests, currentRequest)
			allChanges = append(allChanges, ContainerChanges{currentContainerId, containerRequests})

			requestIndex += 1
			if requestIndex < len(additions) {
				currentRequest = additions[requestIndex]
			} else {
				currentRequest = CardRequest{}
			}

			allocationIndex += 1
			if allocationIndex < len(allocations) {
				currentContainerId = allocations[allocationIndex].ContainerId
			} else {
				currentContainerId = 0
			}

			remainingAssignments = assignments[currentContainerId]
			containerRequests = []CardRequest{}
		}
	}

	return allChanges
}
