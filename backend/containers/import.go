package containers

import (
	"cmp"
	"errors"
	"iter"
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
	allChanges := assignContainerChanges(additions, additionAssignments)

	return allChanges, nil
}

func compareRemainingAllocations(a, b ContainerAllocation) int {
	return cmp.Compare(a.Remaining(), b.Remaining())
}

func findBestFitAssignments(totalAdds int, allocations []ContainerAllocation) []ContainerAllocation {
	var bestCombo *allocationPair
	minSize := len(allocations)

	leftCombinations := getAllocationCombinations(0, 0, totalAdds, nil, allocations[:len(allocations)/2])
	rightCombinations := getAllocationCombinations(0, 0, totalAdds, nil, allocations[len(allocations)/2:])
	slices.SortFunc(rightCombinations, compareRemainingCombinations)

	for _, firstCombo := range leftCombinations {
		remaining := max(totalAdds-firstCombo.TotalRemaining, 0)
		secondComboIndex, _ := slices.BinarySearchFunc(rightCombinations, remaining, checkRemainingCombinations)

		if secondComboIndex == len(rightCombinations) {
			// too small
			continue
		}

		possibleCombo := allocationPair{firstCombo, rightCombinations[secondComboIndex]}

		if possibleCombo.Size() < minSize {
			bestCombo = &possibleCombo
			minSize = possibleCombo.Size()
		}
	}

	if bestCombo == nil {
		return nil
	}

	chosenContainerIds := map[int]bool{}
	for containerId := range bestCombo.ContainerIds() {
		chosenContainerIds[containerId] = true
	}

	assignments := []ContainerAllocation{}
	for _, allocation := range allocations {
		if chosenContainerIds[allocation.ContainerId] {
			assignments = append(assignments, allocation)
		}
	}

	return assignments
}

type allocationGroup struct {
	ContainerId int
	Size        int
	Next        *allocationGroup
}

type allocationCombination struct {
	TotalRemaining int
	Items          *allocationGroup
}

type allocationPair struct {
	First  allocationCombination
	Second allocationCombination
}

func (pair allocationPair) Size() int {
	size := 0

	firstItems := pair.First.Items
	if firstItems != nil {
		size += firstItems.Size
	}

	secondItems := pair.Second.Items
	if secondItems != nil {
		size += secondItems.Size
	}

	return size
}

func (pair allocationPair) ContainerIds() iter.Seq[int] {
	return func(yield func(int) bool) {
		for group := pair.First.Items; group != nil; group = group.Next {
			if !yield(group.ContainerId) {
				return
			}
		}
		for group := pair.Second.Items; group != nil; group = group.Next {
			if !yield(group.ContainerId) {
				return
			}
		}
	}
}

func compareRemainingCombinations(a, b allocationCombination) int {
	return cmp.Compare(a.TotalRemaining, b.TotalRemaining)
}

func checkRemainingCombinations(a allocationCombination, target int) int {
	return cmp.Compare(a.TotalRemaining, target)
}

func getAllocationCombinations(i, totalRemaining, target int, items *allocationGroup, allocations []ContainerAllocation) []allocationCombination {
	if i == len(allocations) || totalRemaining >= target {
		return []allocationCombination{{totalRemaining, items}}
	}

	excludedCombos := getAllocationCombinations(i+1, totalRemaining, target, items, allocations)
	currentAllocation := allocations[i]
	remaining := currentAllocation.Remaining()

	if remaining == 0 {
		return excludedCombos
	}

	var groupSize int
	if items != nil {
		groupSize = items.Size + 1
	} else {
		groupSize = 1
	}

	withAllocation := allocationGroup{currentAllocation.ContainerId, groupSize, items}
	includedCombos := getAllocationCombinations(i+1, totalRemaining+remaining, target, &withAllocation, allocations)

	return slices.Concat(excludedCombos, includedCombos)
}

func assignContainerChanges(additions []CardRequest, assignments []ContainerAllocation) []ContainerChanges {
	allChanges := []ContainerChanges{}

	requestIndex := 0
	assignmentIndex := 0

	containerRequests := []CardRequest{}
	currentRequest := additions[requestIndex]
	currentAssignment := assignments[assignmentIndex]

	for requestIndex < len(additions) && assignmentIndex < len(assignments) {
		assignmentRemaining := currentAssignment.Remaining()

		if currentRequest.Delta < assignmentRemaining {
			containerRequests = append(containerRequests, currentRequest)
			currentAssignmentUsed := currentRequest.Delta + currentAssignment.Used

			requestIndex += 1
			if requestIndex < len(additions) {
				currentRequest = additions[requestIndex]
			} else {
				currentRequest = CardRequest{}
			}

			currentAssignment = ContainerAllocation{
				currentAssignment.ContainerId,
				currentAssignmentUsed,
				currentAssignment.MaxCapacity,
			}
		} else if currentRequest.Delta > assignmentRemaining {
			remainingRequest := CardRequest{currentRequest.ScryfallId, assignmentRemaining}
			fullRequests := append(containerRequests, remainingRequest)

			newChanges := ContainerChanges{currentAssignment.ContainerId, fullRequests}
			allChanges = append(allChanges, newChanges)

			containerRequests = nil
			leftover := currentRequest.Delta - assignmentRemaining
			currentRequest = CardRequest{currentRequest.ScryfallId, leftover}

			assignmentIndex += 1
			if assignmentIndex < len(assignments) {
				currentAssignment = assignments[assignmentIndex]
			} else {
				currentAssignment = ContainerAllocation{}
			}
		} else {
			fullRequests := append(containerRequests, currentRequest)
			newChanges := ContainerChanges{currentAssignment.ContainerId, fullRequests}

			allChanges = append(allChanges, newChanges)
			containerRequests = nil

			requestIndex += 1
			if requestIndex < len(additions) {
				currentRequest = additions[requestIndex]
			} else {
				currentRequest = CardRequest{}
			}

			assignmentIndex += 1
			if assignmentIndex < len(assignments) {
				currentAssignment = assignments[assignmentIndex]
			} else {
				currentAssignment = ContainerAllocation{}
			}
		}
	}

	if len(containerRequests) > 0 && currentAssignment.ContainerId != 0 {
		newChanges := ContainerChanges{currentAssignment.ContainerId, containerRequests}
		allChanges = append(allChanges, newChanges)
	}

	return allChanges
}
