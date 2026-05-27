package services

import (
	"fmt"
	"sort"
	"time"

	"github.com/bytedance/sonic"
	"github.com/fadilmartias/dilz_code/apps/backend/app/dto"
	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/gofiber/fiber/v3"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type PaymentMethodService struct {
	DB                      *gorm.DB
	PaymentMethodRepository *repositories.PaymentMethodRepository
	Redis                   *redis.Client
}

func NewPaymentMethodService(db *gorm.DB, redis *redis.Client, paymentMethodRepository *repositories.PaymentMethodRepository) *PaymentMethodService {
	return &PaymentMethodService{DB: db, PaymentMethodRepository: paymentMethodRepository, Redis: redis}
}

func (s *PaymentMethodService) GetActivePaymentMethods(
	c fiber.Ctx,
	isCache bool,
	cacheTtl int,
	cacheKey string,
) ([]dto.PaymentMethodDTO, error) {
	ctx := c.Context()

	// 🧩 1. Cek cache dulu
	if isCache {
		if data, err := s.Redis.Get(ctx, config.Key(cacheKey)).Result(); err == nil && data != "" {
			var methods []dto.PaymentMethodDTO
			if err := sonic.Unmarshal([]byte(data), &methods); err == nil {
				// Cache hit 🎯
				return methods, nil
			} else {
				// Log aja error unmarshal, jangan ganggu flow
				fmt.Printf("Redis unmarshal error for key %s: %v\n", cacheKey, err)
			}
		}
	}

	// 🧩 2. Ambil dari repository kalau cache miss
	methods, err := s.PaymentMethodRepository.GetAllActivePaymentMethod(config.LoadAppConfig().Env)
	if err != nil {
		return nil, err
	}

	// 🧩 3. Urutkan berdasarkan kategori
	sort.Slice(methods, func(i, j int) bool {
		return methods[i].Category.Order < methods[j].Category.Order
	})

	// 🧩 4. Mapping ke DTO
	grouped := make([]dto.PaymentMethodDTO, 0, len(methods))
	for _, pm := range methods {
		grouped = append(grouped, dto.PaymentMethodDTO{
			ID:                pm.ID,
			PaymentGatewayID:  pm.PaymentGatewayID,
			CategoryID:        pm.CategoryID,
			Order:             pm.Order,
			Code:              pm.Code,
			InvoiceCode:       pm.InvoiceCode,
			Name:              pm.Name,
			FeeFixed:          pm.FeeFixed,
			FeePercent:        pm.FeePercent,
			PPN:               pm.PPN,
			MinAmount:         pm.MinAmount,
			MaxAmount:         pm.MaxAmount,
			Img:               pm.Img,
			Desc:              pm.Desc,
			IsActive:          pm.IsActive,
			IsReadyProduction: pm.IsReadyProduction,
			IsInternational:   pm.IsInternational,
			CategoryName:      pm.Category.Name,
		})
	}

	// 🧩 5. Simpan ke cache kalau perlu
	if isCache {
		b, err := sonic.Marshal(grouped)
		if err == nil {
			if err := s.Redis.Set(ctx, config.Key(cacheKey), b, time.Duration(cacheTtl)*time.Second); err != nil {
				// Log tapi jangan gagalkan request
				fmt.Printf("Redis set error for key %s: %v\n", cacheKey, err)
			}
		} else {
			fmt.Printf("Redis marshal error for key %s: %v\n", cacheKey, err)
		}
	}

	return grouped, nil
}

// fiber:context-methods migrated
