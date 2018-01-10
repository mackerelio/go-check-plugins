package checkmysql

import (
	"fmt"
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
	psr.Usage = "readonly [OPTIONS] [ON|OFF]"
	args, err := psr.ParseArgs(args)
	if err != nil {
		os.Exit(1)
	}
	if len(args) != 1 {
		fmt.Println("wrong number of arguments")
		os.Exit(1)
	}
	argStatus := args[0]

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

	if readOnlyStatus != argStatus {
		msg := fmt.Sprintf("the expected value of read_only is different. readOnlyStatus:%s", readOnlyStatus)
		return checkers.Critical(msg)
	}

	return checkers.Ok("read_only is the expected value")
}
