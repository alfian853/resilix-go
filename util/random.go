package util

import "math/rand"

func RandBool() bool {
	return rand.Intn(100) % 2 == 1
}