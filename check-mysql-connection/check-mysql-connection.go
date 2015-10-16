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
	optCrit := flag.Int("crit", 250, "critical if the number of connection is over")
	optWarn := flag.Int("warn", 200, "warning if the number of connection is over")
	flag.Parse()

	target := fmt.Sprintf("%s:%s", *optHost, *optPort)
	db := mysql.New("tcp", "", target, *optUser, *optPass, "")
	err := db.Connect()
	if err != nil {
		fmt.Println("UNKNOWN: couldn't connect DB")
		os.Exit(UNKNOWN)
	}
	defer db.Close()

	rows, res, err := db.Query("show global status like 'Threads_Connected'")
	if err != nil {
		fmt.Println("UNKNOWN: couldn't execute query")
		os.Exit(UNKNOWN)
	}

	idxValue := res.Map("Value")
	threadsConnected := rows[0].Int(idxValue)

	if threadsConnected > *optCrit {
		msg := fmt.Sprintf("CRITICAL: %d connections", threadsConnected)
		fmt.Println(msg)
		os.Exit(CRITICAL)
	} else if threadsConnected > *optWarn {
		msg := fmt.Sprintf("WARNING: %d connections", threadsConnected)
		fmt.Println(msg)
		os.Exit(WARNING)
	}

	msg := fmt.Sprintf("OK: %d connections", threadsConnected)
	fmt.Println(msg)
	os.Exit(OK)
}
