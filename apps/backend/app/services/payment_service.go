package services

// import (
// 	"errors"
// 	"fmt"
// 	"strings"
// 	"time"

// 	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/usecases"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
// 	"github.com/fadilmartias/dilz_code/apps/backend/config"
// 	"github.com/gofiber/fiber/v3"
// 	"gorm.io/gorm"
// )

// type PaymentService struct {
// 	DB                           *gorm.DB
// 	Redis                        *config.RedisClient
// 	ProductOrderRepository       *repositories.ProductOrderRepository
// 	ProductOrderDetailRepository *repositories.ProductOrderDetailRepository
// 	ProductTransactionRepository *repositories.ProductTransactionRepository
// 	ProductVariantRepository     *repositories.ProductVariantRepository
// 	TelegramService              *TelegramService
// 	FonnteService                *FonnteService
// }

// type WebhookParams struct {
// 	RefId string `json:"refId"`
// 	Data  any    `json:"data"`
// 	Event string `json:"event"`
// 	Sign  string `json:"sign"`
// }

// func NewPaymentService(db *gorm.DB, redis *config.RedisClient, productOrderRepository *repositories.ProductOrderRepository, productOrderDetailRepository *repositories.ProductOrderDetailRepository, productTransactionRepository *repositories.ProductTransactionRepository, productVariantRepository *repositories.ProductVariantRepository, fonnteService *FonnteService) *PaymentService {
// 	return &PaymentService{DB: db, Redis: redis, ProductOrderRepository: productOrderRepository, ProductOrderDetailRepository: productOrderDetailRepository, ProductTransactionRepository: productTransactionRepository, ProductVariantRepository: productVariantRepository, TelegramService: NewTelegramService(), FonnteService: fonnteService}
// }

// func (s *PaymentService) HandleWebhook(refID string, orderStatus string, paymentStatus string, pgStatus string) (string, error) {

// 	if paymentStatus == "pending" {
// 		return "Berhasil update status transaksi (pending)", nil
// 	}

// 	// Get Transaction
// 	transaction, err := s.ProductTransactionRepository.FindByRefID(refID)
// 	if err != nil {
// 		return "", err
// 	}
// 	transactionType := string(transaction.TransactionType)

// 	productPreorder := &models.ProductOrderPreorder{}
// 	installmentSchedule := &models.ProductOrderInstallmentSchedule{}
// 	nextInstallmentSchedule := &models.ProductOrderInstallmentSchedule{}

// 	var messageWA string
// 	ign := ""
// 	gameUsername := ""
// 	if transaction.ProductOrder.GameUsername != nil && *transaction.ProductOrder.GameUsername != "" {
// 		ign = fmt.Sprintf("IGN : %s\n", *transaction.ProductOrder.GameUsername)
// 		gameUsername = *transaction.ProductOrder.GameUsername
// 	}

// 	if strings.HasPrefix(transactionType, "preorder") {
// 		err = s.DB.Where("product_order_id = ?", transaction.ProductOrderID).First(productPreorder).Error
// 		if err != nil {
// 			return "", err
// 		}
// 		err = s.DB.Where("id = ?", transaction.RelatedInstallmentID).First(installmentSchedule).Error
// 		if err != nil {
// 			return "", err
// 		}

// 		now := time.Now()
// 		productPreorder.PaidInstallments += 1

// 		if productPreorder.PaidInstallments < productPreorder.InstallmentCount-1 || productPreorder.PaidInstallments == productPreorder.InstallmentCount-1 {
// 			err = s.DB.Where("product_order_preorder_id = ? AND installment_number = ?", productPreorder.ID, installmentSchedule.ID).First(nextInstallmentSchedule).Error
// 			if err != nil {
// 				return "", err
// 			}
// 		}
// 		productPreorder.RemainingAmount -= (transaction.Amount - transaction.PGFee)
// 		installmentSchedule.PaidAt = &now
// 		installmentSchedule.Status = "paid"
// 		nextDueDate := nextInstallmentSchedule.DueDate.Format("02/01/2006 15:04:05")

// 		orderDetail := fmt.Sprintf(
// 			"Detail Pesanan:\n"+
// 				"Order ID: %s\n"+
// 				"Product : %s\n"+
// 				"Product Variant : %s\n"+
// 				"Qty : %d\n"+
// 				"%s",
// 			transaction.ProductOrderID,
// 			transaction.ProductOrder.Product.Name,
// 			transaction.ProductOrder.ProductOrderDetails[0].ProductVariant.Name,
// 			transaction.ProductOrder.TotalQty,
// 			ign,
// 		)

// 		if orderStatus == "waiting delivery" {

// 			if productPreorder.PaidInstallments == productPreorder.InstallmentCount &&
// 				productPreorder.RemainingAmount <= 0 {

// 				orderStatus = "waiting delivery"
// 				productPreorder.Status = "full_paid"

// 				messageWA = fmt.Sprintf(
// 					"Pembayaran preorder kamu telah lunas! 🎉\n"+
// 						"Order kamu sekarang masuk tahap *waiting delivery*.\n\n%s",
// 					orderDetail,
// 				)

