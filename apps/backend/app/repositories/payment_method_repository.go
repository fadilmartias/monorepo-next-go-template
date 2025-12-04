package repositories

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"gorm.io/gorm"
)

type PaymentMethodRepository struct {
	BaseRepository[models.PaymentMethod]
	DB *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) *PaymentMethodRepository {
	return &PaymentMethodRepository{BaseRepository: NewBaseRepository[models.PaymentMethod](db), DB: db}
}

func (r *PaymentMethodRepository) WithTx(tx *gorm.DB) *PaymentMethodRepository {
	return &PaymentMethodRepository{BaseRepository: r.BaseRepository.WithTx(tx), DB: tx}
}

func (r *PaymentMethodRepository) GetAllActivePaymentMethod(env string) ([]models.PaymentMethod, error) {
	var methods []models.PaymentMethod

	query := r.DB.
		Preload("PaymentGateway").
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Order("categories.order ASC")
		}).
		Where("is_active = ?", true)

	if env == "production" {
		query = query.Where("is_ready_production = ?", true)
	}

	if err := query.Find(&methods).Error; err != nil {
		return nil, err
	}

	return methods, nil
}

func (r *PaymentMethodRepository) GetActiveMidtransPaymentMethod() (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	if err := r.DB.Where("payment_gateway_id = ?", "PG1").Where("is_active = ?", true).First(&paymentMethod).Error; err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

func (r *PaymentMethodRepository) GetActiveIPaymuPaymentMethod() (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	if err := r.DB.Where("payment_gateway_id = ?", "PG2").Where("is_active = ?", true).First(&paymentMethod).Error; err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

func (r *PaymentMethodRepository) FindActive(id string) (*models.PaymentMethod, error) {
	var pm models.PaymentMethod
	query := r.DB.Where("id = ?", id).
		Preload("PaymentGateway", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name").Where("is_active = ?", true)
		}).Where("is_active = ?", true)

	if config.LoadAppConfig().Env == "production" {
		query = query.Where("is_ready_production = ?", true)
	}
	if err := query.First(&pm).Error; err != nil {
		return nil, err
	}
	return &pm, nil
}
