package main

import (
	"flag"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"os"
)

// The exit status of the commands
const (
	OK       = 0
	WARNING  = 1
	CRITICAL = 2
	UNKNOWN  = 3
)

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "3306", "Port")
	optUser := flag.String("username", "root", "Username")
	optPass := flag.String("password", "", "Password")
	optCrit := flag.Int("crit", 1, "critical if the second behind master is over")
	optWarn := flag.Int("warn", 1, "warning if the second behind master is over")
	flag.Parse()

	target := fmt.Sprintf("%s:%s", *optHost, *optPort)
	db := mysql.New("tcp", "", target, *optUser, *optPass, "")
	err := db.Connect()
	if err != nil {
		fmt.Println("UNKNOWN: couldn't connect DB")
		os.Exit(UNKNOWN)
	}
	defer db.Close()

	rows, res, err := db.Query("show slave status")
	if err != nil {
		fmt.Println("UNKNOWN: couldn't execute query")
		os.Exit(UNKNOWN)
	}
	if len(rows) == 0 {
		fmt.Println("OK: MySQL is not slave")
		os.Exit(OK)
	}

	idxIoThreadRunning := res.Map("Slave_IO_Running")
	idxSQLThreadRunning := res.Map("Slave_SQL_Running")
	idxSecondsBehindMaster := res.Map("Seconds_Behind_Master")
	ioThreadStatus := rows[0].Str(idxIoThreadRunning)
	sqlThreadStatus := rows[0].Str(idxSQLThreadRunning)
	secondsBehindMaster := rows[0].Int(idxSecondsBehindMaster)

	if ioThreadStatus == "No" || sqlThreadStatus == "No" {
		fmt.Println("CRITICAL: MySQL replication has been stopped")
		os.Exit(CRITICAL)
	}

	if secondsBehindMaster > *optCrit {
		msg := fmt.Sprintf("CRITICAL: MySQL replication behind master %d seconds", secondsBehindMaster)
		fmt.Println(msg)
		os.Exit(CRITICAL)
	} else if secondsBehindMaster > *optWarn {
		msg := fmt.Sprintf("WARNING: MySQL replication behind master %d seconds", secondsBehindMaster)
		fmt.Println(msg)
		os.Exit(WARNING)
	}

	fmt.Println("OK: MySQL replication works well")
	os.Exit(OK)
}
