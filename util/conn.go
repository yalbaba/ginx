package util

import (
	"math/rand"
	"time"
)

func GetConnId() uint32 {
	rand.Seed(time.Now().Unix())
	return rand.Uint32()
}
