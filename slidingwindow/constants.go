package slidingwindow

type SwLock *int32

const (
	// ordered by priority asc for @SwLock type
	SwLock_Available   = 0
	SwLock_Clearing    = 1
	SwLock_ClearingAll = 2
)
