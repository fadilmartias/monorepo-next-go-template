package services

// import (
// 	"errors"
// 	"os"
// 	"strconv"
// 	"time"

// 	"github.com/fadilmartias/dilz_code/apps/backend/app/dto"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/logger"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
// 	"github.com/fadilmartias/dilz_code/apps/backend/config"
// 	"github.com/go-resty/resty/v2"
// 	"github.com/tidwall/gjson"
// 	"gorm.io/gorm"
// )

// type MidtransService struct {
// 	DB                           *gorm.DB
// 	Redis                        *config.RedisClient
// 	PaymentService               *PaymentService
// 	ProductTransactionRepository *repositories.ProductTransactionRepository
// }

// func NewMidtransService(db *gorm.DB, redis *config.RedisClient, paymentService *PaymentService, productTransactionRepository *repositories.ProductTransactionRepository) *MidtransService {
// 	return &MidtransService{DB: db, Redis: redis, PaymentService: paymentService, ProductTransactionRepository: productTransactionRepository}
// }

// type TransactionDetails struct {
// 	OrderID     string `json:"order_id"`
// 	GrossAmount int    `json:"gross_amount"`
// }

// type ItemDetail struct {
// 	ID           string `json:"id"`
// 	Price        int    `json:"price"`
// 	Quantity     int    `json:"quantity"`
// 	Name         string `json:"name"`
// 	Brand        string `json:"brand"`
// 	Category     string `json:"category"`
// 	MerchantName string `json:"merchant_name"`
// 	URL          string `json:"url"`
// }

// type MidtransBaseRequestAttribute struct {
// 	TransactionDetails TransactionDetails           `json:"transaction_details"`
// 	ItemDetails        []ItemDetail                 `json:"item_details"`
// 	CustomerDetails    *dto.MidtransCustomerDetails `json:"customer_details,omitempty"`
// }

// type MidtransSNAPRequestAttribute struct {
// 	MidtransBaseRequestAttribute
// 	EnabledPayments []string   `json:"enabled_payments"`
// 	Callbacks       Callbacks  `json:"callbacks"`
// 	PageExpiry      PageExpiry `json:"page_expiry"`
// 	Expiry          Expiry     `json:"expiry"`
// }

// type MidtransCoreAPIRequestAttribute struct {
// 	MidtransBaseRequestAttribute
// 	PaymentType  string        `json:"payment_type"`
// 	CustomExpiry *CustomExpiry `json:"custom_expiry,omitempty"`
// }

// type MidtransOTCRequestAttribute struct {
// 	MidtransCoreAPIRequestAttribute
// 	Cstore Cstore `json:"cstore"`
// }

// type MidtransBankTransferRequestAttribute struct {
// 	MidtransCoreAPIRequestAttribute
// 	BankTransfer BankTransfer `json:"bank_transfer"`
// }

// type MidtransQrisRequestAttribute struct {
// 	MidtransCoreAPIRequestAttribute
// 	Qris Qris `json:"qris"`
// }

// type Cstore struct {
// 	Store             string  `json:"store"`
// 	Message           *string `json:"message,omitempty"`
// 	AlfamartFreeText1 *string `json:"alfamart_free_text_1,omitempty"`
// 	AlfamartFreeText2 *string `json:"alfamart_free_text_2,omitempty"`
// 	AlfamartFreeText3 *string `json:"alfamart_free_text_3,omitempty"`
// }

// type CustomExpiry struct {
// 	OrderTime      string `json:"order_time"`
// 	ExpiryDuration int    `json:"expiry_duration"`
// 	Unit           string `json:"unit"`
// }

// type MidtransShopeepayRequestAttribute struct {
// 	MidtransCoreAPIRequestAttribute
// 	Shopeepay Shopeepay `json:"shopeepay"`
// }

// type Shopeepay struct {
// 	CallbackURL string `json:"callback_url"`
// }

// type Callbacks struct {
// 	Finish string `json:"finish"`
// }

