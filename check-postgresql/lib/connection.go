package checkpostgresql

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type connectionOpts struct {
	postgresqlSetting
	Warn int `short:"w" long:"warning" default:"70" description:"warning if the number of connection is over"`
	Crit int `short:"c" long:"critical" default:"90" description:"critical if the number of connection is over"`
}

func checkConnection(args []string) *checkers.Checker {
	opts := connectionOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Usage = "connection [OPTIONS]"
	_, err := psr.ParseArgs(args)
	if err != nil {
		os.Exit(1)
	}

	db, err := sql.Open(opts.getDriverAndDataSourceName())
	if err != nil {
		return checkers.Unknown(err.Error())
	}
	defer db.Close()

	statActivityCount := 0
	err = db.QueryRow("SELECT COUNT(*) AS cnt FROM pg_stat_activity").Scan(&statActivityCount)
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	checkSt := checkers.OK
	msg := fmt.Sprintf("%d connections", statActivityCount)
	if statActivityCount > opts.Crit {
		checkSt = checkers.CRITICAL
	} else if statActivityCount > opts.Warn {
		checkSt = checkers.WARNING
	}
	return checkers.NewChecker(checkSt, msg)
}
