// +build windows

package main

import (
	"os/exec"
	"strings"

	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding/japanese"
)

func getServiceState() ([]serviceState, error) {
	b, err := exec.Command("wmic", "service", "get", "Caption,ErrorControl,Name,Started,StartMode,State", "/format:csv").Output()
	if err != nil {
		return nil, err
	}
	b, err = japanese.ShiftJIS.NewDecoder().Bytes(b)
	if err != nil {
		return nil, err
	}
	csv := strings.Replace(string(b), "\r", "", -1)
	csv = strings.TrimLeft(csv, " \t\n")
	var state []serviceState
	err = gocsv.UnmarshalString(csv, &state)
	return state, err
}
