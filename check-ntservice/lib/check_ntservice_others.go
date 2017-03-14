// +build !windows

package checkntservice

import (
	"syscall"
)

func getServiceState() ([]serviceState, error) {
	return nil, syscall.ENOSYS
}
