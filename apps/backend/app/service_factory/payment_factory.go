// // service_factory/payment.go
package service_factory

// import (
// 	"errors"
// 	"strings"

// 	"github.com/fadilmartias/dilz_code/apps/backend/app/dto"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
// )

// // bikin interface
// type Midtrans interface {
// 	CreateCoreAPIPayment(body dto.CreateTransactionParams) (*models.ProductTransaction, error)
// }

// type PaymentFactory struct {
// 	MidtransService Midtrans
// }

// func NewPaymentFactory(mid Midtrans) *PaymentFactory {
// 	return &PaymentFactory{
// 		MidtransService: mid,
// 	}
// }

// func (s *PaymentFactory) CreateTransactionByGateway(body dto.CreateTransactionParams) (*models.ProductTransaction, error) {
// 	switch strings.ToLower(body.PaymentMethod.PaymentGateway.Name) {
// 	case "midtrans":
// 		return s.MidtransService.CreateCoreAPIPayment(body)
// 	default:
// 		return nil, errors.New("unsupported payment gateway: " + body.PaymentMethod.PaymentGateway.Name)
// 	}
// }
