// +build !windows

package main

import (
	"syscall"
)

func getServiceState() ([]serviceState, error) {
	return nil, syscall.ENOSYS
}
