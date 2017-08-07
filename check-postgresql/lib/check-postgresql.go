package checkpostgresql

import (
	"fmt"
	"os"
	"strings"

	// PostgreSQL Driver
	_ "github.com/lib/pq"
	"github.com/mackerelio/checkers"
)

var commands = map[string](func([]string) *checkers.Checker){
	"connection": checkConnection,
}

type postgresqlSetting struct {
	Host     string `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Port     string `short:"p" long:"port" default:"5432" description:"Port"`
	User     string `short:"u" long:"user" default:"postgres" description:"Username"`
	Password string `short:"P" long:"password" default:"" description:"Password"`
	Database string `short:"d" long:"database" description:"DBname"`
	SSLmode  string `short:"s" long:"sslmode" default:"disable" description:"SSLmode"`
	Timeout  int    `short:"t" long:"timeout" default:"5" description:"Maximum wait for connection, in seconds."`
}

func (p postgresqlSetting) getDriverAndDataSourceName() (string, string) {
	dbName := p.User
	if p.Database != "" {
		dbName = p.Database
	}
	dataSourceName := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s connect_timeout=%d", p.User, p.Password, p.Host, p.Port, dbName, p.SSLmode, p.Timeout)
	return "postgres", dataSourceName
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
  check-postgresql [subcommand] [OPTIONS]

SubCommands:`)
		for k := range commands {
			fmt.Printf("  %s\n", k)
		}
		os.Exit(1)
	}
	ckr := fn(argv)
	ckr.Name = fmt.Sprintf("PostgreSQL %s", strings.Title(subCmd))
	ckr.Exit()
}
