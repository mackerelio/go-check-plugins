// +build !windows

package checkntservice

import (
	"syscall"
)

func getServiceState() ([]Win32Service, error) {
	return nil, syscall.ENOSYS
}
