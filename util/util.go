package util

import (
	"fmt"
	"strconv"
)

func IntToFloat(num int) float64 {
	return float64(num) / 100
}

func PriceToCents(num string) int {
	price, _ := strconv.Atoi(num)

	return price
}

func CentsToEuros(cents int) string {
	euros := float64(cents) / 100.0
	return fmt.Sprintf("%.2f", euros)
}
