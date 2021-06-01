package retry

type RetryManager interface {
	AcquireAndUpdateRetryPermission() bool
	GetRetryState() RetryState
	GetErrorRate() float32
}
