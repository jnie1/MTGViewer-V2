package transactions

import (
	"database/sql"
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
		SELECT transaction_id, group_id, from_container_id, fc.container_name, to_container_id, tc.container_name, scryfall_id, amount
		FROM transactions as t
		LEFT JOIN containers as fc ON from_container_id = fc.container_id
		LEFT JOIN containers as tc ON to_container_id = tc.container_id;`)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	listOfLogs := []TransactionLogs{}

	for row.Next() {
		logs := TransactionLogs{}

		var fromMaybeBoxId sql.Null[int]
		var fromMaybeBoxName sql.NullString

		var toMaybeBoxId sql.Null[int]
		var toMaybeBoxName sql.NullString

		if err := row.Scan(&logs.TransactionId, &logs.GroupId, &fromMaybeBoxId, &fromMaybeBoxName, &toMaybeBoxId, &toMaybeBoxName, &logs.ScryfallId, &logs.Quantity); err != nil {
			return nil, err
		}

		if fromMaybeBoxId.Valid && fromMaybeBoxName.Valid {
			logs.FromContainer = &TransactionContainer{fromMaybeBoxId.V, fromMaybeBoxName.String}
		}

		if toMaybeBoxId.Valid && toMaybeBoxName.Valid {
			logs.FromContainer = &TransactionContainer{toMaybeBoxId.V, fromMaybeBoxName.String}
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
	now := time.Now().UTC()

	valueStatements := []string{}

	for _, change := range changes {
		for _, request := range change.Requests {

			switch {
			case request.Delta > 0:
				valueRow := fmt.Sprintf("('%s'::uuid, %d, NULL, '%s'::uuid, %d, '%s')", groupId, change.ContainerId, request.ScryfallId, -request.Delta, now.Format(time.RFC3339))
				valueStatements = append(valueStatements, valueRow)

			case request.Delta < 0:
				valueRow := fmt.Sprintf("('%s'::uuid, NULL, %d, '%s'::uuid, %d, '%s')", groupId, change.ContainerId, request.ScryfallId, request.Delta, now.Format(time.RFC3339))
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