// type Qris struct {
// 	Acquirer string `json:"acquirer"`
// }

// type BankTransfer struct {
// 	Bank string `json:"bank"`
// }

// type Gopay struct {
// 	EnableCallback bool   `json:"enable_callback"`
// 	CallbackURL    string `json:"callback_url"`
// }

// type MidtransGopayRequestAttribute struct {
// 	MidtransCoreAPIRequestAttribute
// 	Gopay Gopay `json:"gopay"`
// }

// type Echannel struct {
// 	BillInfo1 string `json:"bill_info1"`
// 	BillInfo2 string `json:"bill_info2"`
// }

// type MidtransEchannelRequestAttribute struct {
// 	MidtransCoreAPIRequestAttribute
// 	Echannel Echannel `json:"echannel"`
// }

// type PageExpiry struct {
// 	Duration int    `json:"duration"`
// 	Unit     string `json:"unit"`
// }

// type Expiry struct {
// 	StartTime string `json:"start_time"`
// 	Unit      string `json:"unit"`
// 	Duration  int    `json:"duration"`
// }

// func SNAPRequest(body any) (*resty.Response, error) {
// 	midtransConfig := config.LoadMidtransConfig()
// 	client := resty.New().
// 		SetHeader("Content-Type", "application/json").
// 		SetHeader("Accept", "application/json").
// 		SetHeader("Authorization", utils.MidtransBasicAuth())
// 	return client.R().
// 		SetBody(body).
// 		Post(midtransConfig.SnapAPIURL)
// }

// func (s *MidtransService) CreateSNAPTransaction(body dto.CreateTransactionParams) (*resty.Response, error) {
// 	requestMidtrans := MidtransSNAPRequestAttribute{
// 		MidtransBaseRequestAttribute: MidtransBaseRequestAttribute{
// 			TransactionDetails: TransactionDetails{
// 				OrderID:     body.OrderId,
// 				GrossAmount: int(body.TotalPrice),
// 			},
// 			ItemDetails: []ItemDetail{
// 				{
// 					ID:           body.ProductVariant.ID,
// 					Price:        int(body.ItemPrice),
// 					Quantity:     body.Qty,
// 					Name:         body.ProductVariant.Name,
// 					Brand:        body.Product.Name,
// 					MerchantName: config.LoadAppConfig().Name,
// 					URL:          os.Getenv("FE_URL") + "/topup/" + body.Product.Slug,
// 				},
// 				{
// 					ID:           "service_fee",
// 					Price:        int(body.TotalFee),
// 					Quantity:     1,
// 					Name:         "Service Fee",
// 					MerchantName: config.LoadAppConfig().Name,
// 				},
// 			},
// 		},
// 		EnabledPayments: []string{body.PaymentMethod.InvoiceCode},
// 		Callbacks: Callbacks{
// 			Finish: os.Getenv("FE_URL") + "/payment/" + body.OrderId,
// 		},
// 		PageExpiry: PageExpiry{
// 			Duration: body.PaymentMethod.ExpirySeconds / 60,
// 			Unit:     "minutes",
// 		},
// 		Expiry: Expiry{
// 			StartTime: time.Now().Format("2006-01-02 15:04:05 +0700"),
// 			Unit:      "minutes",
// 			Duration:  body.PaymentMethod.ExpirySeconds / 60,
// 		},
// 	}

// 	if body.CustomerDetails != nil {
// 		requestMidtrans.CustomerDetails = body.CustomerDetails
// 	}

// 	if body.TotalDiscount != nil && *body.TotalDiscount > 0 {
// 		for _, promo := range body.Promotions {
// 			requestMidtrans.ItemDetails = append(requestMidtrans.ItemDetails, ItemDetail{
// 				ID:           promo.PromotionID,
// 				Price:        int(-promo.TotalDiscount),
// 				Quantity:     1,
// 				Name:         "Promo " + promo.PromotionID,
// 				MerchantName: config.LoadAppConfig().Name,
// 			})
// 		}
// 	}
// 	resp, err := SNAPRequest(requestMidtrans)
// 	return resp, err
// }

