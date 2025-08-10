package containers

import (
	"fmt"
	"strings"

	"github.com/jnie1/MTGViewer-V2/database"
)

func GetAllocations() ([]ContainerAllocation, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT c.container_id, SUM(cd.amount), c.capacity
		FROM containers c
		LEFT JOIN card_deposits cd ON c.container_id = cd.container_id
		GROUP BY c.container_id`)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	allocations := []ContainerAllocation{}

	for row.Next() {
		allocation := ContainerAllocation{}
		if err := row.Scan(&allocation.ContainerId, &allocation.Used, &allocation.MaxCapcity); err != nil {
			return nil, err
		}

		allocations = append(allocations, allocation)
	}

	return allocations, nil
}

func GetContainer(containerId int) (Container, error) {
	db := database.Instance()

	row := db.QueryRow(`
		SELECT container, capacity, deletion_mark
		FROM container
		WHERE container_id = $1;`, containerId)

	container := Container{}
	err := row.Scan(&container.Name, &container.Capacity, &container.MarkForDeletion)

	return container, err
}

func GetDeposits(containerId int) ([]CardDeposit, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT container_id, scryfall_id, amount
		FROM card_deposits
		WHERE container_id = $1`, containerId)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	deposits := []CardDeposit{}

	for row.Next() {
		deposit := CardDeposit{}
		if err := row.Scan(&deposit.ContainerId, &deposit.ScryfallId, &deposit.Amount); err != nil {
			return nil, err
		}

		deposits = append(deposits, deposit)
	}

	return deposits, nil
}

func AddContainer(container Container) error {
	db := database.Instance()

	_, err := db.Exec(`
		INSERT INTO container (container_name, capacity, deletion_mark) 
		VALUES ($1, $2, FALSE)`, container.Name, container.Capacity)

	return err
}

func UpdateContainer(containerId int, container Container) error {
	db := database.Instance()

	_, err := db.Exec(`
		UPDATE container
		SET container_name = $2, capacity = $3, deletion_mark = $4
		WHERE container_id = $1;`, containerId, container.Name, container.Capacity, container.MarkForDeletion)

	return err
}

func UpdateDeposits(changes []ContainerChanges) error {
	db := database.Instance()

	valueStatements := []string{}

	for _, change := range changes {
		for _, request := range change.Requests {
			valueRow := fmt.Sprintf("(%d, '%s'::uuid, %d)", change.ContainerId, request.ScryfallId, request.Delta)
			valueStatements = append(valueStatements, valueRow)
		}
	}

	allValues := strings.Join(valueStatements, ", ")

	_, err := db.Exec(`
		MERGE INTO card_deposits AS cd
		USING (VALUES ` + allValues + `) AS ds (container_id, scryfall_id, delta)
		ON cd.container_id = ds.container_id AND cd.scryfall_id = ds.scryfall_id
		WHEN NOT MATCHED THEN
			INSERT (container_id, scryfall_id, amount) VALUES (ds.container_id, ds.scryfall_id, ds.delta)
		WHEN MATCHED AND cd.amount + ds.delta > 0 THEN
			UPDATE SET amount = cd.amount + ds.delta
		WHEN MATCHED THEN
			DELETE`)

	return err
}

func DeleteContainer(containerId int) error {

	db := database.Instance()

	_, err := db.Exec(`
		DELETE FROM container
		WHERE container_id = $1`, containerId)

	return err
}
