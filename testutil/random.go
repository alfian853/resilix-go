package testutil

import (
	"fmt"
	"math/rand"
)

func RandBool() bool {
	return rand.Intn(100)%2 == 1
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
