package checkelasticsearch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type healthStat struct {
	ClusterName string `json:"cluster_name"`
	Status      string `json:"status"`
}

var opts struct {
	Scheme string `short:"s" long:"scheme" default:"http" description:"Elasticsearch scheme"`
	Host   string `short:"H" long:"host" default:"localhost" description:"Elasticsearch host"`
	Port   int64  `short:"p" long:"port" default:"9200" description:"Elasticsearch port"`
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Elasticsearch"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	client := &http.Client{}
	url := fmt.Sprintf("%s://%s:%d/_cluster/health", opts.Scheme, opts.Host, opts.Port)

	stTime := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return checkers.Critical(err.Error())
	}
	elapsed := time.Since(stTime)
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var health healthStat
	dec.Decode(&health)

	checkSt := checkers.UNKNOWN
	switch health.Status {
	case "green":
		checkSt = checkers.OK
	case "yellow":
		checkSt = checkers.WARNING
	case "red":
		checkSt = checkers.CRITICAL
	default:
		checkSt = checkers.UNKNOWN
	}

	msg := fmt.Sprintf("%s (cluster: %s) - %f second respons time",
		health.Status, health.ClusterName, elapsed.Seconds())

	return checkers.NewChecker(checkSt, msg)
}
