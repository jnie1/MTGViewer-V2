package transactions

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/containers"
	"github.com/jnie1/MTGViewer-V2/database"
)

func FetchLogs() ([]TransactionLogs, error) {
	db := database.Instance()
	row, err := db.Query(`
		SELECT transaction_id, group_id, from_container_id, to_container_id, scryfall_id, amount 	
		FROM transactions`)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	listOfLogs := []TransactionLogs{}

	for row.Next() {
		logs := TransactionLogs{}
		if err := row.Scan(&logs.TransactionId, &logs.GroupId, &logs.FromContainer, &logs.ToContainer, &logs.ScryfallId, &logs.Quantity); err != nil {
			return nil, err
		}

		listOfLogs = append(listOfLogs, logs)
	}

	return listOfLogs, nil
}

func LogCollectionChanges(changes []containers.ContainerChanges) error {
	groupId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	db := database.Instance()
	now := time.Now()

	valueStatements := []string{}

	for _, change := range changes {
		for _, request := range change.Requests {

			switch {
			case request.Delta > 0:
				valueRow := fmt.Sprintf("('%s'::uuid, %d, NULL, '%s'::uuid, %d, '%v')", groupId, change.ContainerId, request.ScryfallId, -request.Delta, now)
				valueStatements = append(valueStatements, valueRow)

			case request.Delta < 0:
				valueRow := fmt.Sprintf("('%s'::uuid, NULL, %d, '%s'::uuid, %d, '%v')", groupId, change.ContainerId, request.ScryfallId, request.Delta, now)
				valueStatements = append(valueStatements, valueRow)
			}
		}
	}
	allValues := strings.Join(valueStatements, ", ")

	_, err = db.Exec(`
		INSERT INTO transactions (group_id, from_container_id, to_container_id, scryfall_id, amount, time)
		VALUES ` + allValues + `;`)

	return err
}
