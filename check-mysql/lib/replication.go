package checkmysql

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type replicationOpts struct {
	mysqlSetting
	Crit int64 `short:"c" long:"critical" default:"250" description:"critical if the seconds behind master is over"`
	Warn int64 `short:"w" long:"warning" default:"200" description:"warning if the seconds behind master is over"`
}

func checkReplication(args []string) *checkers.Checker {
	opts := replicationOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Usage = "replication [OPTIONS]"
	_, err := psr.ParseArgs(args)
	if err != nil {
		os.Exit(1)
	}
	db := newMySQL(opts.mysqlSetting)
	err = db.Connect()
	if err != nil {
		return checkers.Unknown("couldn't connect DB")
	}
	defer db.Close()

	rows, res, err := db.Query("SHOW SLAVE STATUS")
	if err != nil {
		return checkers.Unknown("couldn't execute query")
	}

	if len(rows) == 0 {
		return checkers.Ok("MySQL is not slave")
	}

	idxIoThreadRunning := res.Map("Slave_IO_Running")
	idxSQLThreadRunning := res.Map("Slave_SQL_Running")
	idxSecondsBehindMaster := res.Map("Seconds_Behind_Master")
	ioThreadStatus := rows[0].Str(idxIoThreadRunning)
	sqlThreadStatus := rows[0].Str(idxSQLThreadRunning)
	secondsBehindMaster := rows[0].Int64(idxSecondsBehindMaster)

	if !(ioThreadStatus == "Yes" && sqlThreadStatus == "Yes") {
		return checkers.Critical("MySQL replication has been stopped")
	}

	checkSt := checkers.OK
	msg := fmt.Sprintf("MySQL replication behind master %d seconds", secondsBehindMaster)
	if secondsBehindMaster > opts.Crit {
		checkSt = checkers.CRITICAL
	} else if secondsBehindMaster > opts.Warn {
		checkSt = checkers.WARNING
	}
	return checkers.NewChecker(checkSt, msg)
}
