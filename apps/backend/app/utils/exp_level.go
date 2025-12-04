package utils

// Hitung EXP dari transaksi topup
func CalculateExpFromTopup(amount float64) int {
	return int(amount / 100) // 1 EXP per Rp100
}

// Rumus total exp dibutuhkan utk naik level berikut
func ExpToNextLevel(level int) int {
	base := 200
	return base * level * level
}
