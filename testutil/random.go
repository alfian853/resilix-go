package testutil

import (
	"fmt"
	"math/rand"
	"time"
)

func RandBool() bool {
	return rand.Intn(100)%2 == 1
}

func RandInt(min int, max int) int {
	if max < min {
		panic("max should be greater than min")
	}
	return min + rand.Intn(max-min+ 1)
}

func RandSleep(min int, max int) {
	time.Sleep(time.Duration(RandInt(min, max)) * time.Millisecond)
}

func RandPanicMessage() string {
	return fmt.Sprintf("panic #%d", rand.Intn(1000000000))
}

func RandErrorWithMessage() *IntendedPanic {
	return &IntendedPanic{Message: RandPanicMessage()}
}

func RandStringer() *PanicData {
	return &PanicData{Message: RandPanicMessage()}
}
