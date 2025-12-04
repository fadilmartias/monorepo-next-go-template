// // services/checkout_service.go
package services

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"math"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"

// 	"github.com/fadilmartias/dilz_code/apps/backend/app/dto"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/service_factory"
// 	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
// 	"github.com/fadilmartias/dilz_code/apps/backend/config"
// )

// type CheckoutService struct {
// 	db                  *gorm.DB
// 	redis               *config.RedisClient
// 	paymentService      *PaymentService
// 	telegramService     *TelegramService
// 	settingRepo         *repositories.SettingRepository
// 	paymentFactory      *service_factory.PaymentFactory

// 	productRepo       *repositories.ProductRepository
// 	paymentMethodRepo *repositories.PaymentMethodRepository
// }

// func NewCheckoutService(db *gorm.DB, redis *config.RedisClient, paymentService *PaymentService, telegramService *TelegramService, digiflazzService *DigiflazzService, voucherService *VoucherService, promotionService *PromotionService, productOrderService *ProductOrderService, productRepo *repositories.ProductRepository, paymentMethodRepo *repositories.PaymentMethodRepository, settingRepo *repositories.SettingRepository, paymentFactory *service_factory.PaymentFactory) *CheckoutService {
// 	return &CheckoutService{
// 		db:                  db,
// 		redis:               redis,
// 		paymentService:      paymentService,
// 		telegramService:     telegramService,
// 		digiflazzService:    digiflazzService,
// 		voucherService:      voucherService,
// 		promotionService:    promotionService,
// 		productOrderService: productOrderService,
// 		productRepo:         productRepo,
// 		paymentMethodRepo:   paymentMethodRepo,
// 		settingRepo:         settingRepo,
// 		paymentFactory:      paymentFactory,
// 	}
// }

// func (s *CheckoutService) Checkout(ctx context.Context, userId *string, tenantId string, input requests.CheckoutInput) (map[string]any, error) {
// 	// --- STEP 1: ambil price list dari redis ---
// 	// priceListBytes, err := s.redis.GetBytes(ctx, "digiflazz:price-list:prepaid")
// 	// if err == s.redis.Nil() {

// 	// 	// kalau kosong, panggil updater
// 	// 	s.digiflazzService.UpdatePrepaidProductDigiflazz()
// 	// 	// coba lagi
// 	// 	priceListBytes, err = s.redis.GetBytes(ctx, "digiflazz:price-list:prepaid")
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// } else if err != nil {
// 	// 	return nil, err
// 	// }

// 	// var priceList map[string]any
// 	// if err := sonic.Unmarshal(priceListBytes, &priceList); err != nil {
// 	// 	return nil, err
// 	// }

// 	// --- STEP 2: ambil produk + variant ---
// 	product, err := s.productRepo.FindWithVariant(input.ProductID, input.ProductVariantID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	productVariant := product.ProductVariants[0]

// 	orderType := productVariant.OrderType
// 	costPrice := productVariant.CostPrice

// 	itemPricePlusMargin := productVariant.Price
// 	baseFullPrice := costPrice * float64(input.Qty)

// 	// --- STEP 3: cek saldo Digiflazz ---
// 	balanceDigiflazz, err := s.digiflazzService.GetDigiflazzBalance()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if productVariant.ProviderName == "digiflazz" && (baseFullPrice > float64(balanceDigiflazz)) {
// 		message := fmt.Sprintf("❌ [ERROR] Saldo digiflazz tidak mencukupi. Saldo tersisa %d. Segera topup di %s",
// 			balanceDigiflazz, "https://member.digiflazz.com/buyer-area/topup")
// 		s.telegramService.SendMessage(message)
// 		return nil, errors.New("saldo tidak mencukupi")
// 	}

// 	// --- STEP 4: cek product readiness ---
// 	// sku := product.ProductVariants[0].ProviderSKU
// 	// isReady, ok := priceList[sku+"_is_ready"].(bool)
// 	// if !ok || !isReady {
// 	// 	message := fmt.Sprintf("❌ [ERROR] Product with SKU %s is not ready, detail:\n%s",
// 	// 		sku, utils.Sdump(product))
// 	// 	s.telegramService.SendMessage(message)
// 	// 	return nil, errors.New("product not ready")
// 	// }

