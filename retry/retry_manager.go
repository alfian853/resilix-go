package retry

import (
	"resilix-go/consts"
)

type RetryManager interface {
	AcquireAndUpdateRetryPermission() bool
	GetRetryState() consts.RetryState
	GetErrorRate() float32
}
