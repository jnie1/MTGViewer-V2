package transactions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/containers"
	"github.com/jnie1/MTGViewer-V2/database"
)

func GetTimeRange(ctx context.Context, group1, group2 uuid.UUID) (LogRange, error) {
	db := database.Instance()

	row := db.QueryRowContext(ctx, `
		SELECT MIN(time) AS start, MAX(time) AS end
		FROM transactions
		WHERE group_id = $1 OR group_id = $2;`, group1, group2)

	var logRange LogRange
	err := row.Scan(&logRange.Start, &logRange.End)

	return logRange, err
}

func GetTransactions(ctx context.Context) ([]CardTransaction, error) {
	db := database.Instance()
	row, err := db.QueryContext(ctx, `
		SELECT group_id, time, SUM(amount) AS total
		FROM transactions
		GROUP BY group_id, time
		ORDER BY time DESC;`)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	var transactions []CardTransaction

	for row.Next() {
		var transaction CardTransaction
		if err := row.Scan(&transaction.GroupId, &transaction.Time, &transaction.Total); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func GetLogs(ctx context.Context, groupId uuid.UUID) ([]CardLogPreview, error) {
	db := database.Instance()
	row, err := db.QueryContext(ctx, `
		SELECT scryfall_id, from_container_id, to_container_id, amount
		FROM transactions
		WHERE group_id = $1;`, groupId)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	var logs []CardLogPreview

	for row.Next() {
		var log CardLogPreview
		if err := row.Scan(&log.ScryfallId, &log.FromContainerId, &log.ToContainerId, &log.Amount); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func GetLogsFromRange(ctx context.Context, logRange LogRange) ([]CardLogPreview, error) {
	db := database.Instance()
	row, err := db.QueryContext(ctx, `
		SELECT scryfall_id, from_container_id, to_container_id, amount
		FROM transactions
		WHERE time >= $1 AND time <= $2;`, logRange.Start, logRange.End)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	var logs []CardLogPreview

	for row.Next() {
		var log CardLogPreview
		if err := row.Scan(&log.ScryfallId, &log.FromContainerId, &log.ToContainerId, &log.Amount); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func LogCollectionChanges(ctx context.Context, changes []containers.ContainerChanges) error {
	now := time.Now().UTC()
	groupId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	var valueStatements []string
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

	db := database.Instance()
	allValues := strings.Join(valueStatements, ", ")

	_, err = db.ExecContext(ctx, `
		INSERT INTO transactions (group_id, from_container_id, to_container_id, scryfall_id, amount, time)
		VALUES `+allValues+`;`)

	return err
}
