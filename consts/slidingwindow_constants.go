package consts

type SwStrategy string
type SwLock *int32

const (
	// sliding window type
	SwStrategy_CountBased SwStrategy = "count-based"
	SwStrategy_TimeBased  SwStrategy = "time-based"
)

const (
	// ordered by priority asc for @SwLock type
	SwLock_Available   = 0
	SwLock_Clearing    = 1
	SwLock_ClearingAll = 2
)
