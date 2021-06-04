package statehandler



type Executor interface {
	Execute(fun func()) bool
	ExecuteSupplier(fun func() interface{}) (bool,interface{})
}
