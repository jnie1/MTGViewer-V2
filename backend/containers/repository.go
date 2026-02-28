package containers

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/database"
	"github.com/lib/pq"
)

func GetAllocations() ([]ContainerAllocation, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT c.container_id, COALESCE(SUM(cd.amount), 0), c.capacity
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
		if err := row.Scan(&allocation.ContainerId, &allocation.Used, &allocation.MaxCapacity); err != nil {
			return nil, err
		}

		allocations = append(allocations, allocation)
	}

	return allocations, nil
}

func GetContainers() ([]ContainerPreview, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT container_id, container_name, capacity
		FROM containers
		ORDER BY container_name;`)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	containers := []ContainerPreview{}

	for row.Next() {
		container := ContainerPreview{}
		if err := row.Scan(&container.ContainerId, &container.Name, &container.Capacity); err != nil {
			return nil, err
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func GetContainer(containerId int) (Container, error) {
	db := database.Instance()

	row := db.QueryRow(`
		SELECT c2.container_name, c1.used, c2.capacity, c2.deletion_mark
		FROM (
			SELECT c.container_id, COALESCE(SUM(cd.amount), 0) as used
			FROM containers c
			LEFT JOIN card_deposits cd ON c.container_id = cd.container_id
			WHERE c.container_id = $1
			GROUP BY c.container_id
		) AS c1
		JOIN containers AS c2 ON c1.container_id = c2.container_id;`, containerId)

	container := Container{}
	err := row.Scan(&container.Name, &container.Used, &container.Capacity, &container.IsDeleted)

	return container, err
}

func GetDeposits(containerId int) ([]CardDepositPreview, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT scryfall_id, amount
		FROM card_deposits
		WHERE container_id = $1`, containerId)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	deposits := []CardDepositPreview{}

	for row.Next() {
		deposit := CardDepositPreview{}
		if err := row.Scan(&deposit.ScryfallId, &deposit.Amount); err != nil {
			return nil, err
		}

		deposits = append(deposits, deposit)
	}

	return deposits, nil
}

func SearchCards(scryfallIds uuid.UUIDs) ([]CardDeposit, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT cd.container_id, c.container_name, cd.scryfall_id, cd.amount
		FROM card_deposits AS cd
		JOIN containers AS c ON cd.container_id = c.container_id
		WHERE cd.scryfall_id = ANY($1);`, pq.Array(scryfallIds))

	if err != nil {
		return nil, err
	}

	defer row.Close()

	deposits := []CardDeposit{}

	for row.Next() {
		deposit := CardDeposit{}
		if err := row.Scan(&deposit.ContainerId, &deposit.ContainerName, &deposit.ScryfallId, &deposit.Amount); err != nil {
			return nil, err
		}

		deposits = append(deposits, deposit)
	}

	return deposits, nil
}

func AddContainer(container Container) error {
	db := database.Instance()

	_, err := db.Exec(`
		INSERT INTO containers (container_name, capacity, deletion_mark) 
		VALUES ($1, $2, FALSE)`, container.Name, container.Capacity)

	return err
}

func UpdateContainer(containerId int, container Container) error {
	db := database.Instance()

	_, err := db.Exec(`
		UPDATE containers
		SET container_name = $2, capacity = $3, deletion_mark = $4
		WHERE container_id = $1;`, containerId, container.Name, container.Capacity, container.IsDeleted)

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
		DELETE FROM containers
		WHERE container_id = $1`, containerId)

	return err
}
