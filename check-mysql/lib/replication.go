package checkmysql

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx"
	"github.com/mackerelio/checkers"
)

type replicationOpts struct {
	mysqlSetting
	Crit int64 `short:"c" long:"critical" default:"250" description:"critical if the seconds behind master is over"`
	Warn int64 `short:"w" long:"warning" default:"200" description:"warning if the seconds behind master is over"`
}

type status interface {
	ioRunning() string
	sqlRunning() string
	secondsBehind() sql.NullInt64
}

type replicationStatus struct {
	ReplicaIORunning    string        `db:"Replica_IO_Running"`
	ReplicaSQLRunning   string        `db:"Replica_SQL_Running"`
	SecondsBehindSource sql.NullInt64 `db:"Seconds_Behind_Source"`
}

func (r *replicationStatus) ioRunning() string {
	return r.ReplicaIORunning
}

func (r *replicationStatus) sqlRunning() string {
	return r.ReplicaSQLRunning
}

func (r *replicationStatus) secondsBehind() sql.NullInt64 {
	return r.SecondsBehindSource
}

type slaveStatus struct {
	SlaveIORunning      string        `db:"Slave_IO_Running"`
	SlaveSQLRunning     string        `db:"Slave_SQL_Running"`
	SecondsBehindMaster sql.NullInt64 `db:"Seconds_Behind_Master"`
}

func (r *slaveStatus) ioRunning() string {
	return r.SlaveIORunning
}

func (r *slaveStatus) sqlRunning() string {
	return r.SlaveSQLRunning
}

func (r *slaveStatus) secondsBehind() sql.NullInt64 {
	return r.SecondsBehindMaster
}

func checkReplication(args []string) *checkers.Checker {
	opts := replicationOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Usage = "replication [OPTIONS]"
	_, err := psr.ParseArgs(args)
	if err != nil {
		os.Exit(1)
	}
	db, err := newDB(opts.mysqlSetting)
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Couldn't open DB: %s", err))
	}
	defer db.Close()

	mySQLVersion, err := getMySQLVersion(db)
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Couldn't get MySQL Version: %s", err))
	}

	// MySQL > 8.0.22 supports `SHOW REPLICA STATUS`
	replicaSupport := !(mySQLVersion.major < 8 || (mySQLVersion.major == 8 && mySQLVersion.minor == 0 && mySQLVersion.patch < 22))

	sqlxDb := sqlx.NewDb(db, "mysql")
	defer sqlxDb.Close()

	var queryShowStatus string
	if replicaSupport {
		queryShowStatus = "SHOW REPLICA STATUS"
	} else {
		queryShowStatus = "SHOW SLAVE STATUS"
	}

	// Ignore columns which does not exist in structs.
	rows, err := sqlxDb.Unsafe().Queryx(queryShowStatus)
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Couldn't execute query: %s", err))
	}
	defer rows.Close()

	if !rows.Next() {
		return checkers.Ok("MySQL is not a replica")
	}

	var status status
	if replicaSupport {
		status = &replicationStatus{}
	} else {
		status = &slaveStatus{}
	}
	err = rows.StructScan(status)
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Couldn't scan row: %s", err))
	}

	if !(status.ioRunning() == "Yes" && status.sqlRunning() == "Yes") {
		return checkers.Critical("MySQL replication has been stopped")
	}

	checkSt := checkers.OK
	secondsBehind := status.secondsBehind()
	if !secondsBehind.Valid {
		return checkers.Unknown("Unknown seconds behind in MySQL replication")
	}

	msg := fmt.Sprintf("MySQL replication behind master %d seconds", secondsBehind.Int64)
	if secondsBehind.Int64 > opts.Crit {
		checkSt = checkers.CRITICAL
	} else if secondsBehind.Int64 > opts.Warn {
		checkSt = checkers.WARNING
	}
	return checkers.NewChecker(checkSt, msg)
}
