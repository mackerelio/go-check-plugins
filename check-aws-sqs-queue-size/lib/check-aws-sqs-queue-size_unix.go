package checkawssqsqueuesize

import (
	"syscall"
)

func init() {
	defaultSignal = syscall.SIGTERM
}
