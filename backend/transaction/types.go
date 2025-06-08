package transactions

type TransactionLogs struct {
	Transaction_id int `json:"transaction_id"`
	Group_id       int `json:"group_id"`
	From_container int `json:"from_container"`
	To_container   int `json:"to_container"`
	Scryfall_id    int `json:"scryfall_id"`
	Quantity       int `json:"quantity"`
}
