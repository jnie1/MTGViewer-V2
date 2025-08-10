package containers

import (
	"github.com/jnie1/MTGViewer-V2/database"
)

func GetContainer(containerId int) (Container, error) {
	db := database.Instance()

	row := db.QueryRow(`
		SELECT container, capacity
		FROM container
		WHERE container_id = $1;`, containerId)

	container := Container{}
	err := row.Scan(&container.Name, &container.Capacity, &container.MarkForDeletion)

	return container, err
}

func GetDeposits(containerId int) ([]CardDeposit, error) {
	db := database.Instance()

	row, err := db.Query(`
		SELECT scryfall_id, amount
		FROM card_deposits
		WHERE container_id = $1`, containerId)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	deposits := []CardDeposit{}

	for row.Next() {
		deposit := CardDeposit{}
		if err := row.Scan(&deposit.ScryfallId, &deposit.Amount); err != nil {
			return nil, err
		}

		deposits = append(deposits, deposit)
	}

	return deposits, nil
}

func AddContainer(cont Container) error {
	db := database.Instance()

	_, err := db.Exec(`
		INSERT INTO container (container_name, capacity, deletion_mark) 
		VALUES ($1, $2, FALSE)`, cont.Name, cont.Capacity)

	return err
}

func UpdateContainer(containerId int, container Container) error {
	db := database.Instance()

	_, err := db.Exec(`
		UPDATE container
		SET container_name = $2, capacity = $3, mark_for_deletion = $4
		WHERE container_id = $1;`, containerId, container.Name, container.Capacity, container.MarkForDeletion)

	return err
}

func DeleteContainer(containerId int) error {

	db := database.Instance()

	_, err := db.Exec(`
		DELETE FROM container
		WHERE container_id = $1`, containerId)

	return err
}
