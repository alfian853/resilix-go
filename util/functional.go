package util


func CheckedRunnable()func() error {
	return func() error {
		return nil
	}
}
func DoNothingRunnable() func() {
	return func() {}
}

func PanicCheckedRunnable(any interface{}) func() error {
	return func() error {
		panic(any)
	}
}

func ErrorCheckedRunnable(err error) func() error {
	return func() error {
		return err
	}
}

func PanicCheckedSupplier(any interface{}) func() (interface{},error) {
	return func() (interface{},error) {
		panic(any)
	}
}

func ErrorCheckedSupplier(err error) func() (interface{},error) {
	return func() (interface{},error) {
		return nil, err
	}
}

func TrueSupplier() func() bool {
	return func() bool {
		return true
	}
}