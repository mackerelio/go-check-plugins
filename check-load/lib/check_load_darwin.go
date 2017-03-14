package checkload

import (
	"os/exec"
	"strconv"
	"strings"
)

func getloadavg() (loadavgs [3]float64, err error) {
	outputBytes, err := exec.Command("sysctl", "-n", "vm.loadavg").Output()
	if err != nil {
		return loadavgs, err
	}

	output := string(outputBytes)

	// fields will be "{", <loadavg1>, <loadavg5>, <loadavg15>, "}"
	fields := strings.Fields(output)
	if len(fields) != 5 || fields[0] != "{" || fields[len(fields)-1] != "}" {
		return loadavgs, err
	}

	for i := 0; i < 3; i++ {
		loadavg, err := strconv.ParseFloat(fields[i+1], 64)
		if err != nil {
			return loadavgs, err
		}
		loadavgs[i] = loadavg
	}
	return loadavgs, nil
}
