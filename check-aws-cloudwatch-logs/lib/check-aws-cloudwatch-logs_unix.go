package checkawscloudwatchlogs

import (
	"syscall"
)

func init() {
	defaultSignal = syscall.SIGTERM
}
