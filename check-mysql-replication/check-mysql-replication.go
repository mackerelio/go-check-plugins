package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

var opts struct {
	Host string `short:"h" long:"host" default:"localhost" description:"Hostname"`
	Port string `short:"p" long:"port" default:"3306" description:"Port"`
	User string `short:"u" long:"user" default:"root" description:"Username"`
	Pass string `short:"P" long:"password" default:"" description:"Password"`
	Crit int64  `short:"c" long:"critical" default:"250" description:"critical if the seconds behind master is over"`
	Warn int64  `short:"w" long:"warning" default:"200" description:"warning if the seconds behind master is over"`
}

func main() {
	ckr := run(os.Args[1:])
	ckr.Name = "MySQL Replication"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}
	target := fmt.Sprintf("%s:%s", opts.Host, opts.Port)
	db := mysql.New("tcp", "", target, opts.User, opts.Pass, "")
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

	if ioThreadStatus == "No" || sqlThreadStatus == "No" {
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
