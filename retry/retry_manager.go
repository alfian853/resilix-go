package retry

import (
	"resilix-go/consts"
	"resilix-go/util"
)

type RetryExecutor interface {
	util.CheckedExecutor
	GetRetryState() consts.RetryState
}
