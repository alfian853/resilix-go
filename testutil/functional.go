package testutil

func CheckedRunnable() func() error {
	return func() error {
		return nil
	}
}

func ErrorCheckedRunnable(err error) func() error {
	return func() error {
		return err
	}
}

func SpyRunnable(isRun *bool) func() {
	return func() {
		*isRun = true
	}
}

func SpyCheckedRunnable(isRun *bool) func() error {
	return func() error {
		*isRun = true
		return nil
	}
}

func PanicRunnable(any interface{}) func() {
	return func() {
		panic(any)
	}
}

func PanicCheckedRunnable(any interface{}) func() error {
	return func() error {
		panic(any)
	}
}

func PanicSupplier(any interface{}) func() interface{} {
	return func() interface{} {
		panic(any)
	}
}

func PanicCheckedSupplier(any interface{}) func() (interface{}, error) {
	return func() (interface{}, error) {
		panic(any)
	}
}

func ErrorCheckedSupplier(err error) func() (interface{}, error) {
	return func() (interface{}, error) {
		return nil, err
	}
}

func Supplier(any interface{}) func() interface{} {
	return func() interface{} {
		return any
	}
}

func CheckedSupplier(any interface{}) func() (interface{}, error) {
	return func() (interface{}, error) {
		return any, nil
	}
}

func TrueSupplier() func() interface{} {
	return func() interface{} {
		return true
	}
}

func TrueCheckedSupplier() func() (interface{}, error) {
	return func() (interface{}, error) {
		return true, nil
	}
}
