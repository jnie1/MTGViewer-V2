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

func FetchUpdateLogs() ([]UpdateLogs, error) {
	db := database.Instance()
	row, err := db.Query(`
		SELECT group_id, time
		FROM transactions
		GROUP BY group_id, time;`)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	listOfLogs := []UpdateLogs{}

	for row.Next() {
		log := UpdateLogs{}

		if err := row.Scan(&log.GroupId, &log.Time); err != nil {
			return nil, err
		}

		listOfLogs = append(listOfLogs, log)
	}

	return listOfLogs, nil
}

func FetchLogs(groupId uuid.UUID) ([]TransactionLogs, error) {
	db := database.Instance()
	row, err := db.Query(`
		SELECT transaction_id, group_id, from_container_id, fc.container_name, to_container_id, tc.container_name, scryfall_id, amount
		FROM transactions AS t
		LEFT JOIN containers AS fc ON from_container_id = fc.container_id
		LEFT JOIN containers AS tc ON to_container_id = tc.container_id
		WHERE t.group_id = $1;`, groupId)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	listOfLogs := []TransactionLogs{}

	for row.Next() {
		log := TransactionLogs{}

		var fromMaybeBoxId sql.Null[int]
		var fromMaybeBoxName sql.NullString

		var toMaybeBoxId sql.Null[int]
		var toMaybeBoxName sql.NullString

		if err := row.Scan(&log.TransactionId, &log.GroupId, &fromMaybeBoxId, &fromMaybeBoxName, &toMaybeBoxId, &toMaybeBoxName, &log.ScryfallId, &log.Quantity); err != nil {
			return nil, err
		}

		if fromMaybeBoxId.Valid && fromMaybeBoxName.Valid {
			log.FromContainer = &TransactionContainer{fromMaybeBoxId.V, fromMaybeBoxName.String}
		}

		if toMaybeBoxId.Valid && toMaybeBoxName.Valid {
			log.ToContainer = &TransactionContainer{toMaybeBoxId.V, toMaybeBoxName.String}
		}

		listOfLogs = append(listOfLogs, log)
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
				valueRow := fmt.Sprintf("('%s'::uuid, NULL, %d, '%s'::uuid, %d, '%s')", groupId, change.ContainerId, request.ScryfallId, request.Delta, now.Format(time.RFC3339))
				valueStatements = append(valueStatements, valueRow)

			case request.Delta < 0:
				valueRow := fmt.Sprintf("('%s'::uuid, %d, NULL, '%s'::uuid, %d, '%s')", groupId, change.ContainerId, request.ScryfallId, -request.Delta, now.Format(time.RFC3339))
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
