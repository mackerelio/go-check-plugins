//go:build windows

package checkdns

import (
	"golang.org/x/sys/windows"
	"os"
	"syscall"
	"unsafe"
)

// https://go.dev/src/net/interface_windows.go
func adapterAddress() (string, error) {
	var b []byte
	l := uint32(15000) // recommended initial size
	for {
		b = make([]byte, l)
		err := windows.GetAdaptersAddresses(syscall.AF_UNSPEC, windows.GAA_FLAG_INCLUDE_PREFIX, 0, (*windows.IpAdapterAddresses)(unsafe.Pointer(&b[0])), &l)
		if err == nil {
			if l == 0 {
				return "", nil
			}
			break
		}
		if err.(syscall.Errno) != syscall.ERROR_BUFFER_OVERFLOW {
			return "", os.NewSyscallError("getadaptersaddresses", err)
		}
		if l <= uint32(len(b)) {
			return "", os.NewSyscallError("getadaptersaddresses", err)
		}
	}
	var aas []*windows.IpAdapterAddresses
	for aa := (*windows.IpAdapterAddresses)(unsafe.Pointer(&b[0])); aa != nil; aa = aa.Next {
		aas = append(aas, aa)
	}
	return aas[0].FirstDnsServerAddress.Address.IP().String(), nil
}