// func (s *MidtransService) CoreAPIRequest(endpoint string, body any) (*resty.Response, error) {
// 	midtransConfig := config.LoadMidtransConfig()
// 	client := resty.New().
// 		SetBaseURL(midtransConfig.BaseURL).
// 		SetHeader("Content-Type", "application/json").
// 		SetHeader("Accept", "application/json").
// 		SetHeader("Authorization", utils.MidtransBasicAuth())
// 	return client.R().
// 		SetBody(body).
// 		Post(endpoint)
// }

// func (s *MidtransService) CreateCoreAPIPayment(body dto.CreateTransactionParams) (*models.ProductTransaction, error) {
// 	// Base request
// 	requestMidtrans := MidtransCoreAPIRequestAttribute{
// 		MidtransBaseRequestAttribute: MidtransBaseRequestAttribute{
// 			TransactionDetails: TransactionDetails{
// 				OrderID:     body.RefId,
// 				GrossAmount: int(body.TotalPrice),
// 			},
// 			ItemDetails: []ItemDetail{
// 				{
// 					ID:           body.ProductVariant.ID,
// 					Price:        int(body.ItemPrice),
// 					Quantity:     body.Qty,
// 					Name:         body.ProductVariant.Name,
// 					Brand:        body.Product.Name,
// 					MerchantName: config.LoadAppConfig().Name,
// 					URL:          os.Getenv("FE_URL") + "/topup/" + body.Product.Slug,
// 				},
// 				{
// 					ID:           "service_fee",
// 					Price:        int(body.TotalFee),
// 					Quantity:     1,
// 					Name:         "Service Fee",
// 					MerchantName: config.LoadAppConfig().Name,
// 				},
// 			},
// 		},
// 		PaymentType: body.PaymentMethod.Type,
// 		CustomExpiry: &CustomExpiry{
// 			OrderTime:      time.Now().Format("2006-01-02 15:04:05 +0700"),
// 			ExpiryDuration: body.PaymentMethod.ExpirySeconds,
// 			Unit:           "second",
// 		},
// 	}

// 	if body.CustomerDetails != nil {
// 		requestMidtrans.CustomerDetails = body.CustomerDetails
// 	}

// 	if body.TotalDiscount != nil && *body.TotalDiscount > 0 {
// 		for _, promo := range body.Promotions {
// 			requestMidtrans.ItemDetails = append(requestMidtrans.ItemDetails, ItemDetail{
// 				ID:           promo.PromotionID,
// 				Price:        int(-promo.TotalDiscount),
// 				Quantity:     1,
// 				Name:         "Promo " + promo.PromotionID,
// 				MerchantName: config.LoadAppConfig().Name,
// 			})
// 		}
// 	}

// 	// Pilih struct final berdasarkan payment type
// 	var finalRequest any = requestMidtrans
// 	switch body.PaymentMethod.Type {
// 	case "qris":
// 		finalRequest = MidtransQrisRequestAttribute{
// 			MidtransCoreAPIRequestAttribute: requestMidtrans,
// 			Qris:                            Qris{Acquirer: "gopay"},
// 		}
// 	case "bank_transfer":
// 		finalRequest = MidtransBankTransferRequestAttribute{
// 			MidtransCoreAPIRequestAttribute: requestMidtrans,
// 			BankTransfer:                    BankTransfer{Bank: body.PaymentMethod.Code},
// 		}
// 	case "shopeepay":
// 		finalRequest = MidtransShopeepayRequestAttribute{
// 			MidtransCoreAPIRequestAttribute: requestMidtrans,
// 			Shopeepay:                       Shopeepay{CallbackURL: os.Getenv("FE_URL") + "/payment/" + body.OrderId},
// 		}
// 	case "cstore":
// 		finalRequest = MidtransOTCRequestAttribute{
// 			MidtransCoreAPIRequestAttribute: requestMidtrans,
// 			Cstore:                          Cstore{Store: body.PaymentMethod.Code},
// 		}
// 	case "gopay":
// 		finalRequest = MidtransGopayRequestAttribute{
// 			MidtransCoreAPIRequestAttribute: requestMidtrans,
// 			Gopay: Gopay{
// 				EnableCallback: true,
// 				CallbackURL:    os.Getenv("FE_URL") + "/payment/" + body.OrderId,
// 			},
// 		}
// 	case "echannel":
// 		finalRequest = MidtransEchannelRequestAttribute{
// 			MidtransCoreAPIRequestAttribute: requestMidtrans,
// 			Echannel: Echannel{
// 				BillInfo1: "Payment For " + body.ProductVariant.Name,
// 				BillInfo2: config.LoadAppConfig().Name,
// 			},
// 		}
// 	}

