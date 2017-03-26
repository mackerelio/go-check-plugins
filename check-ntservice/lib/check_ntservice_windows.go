// +build windows

package checkntservice

import (
	"github.com/StackExchange/wmi"
)

func getServiceState() ([]Win32Service, error) {
	var records []Win32Service
	err := wmi.Query("SELECT * FROM Win32_Service", &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}
