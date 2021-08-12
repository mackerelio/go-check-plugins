package checkmysql

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/mackerelio/checkers"
)

type mysqlSetting struct {
	Host   string `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Port   string `short:"p" long:"port" default:"3306" description:"Port"`
	Socket string `short:"S" long:"socket" default:"" description:"Path to unix socket"`
	User   string `short:"u" long:"user" default:"root" description:"Username"`
	Pass   string `short:"P" long:"password" default:"" description:"Password" env:"MYSQL_PASSWORD"`

	EnableTLS     bool   `long:"tls" description:"Enables TLS connection"`
	TLSRootCert   string `long:"tls-root-cert" default:"" description:"The root certificate used for TLS certificate verification"`
	TLSSkipVerify bool   `long:"tls-skip-verify" description:"Disable TLS certificate verification"`
}

type mysqlVersion struct {
	major int
	minor int
	patch int
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

func newDB(m mysqlSetting) (*sql.DB, error) {
	proto := "tcp"
	target := fmt.Sprintf("%s:%s", m.Host, m.Port)
	if m.Socket != "" {
		proto = "unix"
		target = m.Socket
	}
	cfg := &mysql.Config{
		User:                 m.User,
		Passwd:               m.Pass,
		Net:                  proto,
		Addr:                 target,
		AllowNativePasswords: true,
	}
	if m.EnableTLS {
		var c tls.Config
		if m.TLSRootCert != "" {
			certPool := x509.NewCertPool()
			pem, err := os.ReadFile(m.TLSRootCert)
			if err != nil {
				return nil, fmt.Errorf("cannot read %s: %v", m.TLSRootCert, err)
			}
			certPool.AppendCertsFromPEM(pem)
			c.RootCAs = certPool
		}
		c.InsecureSkipVerify = m.TLSSkipVerify
		mysql.RegisterTLSConfig("custom", &c)
		cfg.TLSConfig = "custom"
	}

	return sql.Open("mysql", cfg.FormatDSN())
}

func getMySQLVersion(db *sql.DB) (*mysqlVersion, error) {
	var rawVersion string
	err := db.QueryRow("SELECT VERSION()").Scan(&rawVersion)
	if err != nil {
		return nil, fmt.Errorf("Failed to query SELECT VERSION(): %s", err)
	}

	// Version example: 5.5.44-0+deb8u1-log
	var major, minor, patch int
	_, err = fmt.Sscanf(rawVersion, "%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse version: %s", err)
	}

	return &mysqlVersion{
		major: major,
		minor: minor,
		patch: patch,
	}, nil
}