// 			} else if productPreorder.PaidInstallments < productPreorder.InstallmentCount-1 {

// 				orderStatus = "on_progress"
// 				productPreorder.Status = "installment_in_progress"

// 				messageWA = fmt.Sprintf(
// 					"Pembayaran preorder kamu berhasil! 🎉\n"+
// 						"Tanggal jatuh tempo berikutnya: %s, pastikan membayar sebelum jatuh tempo!\n\n%s",
// 					nextDueDate,
// 					orderDetail,
// 				)

// 			} else if productPreorder.PaidInstallments == productPreorder.InstallmentCount-1 {

// 				orderStatus = "on_progress"
// 				productPreorder.Status = "waiting_full_payment"

// 				messageWA = fmt.Sprintf(
// 					"Pembayaran preorder kamu berhasil! 🎉\n"+
// 						"Tanggal jatuh tempo pelunasan: %s, pastikan kamu membayar sebelum jatuh tempo!\n\n%s",
// 					nextDueDate,
// 					orderDetail,
// 				)
// 			}
// 		}

// 	}

// 	// Update ProductOrder
// 	rowsAffected, err := s.ProductOrderRepository.UpdateByID(transaction.ProductOrderID, map[string]any{
// 		"status": orderStatus,
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	if rowsAffected == 0 {
// 		return "", errors.New("product order not found")
// 	}
// 	// Update ProductOrderDetail
// 	rowsAffected, err = s.ProductOrderDetailRepository.UpdateByOrderID(transaction.ProductOrderID, map[string]any{
// 		"status": orderStatus,
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	if rowsAffected == 0 {
// 		return "", errors.New("product order detail not found")
// 	}
// 	// Update ProductTransaction
// 	transaction.PgStatus = &pgStatus
// 	transaction.PaymentStatus = paymentStatus
// 	rowsAffected, err = s.ProductTransactionRepository.UpdateByModel(transaction)
// 	if err != nil {
// 		return "", err
// 	}
// 	if rowsAffected == 0 {
// 		return "", errors.New("product transaction not found")
// 	}
// 	if orderStatus == "refunded" {
// 		if err := s.ProductVariantRepository.DecreaseTotalSalesByProductOrderID(transaction.ProductOrderID); err != nil {
// 			return "", err
// 		}
// 	}
// 	if strings.HasPrefix(transactionType, "preorder") {
// 		// Update ProductOrderPreorder
// 		if err := s.DB.Save(productPreorder).Error; err != nil {
// 			return "", err
// 		}
// 		// Update ProductOrderInstallmentSchedule
// 		if err := s.DB.Save(installmentSchedule).Error; err != nil {
// 			return "", err
// 		}
// 	}
// 	providerName := transaction.ProductOrder.ProductOrderDetails[0].ProductVariant.ProviderName
// 	isSuccess := paymentStatus == "paid" && orderStatus == "waiting delivery"
// 	// Proses digiflazz
// 	if isSuccess {
// 		if providerName == "digiflazz" {
// 			for i := 0; i < transaction.ProductOrder.TotalQty; i++ {
// 				usecase := usecases.NewCreateOrderToProviderUsecase()
// 				usecase.Execute(transaction.ProductOrder.ID, transaction.ProductOrder.ProductOrderDetails[i].ProductVariant.ProviderSKU, transaction.ProductOrder.ProductOrderDetails[i].ReferenceID, "")
// 			}
// 		}
// 		if providerName == "manual" {

// 			s.TelegramService.SendMessage(fmt.Sprintf(`‼️ [%s] Ada pesanan yang membutuhkan aksi lebih lanjut!
// 			Detail Pesanan:
// 			Order ID: %s
// 			Product : %s
// 			Product Variant : %s
// 			Qty : %d
// 			IGN : %s
// 			Status: %s
// 			Payment Status: %s
// 			PG Status: %s
// 			`, strings.ToUpper(transaction.ProductOrder.ProductOrderDetails[0].ProductVariant.OrderType), transaction.ProductOrderID, transaction.ProductOrder.Product.Name, transaction.ProductOrder.ProductOrderDetails[0].ProductVariant.Name, transaction.ProductOrder.TotalQty, gameUsername, orderStatus, paymentStatus, pgStatus))

// 			_, err := s.FonnteService.FonnteSendMessage(requests.FonnteSendMessageRequest{
// 				Target:  utils.StringPtr(transaction.ProductOrder.WaNumber),
// 				Message: utils.StringPtr(messageWA),
// 			})
// 			if err != nil {
// 				return "", err
// 			}

// 		}
// 	}

// 	utils.WebsocketBroadcast(transaction.ProductOrderID, fiber.Map{
// 		"status_pesanan":    orderStatus,
// 		"status_pembayaran": paymentStatus,
// 	})

// 	message := fmt.Sprintf("Berhasil update status transaksi & membuat transaksi di digiflazz (%s)", paymentStatus)
// 	if orderStatus == "refunded" {
// 		message = fmt.Sprintf("Berhasil update status transaksi & update total sales di product variant (%s)", paymentStatus)
// 	}

// 	return message, nil
// }
