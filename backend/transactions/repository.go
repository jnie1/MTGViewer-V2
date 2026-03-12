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

func FetchLogRange(group1, group2 uuid.UUID) (LogRange, error) {
	db := database.Instance()

	row := db.QueryRow(`
		SELECT MIN(time) AS start, MAX(time) AS end
		FROM transactions
		WHERE group_id = $1 OR group_id = $2;`, group1, group2)

	logRange := LogRange{}
	err := row.Scan(&logRange.start, &logRange.end)

	return logRange, err
}

func FetchUpdateLogs() ([]UpdateLogs, error) {
	db := database.Instance()
	row, err := db.Query(`
		SELECT group_id, time, SUM(amount) AS amount
		FROM transactions
		GROUP BY group_id, time
		ORDER BY time DESC;`)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	logs := []UpdateLogs{}

	for row.Next() {
		log := UpdateLogs{}

		if err := row.Scan(&log.GroupId, &log.Time, &log.Amount); err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func FetchLogs(groupId uuid.UUID) ([]TransactionLogs, error) {
	db := database.Instance()
	row, err := db.Query(`
		SELECT fc.container_id, fc.container_name, tc.container_id, tc.container_name, scryfall_id, amount
		FROM transactions
		LEFT JOIN containers AS fc ON from_container_id = fc.container_id
		LEFT JOIN containers AS tc ON to_container_id = tc.container_id
		WHERE group_id = $1;`, groupId)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	return fetchLogsFromQuery(row)
}

func FetchLogsFromRange(logRange LogRange) ([]TransactionLogs, error) {
	db := database.Instance()
	row, err := db.Query(`
		SELECT fc.container_id, fc.container_name, tc.container_id, tc.container_name, scryfall_id, amount
		FROM transactions
		LEFT JOIN containers AS fc ON from_container_id = fc.container_id
		LEFT JOIN containers AS tc ON to_container_id = tc.container_id
		WHERE time >= $1 AND time <= $2;`, logRange.start, logRange.end)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	return fetchLogsFromQuery(row)
}

func fetchLogsFromQuery(row *sql.Rows) ([]TransactionLogs, error) {
	logs := []TransactionLogs{}

	for row.Next() {
		log := TransactionLogs{}

		var fromMaybeBoxId sql.Null[int]
		var fromMaybeBoxName sql.NullString

		var toMaybeBoxId sql.Null[int]
		var toMaybeBoxName sql.NullString

		if err := row.Scan(&fromMaybeBoxId, &fromMaybeBoxName, &toMaybeBoxId, &toMaybeBoxName, &log.ScryfallId, &log.Quantity); err != nil {
			return nil, err
		}

		if fromMaybeBoxId.Valid && fromMaybeBoxName.Valid {
			log.FromContainer = &TransactionContainer{fromMaybeBoxId.V, fromMaybeBoxName.String}
		}

		if toMaybeBoxId.Valid && toMaybeBoxName.Valid {
			log.ToContainer = &TransactionContainer{toMaybeBoxId.V, toMaybeBoxName.String}
		}

		logs = append(logs, log)
	}

	return logs, nil
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
