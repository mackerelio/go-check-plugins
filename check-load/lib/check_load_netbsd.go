package checkload

import (
  "github.com/mackerelio/go-osstat/loadavg"
)

func getloadavg() (loadavgs [3]float64, err error) {
  output, err := loadavg.Get()
  loadavgs = [3]float64{output.Loadavg1,output.Loadavg5,output.Loadavg15}

  return loadavgs, nil
}
