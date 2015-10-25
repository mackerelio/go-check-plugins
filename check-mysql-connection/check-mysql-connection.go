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
	Crit int64  `short:"c" long:"critical" default:"250" description:"critical if the number of connection is over"`
	Warn int64  `short:"w" long:"warning" default:"200" description:"warning if the number of connection is over"`
}

func main() {
	ckr := run(os.Args[1:])
	ckr.Name = "MySQL Connection"
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

	rows, res, err := db.Query("SHOW GLOBAL STATUS LIKE 'Threads_Connected'")
	if err != nil {
		return checkers.Unknown("couldn't execute query")
	}

	idxValue := res.Map("Value")
	threadsConnected := rows[0].Int64(idxValue)

	checkSt := checkers.OK
	msg := fmt.Sprintf("%d connections", threadsConnected)
	if threadsConnected > opts.Crit {
		checkSt = checkers.CRITICAL
	} else if threadsConnected > opts.Warn {
		checkSt = checkers.WARNING
	}
	return checkers.NewChecker(checkSt, msg)
}
