package utils

import (
	"math"
)

func ValidateDiscount(price float64, qty int, fixedDiscount, discountPercent, maxDiscount, minOrderAmount float64) float64 {
	totalPrice := price * float64(qty)
	if totalPrice < minOrderAmount {
		return 0
	}

	percentDiscount := math.Floor(totalPrice * discountPercent)
	totalDiscount := fixedDiscount + percentDiscount

	// kalau maxDiscount == 0 artinya unlimited
	if maxDiscount > 0 && totalDiscount > maxDiscount {
		totalDiscount = maxDiscount
	}

	return totalDiscount
}

func CalculateDiscountedPrice(price float64, qty int, fixedDiscount, discountPercent, maxDiscount, minOrderAmount float64) (float64, float64) {
	totalPrice := price * float64(qty)

	discount := ValidateDiscount(price, qty, fixedDiscount, discountPercent, maxDiscount, minOrderAmount)

	finalPrice := totalPrice - discount
	if finalPrice < 0 {
		return 0, discount
	}
	return finalPrice, discount
}

func CalculateServiceFee(basePrice float64, totalPrice float64) float64 {
	fee := totalPrice - basePrice
	if fee < 0 {
		return 0
	}
	return fee
}

func CalculateTotalPrice(
	discountedPrice float64,
	marginFixed float64,
	marginPercent float64,
	fixedFee float64,
	feePercent float64,
	ppnPercent float64,
) (totalPrice, totalFee, pgFee, profit float64) {

	profit = math.Ceil((discountedPrice * marginPercent) + marginFixed)

	// Step 1: hitung gross dengan rumus pembalik fee PG
	// feeMultiplier := 1 + ppnPercent
	// denominator := 1 - (feePercent * feeMultiplier)
	// numerator := ((discountedPrice * (1 + marginPercent)) + marginFixed) + (fixedFee * feeMultiplier)
	feeMultiplier := 1 + ppnPercent
	denominator := 1 - (feePercent * feeMultiplier)
	numerator := discountedPrice + profit + (fixedFee * feeMultiplier)

	gross := numerator / denominator
	totalPrice = math.Ceil(gross)

	// ceiledProfit := math.Ceil(profit)

	// Step 2: hitung PG Fee yang benar
	pgFee = totalPrice - discountedPrice - profit

	// Step 3: total fee (selisih antara harga jual dan harga bersih)
	totalFee = pgFee + profit
	if totalFee < 0 {
		totalFee = 0
	}

	return totalPrice, totalFee, pgFee, profit
}
