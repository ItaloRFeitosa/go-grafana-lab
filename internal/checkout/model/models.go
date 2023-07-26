package model

type Payment struct {
	CorrelationID string `json:"correlationId"`
	Status        string `json:"status"`
}