// 	// --- STEP 5: ambil payment method ---
// 	paymentMethod, err := s.paymentMethodRepo.FindActive(input.PaymentMethod.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// --- STEP 6: validasi dan proses semua promotion (termasuk voucher) ---
// 	var (
// 		totalDiscAmount     float64
// 		appliedPromotions   []models.PivProductOrderPromotion
// 		basePriceAfterPromo float64 = baseFullPrice
// 	)

// 	orderID := models.GenerateID(7)

// 	for _, promo := range input.Promotions {
// 		if promo.Type == "voucher" || promo.ID == "" {
// 			continue
// 		}
// 		isAllowed, item, err := s.promotionService.ValidatePromotionItem(
// 			promo.ID,
// 			input.ProductID,
// 			input.ProductVariantID,
// 			input.PaymentMethod.ID,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if !isAllowed {
// 			return nil, fmt.Errorf("promo %s tidak berlaku untuk item ini", promo.ID)
// 		}

// 		var discFixed, discPercent, promoDiscount float64

// 		// --- PROMO BIASA ---
// 		discFixed = item.DiscFixed
// 		discPercent = item.DiscPercent

// 		_, promoDiscount = utils.CalculateDiscountedPrice(
// 			costPrice, // bukan harga dp
// 			input.Qty,
// 			discFixed,
// 			discPercent,
// 			0,
// 			0,
// 		)

// 		basePriceAfterPromo = math.Max(basePriceAfterPromo-promoDiscount, 0)
// 		totalDiscAmount += promoDiscount

// 		appliedPromotions = append(appliedPromotions, models.PivProductOrderPromotion{
// 			ProductOrderID: orderID,
// 			PromotionID:    promo.ID,
// 			Type:           promo.Type,
// 			DiscFixed:      discFixed,
// 			DiscPercent:    discPercent,
// 			TotalDiscount:  promoDiscount,
// 		})
// 	}

// 	// --- STEP 7: Hitung voucher (jika ada) ---
// 	for _, promo := range input.Promotions {
// 		if promo.Type != "voucher" || !promo.IsRedeemed || promo.ID == "" {
// 			continue
// 		}

// 		vRes, err := s.voucherService.ValidateVoucher(
// 			promo.Code,
// 			userId,
// 			input.ProductID,
// 			input.ProductVariantID,
// 			input.Qty,
// 			input.PaymentMethod.ID,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		discFixed := vRes.Voucher.DiscFixed
// 		discPercent := vRes.Voucher.DiscPercent

// 		// Voucher dihitung dari basePriceAfterPromo
// 		_, voucherDiscount := utils.CalculateDiscountedPrice(
// 			basePriceAfterPromo, // harga satuan terbaru
// 			1,
// 			discFixed,
// 			discPercent,
// 			vRes.Voucher.MaxDiscount,
// 			vRes.Voucher.MinOrderAmount,
// 		)

// 		basePriceAfterPromo = math.Max(basePriceAfterPromo-voucherDiscount, 0)
// 		totalDiscAmount += voucherDiscount

// 		appliedPromotions = append(appliedPromotions, models.PivProductOrderPromotion{
// 			ProductOrderID: orderID,
// 			PromotionID:    promo.ID,
// 			VoucherCode:    models.NewNullString(promo.Code),
// 			Type:           promo.Type,
// 			DiscFixed:      discFixed,
// 			DiscPercent:    discPercent,
// 			TotalDiscount:  voucherDiscount,
// 		})
// 	}

// 	var dpAmount float64
// 	preorderRule := &models.ProductVariantPreorderRule{}
// 	if orderType == "preorder" {
// 		err := s.db.Where("id = ?", input.PreorderRuleID).First(preorderRule).Error
// 		if err != nil {
// 			return nil, err
// 		}
// 		dpAmount = basePriceAfterPromo * preorderRule.MinDPPercent
// 	}

