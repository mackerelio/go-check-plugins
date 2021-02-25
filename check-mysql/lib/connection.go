package checkmysql

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type connectionOpts struct {
	mysqlSetting
	Crit int64 `short:"c" long:"critical" default:"250" description:"critical if the number of connection is over"`
	Warn int64 `short:"w" long:"warning" default:"200" description:"warning if the number of connection is over"`
}

func checkConnection(args []string) *checkers.Checker {
	opts := connectionOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Usage = "connection [OPTIONS]"
	_, err := psr.ParseArgs(args)
	if err != nil {
		os.Exit(1)
	}
	db, err := newDB(opts.mysqlSetting)
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Couldn't open DB: %s", err))
	}
	defer db.Close()

	var (
		variableName     string
		threadsConnected int64
	)
	err = db.QueryRow("SHOW GLOBAL STATUS LIKE 'Threads_Connected'").Scan(&variableName, &threadsConnected)
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Couldn't get 'Threads_Connected' status: %s", err))
	}

	checkSt := checkers.OK
	msg := fmt.Sprintf("%d connections", threadsConnected)
	if threadsConnected > opts.Crit {
		checkSt = checkers.CRITICAL
	} else if threadsConnected > opts.Warn {
		checkSt = checkers.WARNING
	}
	return checkers.NewChecker(checkSt, msg)
}
