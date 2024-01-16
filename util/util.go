package util

import "strconv"

func IntToFloat(num int) float64 {
	return float64(num) / 100
}

func PriceToCents(num string) int {
	price, _ := strconv.Atoi(num)

	return price * 100
}
