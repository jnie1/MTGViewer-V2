package transactions

type TransactionLogs struct {
	TransactionId int `json:"transactionId"`
	GroupId       int `json:"groupId"`
	FromContainer int `json:"fromContainer"`
	ToContainer   int `json:"toContainer"`
	ScryfallId    int `json:"scryfallId"`
	Quantity       int `json:"quantity"`
}
