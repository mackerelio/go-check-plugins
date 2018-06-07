package checkmysql

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type uptimeOpts struct {
	mysqlSetting
	Crit int64 `short:"c" long:"critical" default:"0" description:"critical if the uptime less than"`
	Warn int64 `short:"w" long:"warning" default:"0" description:"warning if the uptime less than"`
}

func uptime2str(uptime int64) string {
	day := uptime / 86400
	hour := (uptime % 86400) / 3600
	min := ((uptime % 86400) % 3600) / 60
	sec := ((uptime % 86400) % 3600) % 60
	return fmt.Sprintf("%d days, %02d:%02d:%02d", day, hour, min, sec)
}

func checkUptime(args []string) *checkers.Checker {
	opts := uptimeOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Usage = "uptime [OPTIONS]"
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

	rows, res, err := db.Query("SHOW GLOBAL STATUS LIKE 'Uptime'")
	if err != nil {
		return checkers.Unknown("couldn't execute query")
	}

	idxValue := res.Map("Value")
	upTime := rows[0].Int64(idxValue)

	checkSt := checkers.OK
	msg := fmt.Sprintf("up %s", uptime2str(upTime))
	if opts.Crit > 0 && upTime < opts.Crit {
		checkSt = checkers.CRITICAL
	} else if opts.Warn > 0 && upTime < opts.Warn {
		checkSt = checkers.WARNING
	}
	return checkers.NewChecker(checkSt, msg)
}
