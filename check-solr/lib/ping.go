package checksolr

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mackerelio/checkers"
)

func checkPing(opts solrOpts) *checkers.Checker {
	uri := opts.createBaseURL() + "/admin/ping?wt=json"
	resp, err := http.Get(uri)
	if err != nil {
		return checkers.Unknown("couldn't get access to " + uri)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var stats map[string]interface{}
	err = dec.Decode(&stats)
	if err != nil {
		return checkers.Unknown("couldn't parse JSON at " + uri)
	}

	status, ok := stats["status"].(string)
	if !ok {
		return checkers.Unknown("couldn't find status in JSON at " + uri)
	}

	checkSt := checkers.OK
	msg := fmt.Sprintf("%s %s", opts.Core, status)
	if status != "OK" {
		checkSt = checkers.CRITICAL
	}
	return checkers.NewChecker(checkSt, msg)
}
