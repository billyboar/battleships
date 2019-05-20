package helpers

import (
	"math/rand"
	"time"
)

func GenerateRandomInt(limit int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(limit)
}
