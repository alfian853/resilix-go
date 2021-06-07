package retry

import (
	"github.com/alfian853/resilix-go/consts"
	"github.com/alfian853/resilix-go/executor"
)

type RetryExecutor interface {
	executor.CheckedExecutor
	GetRetryState() consts.RetryState
}
