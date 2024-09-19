package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateLoginToken() string {
	rand.Seed(time.Now().UnixNano())

	var digits []int
	for len(digits) < 4 {
		digit := rand.Intn(10)
		if !contains(digits, digit) {
			digits = append(digits, digit)
		}
	}

	token := fmt.Sprintf("%d%d%d%d", digits[0], digits[1], digits[2], digits[3])
	return token
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
