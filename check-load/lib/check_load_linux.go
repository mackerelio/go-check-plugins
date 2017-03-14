package checkload

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func getloadavg() (loadavgs [3]float64, _ error) {
	contentbytes, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return loadavgs, fmt.Errorf("Failed to load /proc/loadavg: %s", err)
	}
	content := string(contentbytes)
	cols := strings.Split(content, " ")
	for i := 0; i < 3; i++ {
		f, err := strconv.ParseFloat(cols[i], 64)
		if err != nil {
			return loadavgs, fmt.Errorf("Failed to parse loadavg metrics: %s", err)
		}
		loadavgs[i] = f
	}
	return loadavgs, nil
}
