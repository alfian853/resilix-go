package consts

type SwStrategy string
type SwLock *int32
const (
	// sliding window type
	COUNT_BASED = "count-based"
	TIME_BASED = "time-based"
)

const (
	// ordered by priority asc for @SwLock type
	SwLock_Available int32   = 0
	SwLock_Clearing    = 1
	SwLock_ClearingAll = 2
)