// 	logger.Infof("Final Request: %v", finalRequest)
// 	utils.Dump(finalRequest)

// 	// Kirim request
// 	resp, err := s.CoreAPIRequest("/v2/charge", finalRequest)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Pakai gjson untuk parsing cepat
// 	jsonStr := resp.String()
// 	result := gjson.Parse(jsonStr)

// 	// Validasi
// 	logger.Infof("Midtrans create transaction response: %v", jsonStr)
// 	if resp.StatusCode() != 200 || result.Get("status_code").String() != "201" {
// 		return nil, errors.New("failed to create transaction")
// 	}

// 	// Ambil data VA number
// 	var vaNumber, bank *string
// 	if va := result.Get("va_numbers.0.va_number"); va.Exists() {
// 		vaNumber = utils.StringPtr(va.String())
// 		bank = utils.StringPtr(result.Get("va_numbers.0.bank").String())
// 	}
// 	// Cek permata
// 	if permata := result.Get("permata_va_number"); permata.Exists() {
// 		vaNumber = utils.StringPtr(permata.String())
// 		bank = utils.StringPtr("permata")
// 	}

// 	// Qr String
// 	qrString := utils.StringPtr(result.Get("qr_string").String())

// 	// Ambil QR dan Checkout URL
// 	var qrURL, checkoutURL *string
// 	qr := result.Get(`actions.#(name=="generate-qr-code-v2").url`)
// 	if qr.Exists() {
// 		qrURL = utils.StringPtr(qr.String())
// 	}
// 	checkout := result.Get(`actions.#(name=="deeplink-redirect").url`)
// 	if checkout.Exists() {
// 		checkoutURL = utils.StringPtr(checkout.String())
// 	} else if redirect := result.Get("redirect_url"); redirect.Exists() {
// 		checkoutURL = utils.StringPtr(redirect.String())
// 	}

// 	actionsJSON := result.Get("actions").Raw

// 	// Simpan transaksi
// 	transaction := models.ProductTransaction{
// 		ID:               body.ID,
// 		ProductOrderID:   body.OrderId,
// 		PaymentMethodID:  body.PaymentMethod.ID,
// 		ReferenceID:      body.RefId,
// 		PaymentRequestID: utils.StringPtr(result.Get("transaction_id").String()),
// 		Currency:         "IDR",
// 		Region:           "ID",
// 		PGFee:            body.PGFee,
// 		Fee:              body.TotalFee,
// 		BasePrice:        body.BasePrice,
// 		PriceAfterMargin: body.PriceAfterMargin,
// 		Amount:           body.TotalPrice,
// 		PgStatus:         utils.StringPtr(result.Get("transaction_status").String()),
// 		Actions:          utils.StringPtr(actionsJSON),
// 		InvoiceURL:       utils.StringPtr(result.Get("invoice_url").String()),
// 		PaymentCode:      utils.StringPtr(result.Get("payment_code").String()),
// 		VANumber:         vaNumber,
// 		Bank:             bank,
// 		CheckoutURL:      checkoutURL,
// 		QRString:         qrString,
// 		QRURL:            qrURL,
// 		BillKey:          utils.StringPtr(result.Get("bill_key").String()),
// 		BillerCode:       utils.StringPtr(result.Get("biller_code").String()),
// 		TransactionType:  "normal", // TODO: UBAH KALO MAU PREORDER NANTI
// 		Note:             nil,
// 		Metadata:         nil,
// 		ExpiredAt:        time.Now().Add(time.Second * time.Duration(body.PaymentMethod.ExpirySeconds)),
// 	}

