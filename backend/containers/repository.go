package containers

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/database"
	"github.com/lib/pq"

	"github.com/jnie1/MTGViewer-V2/cards"
)

func GetAllocations() ([]ContainerAllocation, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT c.container_id, COALESCE(SUM(cd.amount), 0) AS used, c.capacity
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
		if err := row.Scan(&allocation.ContainerId, &allocation.Used, &allocation.Capacity); err != nil {
			return nil, err
		}

		allocations = append(allocations, allocation)
	}

	return allocations, nil
}

func GetContainers() ([]Container, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT c.container_id, c.container_name,  COALESCE(SUM(cd.amount), 0) AS used, c.capacity
		FROM containers c
		LEFT JOIN card_deposits cd ON c.container_id = cd.container_id
		GROUP BY c.container_id
		ORDER BY sort_order`)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	containers := []Container{}

	for row.Next() {
		container := Container{}
		if err := row.Scan(&container.ContainerId, &container.Name, &container.Used, &container.Capacity); err != nil {
			return nil, err
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func GetContainer(containerId int) (ContainerEntry, error) {
	db := database.Instance()

	row := db.QueryRow(`
		SELECT c.container_name, COALESCE(SUM(cd.amount), 0) AS used, c.capacity, c.deletion_mark
		FROM containers c
		LEFT JOIN card_deposits cd ON c.container_id = cd.container_id
		WHERE c.container_id = $1
		GROUP BY c.container_id
		LIMIT 1;`, containerId)

	container := ContainerEntry{}
	err := row.Scan(&container.Name, &container.Used, &container.Capacity, &container.IsDeleted)

	return container, err
}

func GetAmounts(containerId int) ([]cards.CardAmountPreview, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT scryfall_id, oracle_id, amount
		FROM card_deposits
		WHERE container_id = $1`, containerId)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	amounts := []cards.CardAmountPreview{}

	for row.Next() {
		amount := cards.CardAmountPreview{}
		if err := row.Scan(&amount.ScryfallId, &amount.OracleId, &amount.Amount); err != nil {
			return nil, err
		}

		amounts = append(amounts, amount)
	}

	return amounts, nil
}

func FindExcessAmounts(count int) ([]cards.CardAmountPreview, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT cd.scryfall_id, cd.oracle_id, SUM(cd.amount)
		FROM card_deposits cd 
		GROUP BY cd.scryfall_id, cd.oracle_id 
		HAVING SUM(cd.amount) > $1;`, count)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	deposits := []cards.CardAmountPreview{}

	for row.Next() {
		deposit := cards.CardAmountPreview{}
		if err := row.Scan(&deposit.ScryfallId, &deposit.OracleId, &deposit.Amount); err != nil {
			return nil, err
		}

		deposits = append(deposits, deposit)
	}

	return deposits, nil
}

func MatchCards(scryfallIds uuid.UUIDs) (uuid.UUIDs, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT DISTINCT cd.scryfall_id
		FROM card_deposits AS cd
		WHERE cd.scryfall_id = ANY($1);`, pq.Array(scryfallIds))

	if err != nil {
		return nil, err
	}

	defer row.Close()
	matches := uuid.UUIDs{}

	for row.Next() {
		var cardId uuid.UUID
		if err := row.Scan(&cardId); err != nil {
			return nil, err
		}

		matches = append(matches, cardId)
	}

	return matches, nil
}

func SearchDeposits(scryfallIds uuid.UUIDs) ([]CardDepositPreview, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT cd.container_id, c.container_name, cd.scryfall_id, cd.oracle_id, cd.amount
		FROM card_deposits AS cd
		JOIN containers AS c ON cd.container_id = c.container_id
		WHERE cd.scryfall_id = ANY($1)
		ORDER BY c.sort_order;`, pq.Array(scryfallIds))

	if err != nil {
		return nil, err
	}

	defer row.Close()

	deposits := []CardDepositPreview{}

	for row.Next() {
		deposit := CardDepositPreview{}
		if err := row.Scan(&deposit.ContainerId, &deposit.ContainerName, &deposit.ScryfallId, &deposit.OracleId, &deposit.Amount); err != nil {
			return nil, err
		}

		deposits = append(deposits, deposit)
	}

	return deposits, nil
}

func AddContainer(container ContainerEntry) error {
	db := database.Instance()

	_, err := db.Exec(`
		INSERT INTO containers (container_name, capacity, deletion_mark) 
		VALUES ($1, $2, FALSE)`, container.Name, container.Capacity)

	return err
}

func UpdateContainer(containerId int, container ContainerEntry) error {
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
