package helpers

import (
	"math/rand"
	"strconv"
	"time"
)

// GenerateVerificationCode generates a random 6 digit number
func GenerateVerificationCode() string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random number between 100000 and 999999 (6 digits) and convert it to string
	return strconv.Itoa(random.Intn(999999) + 100000)
}
