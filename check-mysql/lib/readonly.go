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

	db, err := newDB(opts.mysqlSetting)
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Couldn't open DB: %s", err))
	}
	defer db.Close()

	var (
		variableName   string
		readOnlyStatus string
	)
	err = db.QueryRow("SHOW GLOBAL VARIABLES LIKE 'read_only'").Scan(&variableName, &readOnlyStatus)
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Couldn't get 'read_only' variable status: %s", err))
	}

	if readOnlyStatus != argStatus {
		msg := fmt.Sprintf("the expected value of read_only is different. readOnlyStatus:%s", readOnlyStatus)
		return checkers.Critical(msg)
	}

	return checkers.Ok(fmt.Sprintf("read_only is the expected value (%s)", argStatus))
}
