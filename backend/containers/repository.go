package containers

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/database"
	"github.com/lib/pq"

	"github.com/jnie1/MTGViewer-V2/cards"
)

func GetAllocations(ctx context.Context) ([]ContainerAllocation, error) {
	db := database.Instance()

	row, err := db.QueryContext(ctx, `
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

func GetContainers(ctx context.Context) ([]Container, error) {
	db := database.Instance()

	row, err := db.QueryContext(ctx, `
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

func GetContainer(ctx context.Context, containerId int) (ContainerEntry, error) {
	db := database.Instance()

	row := db.QueryRowContext(ctx, `
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

func GetAmounts(ctx context.Context, containerId int) ([]cards.CardAmountPreview, error) {
	db := database.Instance()

	row, err := db.QueryContext(ctx, `
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

func FindExcessDeposits(ctx context.Context, count int) ([]CardDepositPreview, error) {
	db := database.Instance()

	row, err := db.QueryContext(ctx, `
		SELECT cd.container_id, c.container_name, cd.scryfall_id, cd.oracle_id, cd.amount
		FROM card_deposits cd
		JOIN containers AS c ON cd.container_id = c.container_id
		WHERE cd.oracle_id IN (
			SELECT DISTINCT cd2.oracle_id
			FROM card_deposits cd2
			GROUP BY cd2.oracle_id
			HAVING SUM(cd2.amount) > $1);`, count)

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

func SearchDeposits(ctx context.Context, scryfallIds uuid.UUIDs) ([]CardDepositPreview, error) {
	db := database.Instance()

	row, err := db.QueryContext(ctx, `
		SELECT cd.container_id, c.container_name, cd.scryfall_id, cd.oracle_id, cd.amount
		FROM card_deposits AS cd
		JOIN containers AS c ON cd.container_id = c.container_id
		WHERE cd.scryfall_id = ANY($1);`, pq.Array(scryfallIds))

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

func SearchDepositsByOracleId(ctx context.Context, oracleIds uuid.UUIDs) ([]CardDepositPreview, error) {
	db := database.Instance()

	row, err := db.QueryContext(ctx, `
		SELECT cd.container_id, c.container_name, cd.scryfall_id, cd.oracle_id, cd.amount
		FROM card_deposits AS cd
		JOIN containers AS c ON cd.container_id = c.container_id
		WHERE cd.oracle_id = ANY($1);`, pq.Array(oracleIds))

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

func AddContainer(ctx context.Context, container ContainerEntry) error {
	db := database.Instance()

	_, err := db.ExecContext(ctx, `
		INSERT INTO containers (container_name, capacity, deletion_mark) 
		VALUES ($1, $2, FALSE);`, container.Name, container.Capacity)

	return err
}

func UpdateContainer(ctx context.Context, containerId int, container ContainerEntry) error {
	db := database.Instance()

	_, err := db.ExecContext(ctx, `
		UPDATE containers
		SET container_name = $2, capacity = $3, deletion_mark = $4
		WHERE container_id = $1;`, containerId, container.Name, container.Capacity, container.IsDeleted)

	return err
}

func UpdateDeposits(ctx context.Context, changes []ContainerChanges) error {
	db := database.Instance()

	vals := make([]string, len(changes))
	for _, change := range changes {
		for i, request := range change.Requests {
			vals[i] = fmt.Sprintf("(%d, '%s'::uuid, '%s'::uuid, %d)", change.ContainerId, request.ScryfallId, request.OracleId, request.Delta)
		}
	}

	values := strings.Join(vals, ", ")

	_, err := db.ExecContext(ctx, `
		MERGE INTO card_deposits AS cd
		USING (VALUES `+values+`) AS ds (container_id, scryfall_id, oracle_id, delta)
		ON cd.container_id = ds.container_id AND cd.scryfall_id = ds.scryfall_id
		WHEN NOT MATCHED THEN
			INSERT (container_id, scryfall_id, oracle_id, amount) VALUES (ds.container_id, ds.scryfall_id, ds.oracle_id, ds.delta)
		WHEN MATCHED AND cd.amount + ds.delta > 0 THEN
			UPDATE SET amount = cd.amount + ds.delta
		WHEN MATCHED THEN
			DELETE;`)

	return err
}

func DeleteContainer(ctx context.Context, containerId int) error {
	db := database.Instance()

	_, err := db.ExecContext(ctx, `
		DELETE FROM containers
		WHERE container_id = $1;`, containerId)

	return err
}

func FindMissingOracleIds(ctx context.Context) (uuid.UUIDs, error) {
	db := database.Instance()

	row, err := db.QueryContext(ctx, `
		SELECT DISTINCT cd.scryfall_id
		FROM card_deposits cd
		WHERE cd.oracle_id = $1;`, uuid.Nil)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	ids := uuid.UUIDs{}

	for row.Next() {
		id := uuid.UUID{}
		if err := row.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func UpdateOracleIds(ctx context.Context, oracleIds []cards.ScryfallOracleObj) error {
	db := database.Instance()

	vals := make([]string, len(oracleIds))
	for i, id := range oracleIds {
		vals[i] = fmt.Sprintf("('%s'::uuid,'%s'::uuid)", id.ScryfallId, id.OracleId)
	}
	values := strings.Join(vals, ", ")

	_, err := db.ExecContext(ctx, `
		MERGE INTO card_deposits AS cd
		USING (VALUES `+values+`) AS os (scryfall_id, oracle_id)
		ON cd.scryfall_id = os.scryfall_id
		WHEN MATCHED THEN
			UPDATE SET oracle_id = os.oracle_id;`)

	return err
}
