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

	bestFit := findBestFitAssignments(totalAdds, allocations)
	additionAssignments := slices.Collect(bestFit)

	if len(additionAssignments) == 0 {
		additionAssignments = append(additionAssignments, allocations...)
	}

	allChanges := assignContainerChanges(additions, additionAssignments)

	return allChanges, nil
}

func compareRemainingAllocations(a, b ContainerAllocation) int {
	return cmp.Compare(a.Remaining(), b.Remaining())
}

func findBestFitAssignments(totalAdds int, allocations []ContainerAllocation) iter.Seq[ContainerAllocation] {
	return func(yield func(ContainerAllocation) bool) {
		if len(allocations) <= 1 {
			return
		}

		leftCombinations := getAllocationCombinations(0, 0, totalAdds, nil, allocations[:len(allocations)/2])
		rightCombinations := getAllocationCombinations(0, 0, totalAdds, nil, allocations[len(allocations)/2:])

		slices.SortFunc(leftCombinations, compareRemainingCombinations)
		combos := [2]allocationCombination{}
		minRemainingSpace := totalAdds

		for _, firstCombo := range rightCombinations {
			remaining := max(totalAdds-firstCombo.TotalRemaining, 0)
			secondComboIndex, found := slices.BinarySearchFunc(leftCombinations, remaining, checkRemainingCombinations)

			if found {
				combos[0] = firstCombo
				combos[1] = leftCombinations[secondComboIndex]
				break
			}

			if secondComboIndex == len(leftCombinations) {
				// too small
				continue
			}

			secondCombo := leftCombinations[secondComboIndex]
			remainingSpace := firstCombo.TotalRemaining + secondCombo.TotalRemaining - totalAdds

			if remainingSpace < minRemainingSpace {
				combos[0] = firstCombo
				combos[1] = secondCombo
				minRemainingSpace = remainingSpace
			}
		}

		if combos[0].TotalRemaining == 0 && combos[1].TotalRemaining == 0 {
			return
		}

		allocationMap := map[int]ContainerAllocation{}
		for _, allocation := range allocations {
			allocationMap[allocation.ContainerId] = allocation
		}

		for _, combo := range combos {
			for containerId := range combo.getContainerIds() {
				if alloc, ok := allocationMap[containerId]; ok {
					if !yield(alloc) {
						return
					}
				}
			}
		}
	}
}

type allocationGroup struct {
	ContainerId int
	Next        *allocationGroup
}

type allocationCombination struct {
	TotalRemaining int
	Items          *allocationGroup
}

func (combo allocationCombination) getContainerIds() iter.Seq[int] {
	return func(yield func(int) bool) {
		for group := combo.Items; group != nil; group = group.Next {
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

	withAllocation := allocationGroup{currentAllocation.ContainerId, items}
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

			containerRequests = []CardRequest{}
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
			containerRequests = []CardRequest{}

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