// 	logger.Infof("Midtrans create transaction response: %v", jsonStr)
// 	return &transaction, nil
// }

// func (s *MidtransService) RefundTransaction(reference_id string, amount int, reason string) (*resty.Response, error) {
// 	body := map[string]any{
// 		"amount": amount,
// 		"reason": reason,
// 	}
// 	resp, err := s.CoreAPIRequest("/v2/"+reference_id+"/refund", body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	logger.Infof("Midtrans refund response: %v", resp.String())
// 	if gjson.Get(resp.String(), "status_code").String() != "200" {
// 		return nil, errors.New(gjson.Get(resp.String(), "status_message").String())
// 	}
// 	return resp, nil
// }

// func (s *MidtransService) GetBalance() (int64, error) {
// 	client := resty.New().
// 		SetHeader("Content-Type", "application/json").
// 		SetHeader("Accept", "application/json").
// 		SetHeader("Authorization", utils.MidtransBasicAuth())
// 	resp, err := client.R().
// 		Get("https://app.sandbox.midtrans.com/iris/api/v1/balance")
// 	if err != nil {
// 		return 0, err
// 	}
// 	logger.Infof("Midtrans balance response: %v", resp.String())
// 	balance := gjson.Get(resp.String(), "balance").Int()
// 	return balance, nil
// }

// func (s *MidtransService) HandleWebhook(body requests.MidtransWebhookRequestAttributes, role string) (message string, err error) {
// 	if role != "admin" {
// 		_, err = utils.MidtransVerifyWebhookSignature(body.OrderID, body.StatusCode, body.GrossAmount, body.SignatureKey)
// 		if err != nil {
// 			return "", err
// 		}
// 	}
// 	code, err := strconv.Atoi(body.StatusCode)
// 	if err != nil {
// 		return "", errors.New("invalid status code")
// 	}

// 	if code < 200 || code > 299 {
// 		return "", errors.New("invalid status code")
// 	}

// 	if body.FraudStatus != "accept" {
// 		return "", errors.New("invalid fraud status")
// 	}

// 	var orderStatus string
// 	var paymentStatus string
// 	pgStatus := body.TransactionStatus

// 	switch pgStatus {
// 	case "settlement":
// 		orderStatus = "waiting delivery"
// 		paymentStatus = "paid"
// 	case "capture":
// 		orderStatus = "waiting delivery"
// 		paymentStatus = "paid"
// 	case "expire":
// 		orderStatus = "expired"
// 		paymentStatus = "expired"
// 	case "pending":
// 		orderStatus = "waiting payment"
// 		paymentStatus = "pending"
// 	case "authorize":
// 		orderStatus = "waiting payment"
// 		paymentStatus = "pending"
// 	case "refund":
// 		orderStatus = "refunded"
// 		paymentStatus = "refunded"
// 	case "partial_refund":
// 		orderStatus = "partial_refunded"
// 		paymentStatus = "partial_refunded"
// 	case "failure":
// 		orderStatus = "failed"
// 		paymentStatus = "failed"
// 	case "deny":
// 		orderStatus = "failed"
// 		paymentStatus = "failed"
// 	case "cancel":
// 		orderStatus = "failed"
// 		paymentStatus = "failed"
// 	default:
// 		return "", errors.New("invalid transaction status")
// 	}

// 	message, err = s.PaymentService.HandleWebhook(body.OrderID, orderStatus, paymentStatus, pgStatus)
// 	if err != nil {
// 		return "", err
// 	}

// 	return message, nil
// }
