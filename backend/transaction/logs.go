package transactions

import (
	"github.com/jnie1/MTGViewer-V2/database"
)

func FetchLogs() ([]TransactionLogs, error) {
	listOfLogs := []TransactionLogs{}

	db := database.Instance()
	row, err := db.Query(`
		SELECT *
		FROM transaction
	`)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		logs := TransactionLogs{}

		if err = row.Scan(&logs.Transaction_id); err != nil {
			return nil, err
		}

		listOfLogs = append(listOfLogs, logs)
	}

	return listOfLogs, nil
}
