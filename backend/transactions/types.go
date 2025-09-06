package transactions

import "github.com/google/uuid"

type TransactionLogs struct {
	TransactionId int                   `json:"transactionId"`
	GroupId       uuid.UUID             `json:"groupId"`
	FromContainer *TransactionContainer `json:"fromContainer"`
	ToContainer   *TransactionContainer `json:"toContainer"`
	ScryfallId    uuid.UUID             `json:"scryfallId"`
	Quantity      int                   `json:"quantity"`
}

type TransactionContainer struct {
	ContainerId int    `json:"containerId"`
	Name        string `json:"name"`
}
