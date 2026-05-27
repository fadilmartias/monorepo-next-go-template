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
