package statehandler



type Executor interface {
	Execute(fun func()) bool
	ExecuteWithReturn(fun func() interface{}) (bool,interface{})
}
