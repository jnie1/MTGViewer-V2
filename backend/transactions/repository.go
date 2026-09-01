package transactions

import (
	"context"
	"database/sql"
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
		SELECT MIN(lg.time) AS start, MAX(lg.time) AS end
		FROM log_groups AS lg
		WHERE lg.log_group_id = $1 OR lg.log_group_id = $2;`, group1, group2)

	var logRange LogRange
	err := row.Scan(&logRange.Start, &logRange.End)

	return logRange, err
}

func GetTransactions(ctx context.Context) ([]CardTransaction, error) {
	db := database.Instance()
	row, err := db.QueryContext(ctx, `
		SELECT lg.log_group_id, lg.time, COALESCE(SUM(t.amount), 0) AS total, lg.description
		FROM log_groups AS lg
		LEFT JOIN transactions AS t ON t.log_group_id = lg.log_group_id
		GROUP BY lg.log_group_id
		ORDER BY lg.time DESC;`)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	var transactions []CardTransaction

	for row.Next() {
		var transaction CardTransaction
		if err := row.Scan(&transaction.GroupId, &transaction.Time, &transaction.Total, &transaction.Description); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func GetTransaction(ctx context.Context, groupId uuid.UUID) (CardTransaction, error) {
	db := database.Instance()
	row := db.QueryRowContext(ctx, `
		SELECT lg.log_group_id, lg.time, COALESCE(SUM(t.amount), 0) AS total, lg.description
		FROM log_groups AS lg
		LEFT JOIN transactions AS t ON t.log_group_id = lg.log_group_id
		WHERE lg.log_group_id = $1
		GROUP BY lg.log_group_id;`, groupId)

	var transaction CardTransaction
	err := row.Scan(&transaction.GroupId, &transaction.Time, &transaction.Total, &transaction.Description)

	return transaction, err
}

func GetLogs(ctx context.Context, groupId uuid.UUID) ([]CardLogPreview, error) {
	db := database.Instance()
	row, err := db.QueryContext(ctx, `
		SELECT t.scryfall_id, t.from_container_id, t.to_container_id, t.amount
		FROM transactions AS t
		WHERE t.log_group_id = $1;`, groupId)

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
		SELECT t.scryfall_id, t.from_container_id, t.to_container_id, t.amount
		FROM transactions AS t
		JOIN log_groups AS lg ON lg.log_group_id = t.log_group_id
		WHERE lg.time >= $1 AND lg.time <= $2;`, logRange.Start, logRange.End)

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

func UpdateDescription(ctx context.Context, groupId uuid.UUID, description *string) error {
	db := database.Instance()
	res, err := db.ExecContext(ctx, `
		UPDATE log_groups
		SET description = $1
		WHERE log_group_id = $2
		RETURNING log_group_id;`, description, groupId)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows < 1 {
		return sql.ErrNoRows
	}

	return nil
}

func LogCollectionChanges(ctx context.Context, changes []containers.ContainerChanges) error {
	now := time.Now().UTC()
	groupId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	db := database.Instance()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO log_groups (log_group_id, time)
		VALUES ($1, $2);`, groupId, now)

	if err != nil {
		return err
	}

	var vals []string
	var args []any
	i := 0
	for _, change := range changes {
		for _, request := range change.Requests {
			if request.Delta == 0 {
				continue
			}
			vals = append(vals, fmt.Sprintf("($%d::uuid, $%d, $%d, $%d::uuid, $%d)", i+1, i+2, i+3, i+4, i+5))
			i += 5
			switch {
			case request.Delta > 0:
				args = append(args, groupId, nil, change.ContainerId, request.ScryfallId, request.Delta)
			case request.Delta < 0:
				args = append(args, groupId, change.ContainerId, nil, request.ScryfallId, -request.Delta)
			}
		}
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactions (log_group_id, from_container_id, to_container_id, scryfall_id, amount)
		VALUES `+strings.Join(vals, ", ")+`;`, args...)

	if err != nil {
		return err
	}

	return tx.Commit()
}
