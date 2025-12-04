package requests

type DigiflazzWebhookRequestAttributes struct {
	Data struct {
		RefID          string `json:"ref_id"`
		CustomerNo     string `json:"customer_no"`
		BuyerSKUCode   string `json:"buyer_sku_code"`
		Message        string `json:"message"`
		Status         string `json:"status"`
		RC             string `json:"rc"`
		BuyerLastSaldo int    `json:"buyer_last_saldo"`
		SN             string `json:"sn"`
		Price          int    `json:"price"`
		Tele           string `json:"tele"`
		WA             string `json:"wa"`
	} `json:"data"`
}
