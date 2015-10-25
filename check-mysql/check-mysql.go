package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

var cOpts struct {
	Host string `short:"h" long:"host" default:"localhost" description:"Hostname"`
	Port string `short:"p" long:"port" default:"3306" description:"Port"`
	User string `short:"u" long:"user" default:"root" description:"Username"`
	Pass string `short:"P" long:"password" default:"" description:"Password"`
	Crit int64  `short:"c" long:"critical" default:"250" description:"critical if the number of connection is over"`
	Warn int64  `short:"w" long:"warning" default:"200" description:"warning if the number of connection is over"`
}

var rOpts struct {
	Host string `short:"h" long:"host" default:"localhost" description:"Hostname"`
	Port string `short:"p" long:"port" default:"3306" description:"Port"`
	User string `short:"u" long:"user" default:"root" description:"Username"`
	Pass string `short:"P" long:"password" default:"" description:"Password"`
	Crit int64  `short:"c" long:"critical" default:"250" description:"critical if the seconds behind master is over"`
	Warn int64  `short:"w" long:"warning" default:"200" description:"warning if the seconds behind master is over"`
}

var commands = map[string](func([]string) *checkers.Checker){
	"replication": checkReplication,
	"connection":  checkConnection,
}

func separateSub(argv []string) (string, []string) {
	if len(argv) == 0 || strings.HasPrefix(argv[0], "-") {
		return "", argv
	}
	return argv[0], argv[1:]
}

func main() {
	subCmd, argv := separateSub(os.Args[1:])
	fn, ok := commands[subCmd]
	if !ok {
		fmt.Println(`Usage:
  check-mysql [subcommand] [OPTIONS]

Subcommand:
  connection
  replication`)
		os.Exit(1)
	}
	ckr := fn(argv)
	ckr.Name = fmt.Sprintf("MySQL %s", strings.ToUpper(string(subCmd[0]))+subCmd[1:])
	ckr.Exit()
}

func checkConnection(args []string) *checkers.Checker {
	psr := flags.NewParser(&cOpts, flags.Default)
	psr.Usage = "connection [OPTIONS]"
	_, err := psr.ParseArgs(args)
	if err != nil {
		os.Exit(1)
	}
	target := fmt.Sprintf("%s:%s", cOpts.Host, cOpts.Port)
	db := mysql.New("tcp", "", target, cOpts.User, cOpts.Pass, "")
	err = db.Connect()
	if err != nil {
		return checkers.Unknown("couldn't connect DB")
	}
	defer db.Close()

	rows, res, err := db.Query("SHOW GLOBAL STATUS LIKE 'Threads_Connected'")
	if err != nil {
		return checkers.Unknown("couldn't execute query")
	}

	idxValue := res.Map("Value")
	threadsConnected := rows[0].Int64(idxValue)

	checkSt := checkers.OK
	msg := fmt.Sprintf("%d connections", threadsConnected)
	if threadsConnected > cOpts.Crit {
		checkSt = checkers.CRITICAL
	} else if threadsConnected > cOpts.Warn {
		checkSt = checkers.WARNING
	}
	return checkers.NewChecker(checkSt, msg)
}

func checkReplication(args []string) *checkers.Checker {
	psr := flags.NewParser(&rOpts, flags.Default)
	psr.Usage = "replication [OPTIONS]"
	_, err := psr.ParseArgs(args)
	if err != nil {
		os.Exit(1)
	}
	target := fmt.Sprintf("%s:%s", rOpts.Host, rOpts.Port)
	db := mysql.New("tcp", "", target, rOpts.User, rOpts.Pass, "")
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
	if secondsBehindMaster > rOpts.Crit {
		checkSt = checkers.CRITICAL
	} else if secondsBehindMaster > rOpts.Warn {
		checkSt = checkers.WARNING
	}
	return checkers.NewChecker(checkSt, msg)
}
