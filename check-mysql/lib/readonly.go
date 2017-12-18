package checkmysql

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type readOnlyOpts struct {
	mysqlSetting
}

func checkReadOnly(args []string) *checkers.Checker {
	opts := readOnlyOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Usage = "readonly [OPTIONS]"
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

	rows, res, err := db.Query("SHOW GLOBAL VARIABLES LIKE 'read_only'")
	if err != nil {
		return checkers.Unknown("couldn't execute query")
	}

	idxReadOnly := res.Map("Value")
	readOnlyStatus := rows[0].Str(idxReadOnly)

	if readOnlyStatus == "OFF" {
		return checkers.Critical("MySQL is not read_only")
	}

	return checkers.Ok("MySQL is read_only")
}
