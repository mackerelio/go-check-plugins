package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mackerelio/checkers"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type mysqlSetting struct {
	Host string `short:"h" long:"host" default:"localhost" description:"Hostname"`
	Port string `short:"p" long:"port" default:"3306" description:"Port"`
	User string `short:"u" long:"user" default:"root" description:"Username"`
	Pass string `short:"P" long:"password" default:"" description:"Password"`
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

func newMySQL(m mysqlSetting) mysql.Conn {
	target := fmt.Sprintf("%s:%s", m.Host, m.Port)
	return mysql.New("tcp", "", target, m.User, m.Pass, "")
}
