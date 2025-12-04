package requests

type IPaymuWebhookRequestAttributes struct {
	TrxID       string `json:"trxId"`
	Status      string `json:"status"`
	Amount      int    `json:"amount"`
	StatusCode  string `json:"statusCode"`
	Sid         string `json:"sid"`
	ReferenceID string `json:"referenceId"`
}
