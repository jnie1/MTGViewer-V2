package transactions

type TransactionLogs struct {
	TransactionId int `json:"transaction_id"`
	GroupId       int `json:"group_id"`
	FromContainer int `json:"from_container"`
	ToContainer   int `json:"to_container"`
	ScryfallId    int `json:"scryfall_id"`
	Quantity       int `json:"quantity"`
}
