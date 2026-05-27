package controllers_v1

// import (
// 	"fmt"

// 	"github.com/fadilmartias/dilz_code/apps/backend/app/logger"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
// 	"github.com/fadilmartias/dilz_code/apps/backend/config"
// 	"github.com/go-resty/resty/v2"
// 	"github.com/golang-jwt/jwt"

// 	"github.com/gofiber/fiber/v3"
// 	"gorm.io/gorm"
// )

// type MidtransController struct {
// 	BaseController
// 	MidtransService *services.MidtransService
// 	DB              *gorm.DB
// 	Redis           *config.RedisClient
// }

// func NewMidtransController(db *gorm.DB, redis *config.RedisClient, midtransService *services.MidtransService) *MidtransController {
// 	return &MidtransController{DB: db, Redis: redis, MidtransService: midtransService}
// }

// // RenderQrGopay mengambil stream raw binary gambar QR Code GoPay
// // @Summary Render QR Code GoPay
// // @Description Mengambil dan me-render langsung gambar QR Code GoPay dari Midtrans berdasarkan ID pembayaran.
// // @Tags Midtrans
// // @Produce image/png
// // @Param id path string true "ID Transaksi / Order ID"
// // @Success 200 {file} file "Gambar QR Code dalam format PNG"
// // @Failure 500 {object} utils.OrderedErrorResponse "Gagal request ke Midtrans"
// // @Router /v1/midtrans/qr-gopay/{id}/qr-code [get]
// func (ctrl *MidtransController) RenderQrGopay(c fiber.Ctx) error {
// 	id := c.Params("id")
// 	authHeader := utils.MidtransBasicAuth()
// 	url := fmt.Sprintf("https://api.midtrans.com/v2/gopay/%s/qr-code", id)

// 	resp, err := resty.New().
// 		SetHeader("Authorization", authHeader).
// 		SetHeader("Accept", "image/png").
// 		SetDoNotParseResponse(true). // Supaya dapat raw binary
// 		R().
// 		Get(url)

// 	if err != nil {
// 		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
// 			Code:    fiber.StatusInternalServerError,
// 			Message: "Gagal request ke Midtrans",
// 		}, err)
// 	}

// 	// if resp.StatusCode() != 200 {
// 	// 	return utils.ErrorResponse(c, utils.ErrorResponseFormat{
// 	// 		Code:    fiber.StatusInternalServerError,
// 	// 		Message: "Midtrans error",
// 	// 		Details: err,
// 	// 	})
// 	// }

// 	fmt.Println(resp.StatusCode(), "status", "url", url)

// 	// Set response sebagai image/png
// 	c.Set("Content-Type", "image/png")
// 	return c.SendStream(resp.RawBody())
// }

// // Balance mengambil informasi saldo Midtrans
// // @Summary Ambil saldo Midtrans
// // @Description Mengambil informasi total saldo (balance) yang tersedia di akun Midtrans.
// // @Tags Midtrans
// // @Produce json
// // @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data balance Midtrans"
// // @Failure 500 {object} utils.OrderedErrorResponse "Gagal request ke Midtrans"
// // @Router /v1/midtrans/balance [get]
// func (ctrl *MidtransController) Balance(c fiber.Ctx) error {
// 	balance, err := ctrl.MidtransService.GetBalance()
// 	if err != nil {
// 		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
// 			Code:    fiber.StatusInternalServerError,
// 			Message: "Gagal request ke Midtrans",
// 		}, err)
// 	}
// 	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
// 		Message: "Berhasil mendapatkan data balance Midtrans",
// 		Data:    fiber.Map{"balance": balance},
// 	})
// }

// // RefundTransaction memproses pengembalian dana transaksi
// // @Summary Refund Transaksi Midtrans
// // @Description Memproses pengembalian dana (refund) untuk transaksi yang telah berhasil menggunakan Reference ID.
// // @Tags Midtrans
// // @Accept json
// // @Produce json
// // @Param reference_id path string true "Reference ID transaksi yang akan direfund"
// // @Param payload body map[string]interface{} true "Data refund berisi property 'amount' (int) dan 'reason' (string)"
// // @Success 200 {object} utils.OrderedSuccessResponse "Berhasil refund transaksi Midtrans"
// // @Failure 400 {object} utils.OrderedErrorResponse "Invalid request body"
// // @Failure 500 {object} utils.OrderedErrorResponse "Gagal refund transaksi Midtrans"
// // @Router /v1/midtrans/refund/{reference_id} [post]
// func (ctrl *MidtransController) RefundTransaction(c fiber.Ctx) error {
// 	reference_id := c.Params("reference_id")
// 	var input struct {
// 		Amount int    `json:"amount"`
// 		Reason string `json:"reason"`
// 	}
// 	if err := c.Bind().Body(&input); err != nil {
// 		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
// 			Code:    fiber.StatusBadRequest,
// 			Message: "Invalid request body",
// 		}, err)
// 	}
// 	amount := input.Amount
// 	reason := input.Reason
// 	resp, err := ctrl.MidtransService.RefundTransaction(reference_id, amount, reason)
// 	if err != nil {
// 		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
// 			Code:    fiber.StatusInternalServerError,
// 			Message: "Gagal refund transaksi Midtrans",
// 		}, err)
// 	}
// 	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
// 		Message: "Berhasil refund transaksi Midtrans",
// 		Data:    resp,
// 	})
// }

// // Webhook menangani notifikasi otomatis dari Midtrans
// // @Summary Midtrans Webhook Handler
// // @Description Menerima dan memproses payload notifikasi status pembayaran dari server Midtrans.
// // @Tags Midtrans
// // @Accept json
// // @Produce json
// // @Param payload body requests.MidtransWebhookRequestAttributes true "Data payload notifikasi transaksi dari Midtrans"
// // @Success 200 {object} utils.OrderedSuccessResponse "Berhasil memproses webhook dan update status transaksi"
// // @Failure 400 {object} utils.OrderedErrorResponse "Gagal parsing payload webhook"
// // @Failure 500 {object} utils.OrderedErrorResponse "Gagal memproses dan menyimpan perubahan status transaksi"
// // @Router /v1/midtrans/webhook [post]
// func (ctrl *MidtransController) Webhook(c fiber.Ctx) error {
// 	userData := c.Value("user")
// 	role := ""
// 	if userData != nil {
// 		claims, ok := userData.(jwt.MapClaims)
// 		if !ok {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Unauthorized: invalid token claims",
// 			})
// 		}

// 		if r, ok := claims["role"].(string); ok {
// 			role = r
// 		}
// 	}

// 	logger.Infof("Midtrans webhook received: %s", c.Request().Body)

// 	var input requests.MidtransWebhookRequestAttributes
// 	if err := c.Bind().Body(&input); err != nil {
// 		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
// 			Code:    fiber.StatusBadRequest,
// 			Message: "Gagal update status transaksi",
// 		}, err)
// 	}

// 	message, err := ctrl.MidtransService.HandleWebhook(input, role)
// 	if err != nil {
// 		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
// 			Code:    fiber.StatusInternalServerError,
// 			Message: "Gagal update status transaksi",
// 		}, err)
// 	}

// 	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
// 		Message: message,
// 	})
// }