// 	// --- STEP 8: Hitung harga akhir ---
// 	totalPriceFull, totalFeeFull, pgFeeFull, profitFull := utils.CalculateTotalPrice(
// 		basePriceAfterPromo,
// 		product.ProductVariants[0].MarginFixed,
// 		product.ProductVariants[0].MarginPercent,
// 		paymentMethod.FeeFixed,
// 		paymentMethod.FeePercent,
// 		paymentMethod.PPN,
// 	)

// 	var totalPriceDP, totalFeeDP, pgFeeDP, profitDP, itemPriceDP float64
// 	if orderType == "preorder" {
// 		totalPriceDP, totalFeeDP, pgFeeDP, profitDP = utils.CalculateTotalPrice(
// 			dpAmount,
// 			product.ProductVariants[0].MarginFixed,
// 			product.ProductVariants[0].MarginPercent,
// 			paymentMethod.FeeFixed,
// 			paymentMethod.FeePercent,
// 			paymentMethod.PPN,
// 		)
// 		itemPriceDP = dpAmount
// 		pgFeeFull = 0
// 	} else {
// 		totalPriceDP = totalPriceFull
// 		totalFeeDP = totalFeeFull
// 		pgFeeDP = pgFeeFull
// 		profitDP = profitFull
// 		itemPriceDP = basePriceAfterPromo
// 	}

// 	// --- STEP 9: buat transaksi ke gateway ---
// 	refId := paymentMethod.Code + "-" + uuid.New().String()
// 	transactionID := utils.GenerateShortID(7)

// 	createTransactionParams := dto.CreateTransactionParams{
// 		ID:                  transactionID,
// 		Product:             product,
// 		ProductVariant:      product.ProductVariants[0],
// 		Promotions:          appliedPromotions,
// 		Qty:                 input.Qty,
// 		PaymentMethod:       *paymentMethod,
// 		TotalFee:            totalFeeDP, // pg_fee + profit
// 		PGFee:               pgFeeDP,
// 		Profit:              profitDP,
// 		TotalPrice:          totalPriceDP,        // final_price (amount)
// 		ItemPrice:           itemPriceDP,         // cost_price
// 		ItemPricePlusMargin: itemPricePlusMargin, // price
// 		TotalDiscount:       &totalDiscAmount,
// 		OrderId:             orderID,
// 		RefId:               refId,
// 		NoWA:                input.WaNumber,
// 		BasePrice:           basePriceAfterPromo,            // (cost_price - discount) * qty
// 		PriceAfterMargin:    basePriceAfterPromo + profitDP, // (cost_price - discount) * qty + profit
// 	}

// 	transaction, err := s.paymentFactory.CreateTransactionByGateway(createTransactionParams)
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal membuat transaksi: %w", err)
// 	}

// 	params := CreateOrderParams{
// 		TenantID:            tenantId,
// 		UserID:              userId,
// 		Product:             product,
// 		Promotions:          appliedPromotions,
// 		ProductTransaction:  *transaction,
// 		CheckoutInput:       input,
// 		BasePrice:           baseFullPrice,       // cost_price * qty
// 		BasePriceAfterPromo: basePriceAfterPromo, // (cost_price - discount) * qty
// 		FixedDiscount:       totalDiscAmount,     // total discount
// 		PercentDiscount:     0,
// 		TotalPrice:          totalPriceFull, // final_price
// 		TotalFee:            totalFeeFull,   // pg_fee + profit
// 		PGFee:               pgFeeFull,
// 		Profit:              profitFull,
// 		PriceAfterMargin:    basePriceAfterPromo + profitFull, // (cost_price - discount) * qty + profit
// 		OrderType:           orderType,
// 		PreorderRule:        *preorderRule,
// 	}

// 	if err := s.productOrderService.CreateOrder(params); err != nil {
// 		return nil, err
// 	}

// 	return map[string]any{
// 		"order_id":       orderID,
// 		"transaction_id": transactionID,
// 		"checkout_url":   transaction.CheckoutURL,
// 	}, nil
// }
