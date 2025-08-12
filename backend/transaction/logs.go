package transactions

import (
	"github.com/jnie1/MTGViewer-V2/database"
)

func FetchLogs() ([]TransactionLogs, error) {
	listOfLogs := []TransactionLogs{}

	db := database.Instance()
	row, err := db.Query(`
		SELECT transaction_id, group_id, from_container, to_container, scryfall_id, amount 	
		FROM transactions`)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		logs := TransactionLogs{}

		if err = row.Scan(
			&logs.TransactionId,
			&logs.GroupId,
			&logs.FromContainer,
			&logs.ToContainer,
			&logs.ScryfallId,
			&logs.Quantity); err != nil {
			return nil, err
		}

		listOfLogs = append(listOfLogs, logs)
	}

	return listOfLogs, nil
}
