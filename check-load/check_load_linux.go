package main

// // for getloadavg(2)
/*
#include <stdlib.h>
*/
import "C"

func getloadavg() (loadavgs [3]float64, _ error) {
	var cLoadavg [3]C.double
	C.getloadavg(&cLoadavg[0], 3)
	for i, v := range cLoadavg {
		loadavgs[i] = float64(v)
	}

	return loadavgs, nil
}
