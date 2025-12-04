export const validateDiscount = (
  price: number,
  qty: number,
  fixedDiscount: number,
  discountPercent: number,
  maxDiscount: number,
  minOrderAmount: number
): number => {
  const totalPrice = price * qty;
  if (totalPrice < minOrderAmount) return 0;

  const percentDiscount = Math.floor(totalPrice * discountPercent);
  let totalDiscount = fixedDiscount + percentDiscount;

  if (maxDiscount > 0 && totalDiscount > maxDiscount) {
    totalDiscount = maxDiscount;
  }

  return totalDiscount;
};


export const calculateDiscountedPrice = (
  price: number,
  qty: number,
  fixedDiscount: number,
  discountPercent: number,
  maxDiscount: number,
  minOrderAmount: number
): number => {
  const totalPrice = price * qty;
  const discount = validateDiscount(
    price,
    qty,
    fixedDiscount,
    discountPercent,
    maxDiscount,
    minOrderAmount
  );

  const finalPrice = totalPrice - discount;
  return Math.max(finalPrice, 0);
};


export const calculateServiceFee = (discountedPrice: number, totalPrice: number) => {
  return Math.ceil(totalPrice - discountedPrice);
};

export const calculateTotalPrice = (
  discountedPrice: number, // target bersih yang kamu mau terima
  marginFixed: number, // biaya flat, contoh: 4000
  marginPercent: number, // contoh: 0.011 untuk 1.1%
  fixedFee: number, // biaya flat, contoh: 4000
  feePercent: number, // contoh: 0.011 untuk 1.1%
  ppnPercent: number // contoh: 0.11 untuk 11%
): number => {
  const feeMultiplier = 1 + ppnPercent;
  const denominator = 1 - feePercent * feeMultiplier;
  const numerator = ((discountedPrice * (1 + marginPercent)) + marginFixed) + fixedFee * feeMultiplier;

  const gross = numerator / denominator;

  return Math.ceil(gross);
};

export const ceilNumber = (num: number, fractionDigits = 0): number =>
  Math.ceil(num * 10 ** fractionDigits) / 10 ** fractionDigits;

export const floorNumber = (num: number, fractionDigits = 0): number =>
  Math.floor(num * 10 ** fractionDigits) / 10 ** fractionDigits;

export const roundNumber = (num: number, fractionDigits = 0): number =>
  Math.round(num * 10 ** fractionDigits) / 10 ** fractionDigits;

// Fungsi untuk menghitung waktu baca
export function calculateReadingTime(content: string): number {
  const wordsPerMinute = 200;
  const text = content.replace(/<[^>]*>/g, "");
  const wordCount = text.split(/\s+/).length;
  return Math.ceil(wordCount / wordsPerMinute);
}



