package retry

import (
	"github.com/alfian853/resilix-go/consts"
	"github.com/alfian853/resilix-go/util"
)

type RetryExecutor interface {
	util.CheckedExecutor
	GetRetryState() consts.RetryState
}
