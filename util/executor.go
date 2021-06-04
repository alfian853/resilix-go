package util



type Executor interface {
	Execute(fun func()) bool
	ExecuteSupplier(fun func() interface{}) (bool,interface{})
}

type CheckedExecutor interface {
	ExecuteChecked(fun func() error) (bool, error)
	ExecuteCheckedSupplier(fun func()(interface{}, error)) (bool, interface{}, error)
}
