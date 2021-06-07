package testutil

import "fmt"

type IntendedPanic struct {
	error
	Message string
}

type PanicData struct {
	fmt.Stringer
	fmt.GoStringer
	Message string
}

func (err *IntendedPanic) Error() string {
	return err.Message
}

func (data *PanicData) String() string {
	return data.Message
}

func (data *PanicData) GoString() string {
	return data.Message
}
