package webhook

type Payment struct {
	CorrelationID string `json:"correlationId"`
	Status        string `json:"status"`
}
