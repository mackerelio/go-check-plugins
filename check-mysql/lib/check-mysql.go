package checkmysql

import (
	"fmt"
	"os"
	"strings"

	"github.com/mackerelio/checkers"
	"github.com/ziutek/mymysql/mysql"
	// MySQL Driver
	_ "github.com/ziutek/mymysql/native"
)

type mysqlSetting struct {
	Host   string `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Port   string `short:"p" long:"port" default:"3306" description:"Port"`
	Socket string `short:"S" long:"socket" default:"" description:"Path to unix socket"`
	User   string `short:"u" long:"user" default:"root" description:"Username"`
	Pass   string `short:"P" long:"password" default:"" description:"Password" env:"MYSQL_PASSWORD"`
}

var commands = map[string](func([]string) *checkers.Checker){
	"replication": checkReplication,
	"connection":  checkConnection,
	"uptime":      checkUptime,
	"readonly":    checkReadOnly,
}

func separateSub(argv []string) (string, []string) {
	if len(argv) == 0 || strings.HasPrefix(argv[0], "-") {
		return "", argv
	}
	return argv[0], argv[1:]
}

// Do the plugin
func Do() {
	subCmd, argv := separateSub(os.Args[1:])
	fn, ok := commands[subCmd]
	if !ok {
		fmt.Println(`Usage:
  check-mysql [subcommand] [OPTIONS]

SubCommands:`)
		for k := range commands {
			fmt.Printf("  %s\n", k)
		}
		os.Exit(1)
	}
	ckr := fn(argv)
	ckr.Name = fmt.Sprintf("MySQL %s", strings.ToUpper(string(subCmd[0]))+subCmd[1:])
	ckr.Exit()
}

func newMySQL(m mysqlSetting) mysql.Conn {
	proto := "tcp"
	target := fmt.Sprintf("%s:%s", m.Host, m.Port)
	if m.Socket != "" {
		proto = "unix"
		target = m.Socket
	}
	return mysql.New(proto, "", target, m.User, m.Pass, "")
}
