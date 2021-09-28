package checkredis

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type redisSetting struct {
	Host     string `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Socket   string `short:"s" long:"socket" default:"" description:"Server socket"`
	Port     string `short:"p" long:"port" default:"6379" description:"Port"`
	Password string `short:"P" long:"password" default:"" description:"Password"`
	Timeout  uint64 `short:"t" long:"timeout" default:"5" description:"Dial Timeout in sec"`
}

var commands = map[string](func([]string) *checkers.Checker){
	"reachable":   checkReachable,
	"replication": checkReplication,
	"slave":       checkSlave, // deprecated command
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
  check-redis [subcommand] [OPTIONS]

SubCommands:`)
		for k := range commands {
			fmt.Printf("  %s\n", k)
		}
		os.Exit(1)
	}
	ckr := fn(argv)
	ckr.Name = fmt.Sprintf("Redis %s", strings.ToUpper(string(subCmd[0]))+subCmd[1:])
	ckr.Exit()
}

func connectRedis(m redisSetting) (redis.Conn, error) {
	network := "tcp"
	address := net.JoinHostPort(m.Host, m.Port)
	if m.Socket != "" {
		network = "unix"
		address = m.Socket
	}
	c, err := redis.Dial(network, address, redis.DialConnectTimeout(time.Duration(m.Timeout)*time.Second))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect: %s", err)
	}

	password := ""

	if os.Getenv("REDIS_PASSWORD") != "" {
		password = os.Getenv("REDIS_PASSWORD")
	}

	if m.Password != "" {
		password = m.Password
	}

	if password != "" {
		_, err := c.Do("AUTH", password)
		if err != nil {
			return nil, fmt.Errorf("couldn't authenticate: %v", err)
		}
	}

	return c, nil
}

func getRedisInfo(c redis.Conn) (*map[string]string, error) {
	info := make(map[string]string)

	str, err := redis.String(c.Do("info"))
	if err != nil {
		return nil, errors.New("couldn't execute query")
	}

	for _, line := range strings.Split(str, "\r\n") {
		if line == "" {
			continue
		}
		if re, _ := regexp.MatchString("^#", line); re {
			continue
		}

		record := strings.SplitN(line, ":", 2)
		if len(record) < 2 {
			continue
		}
		key, value := record[0], record[1]
		info[key] = value
	}

	return &info, nil
}

func connectRedisGetInfo(opts redisSetting) (redis.Conn, *map[string]string, error) {
	c, err := connectRedis(opts)
	if err != nil {
		return nil, nil, err
	}

	info, err := getRedisInfo(c)
	if err != nil {
		return nil, nil, err
	}

	return c, info, nil
}

func checkReachable(args []string) *checkers.Checker {
	opts := redisSetting{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Usage = "reachable [OPTIONS]"
	_, err := psr.ParseArgs(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c, info, err := connectRedisGetInfo(opts)
	if err != nil {
		return checkers.Unknown(err.Error())
	}
	defer c.Close()

	if _, ok := (*info)["redis_version"]; !ok {
		return checkers.Unknown("couldn't get redis_version")
	}

	return checkers.Ok(
		fmt.Sprintf("version: %s", (*info)["redis_version"]),
	)
}

type replicationOpts struct {
	redisSetting
	SkipMaster bool `long:"skip-master" description:"return ok if redis role is master"`
}

func checkReplication(args []string) *checkers.Checker {
	opts := replicationOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Usage = "replication [OPTIONS]"
	_, err := psr.ParseArgs(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c, info, err := connectRedisGetInfo(opts.redisSetting)
	if err != nil {
		return checkers.Unknown(err.Error())
	}
	defer c.Close()

	if role, ok := (*info)["role"]; ok {
		if role != "slave" && opts.SkipMaster {
			return checkers.Ok("role is not slave")
		}
	} else {
		return checkers.Unknown("couldn't get role")
	}

	if status, ok := (*info)["master_link_status"]; ok {
		msg := fmt.Sprintf("master_link_status: %s", status)

		switch status {
		case "up":
			return checkers.Ok(msg)
		case "down":
			return checkers.Critical(msg)
		default:
			return checkers.Unknown(msg)
		}

	} else {
		return checkers.Unknown("couldn't get master_link_status")
	}
}

// Deprecated: For backward compatibility.
func checkSlave(args []string) *checkers.Checker {
	opts := redisSetting{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Usage = `slave [OPTIONS]

DEPRECATED: For backward compatibility. Use 'replication' command.`
	_, err := psr.ParseArgs(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c, info, err := connectRedisGetInfo(opts)
	if err != nil {
		return checkers.Unknown(err.Error())
	}
	defer c.Close()

	if status, ok := (*info)["master_link_status"]; ok {
		msg := fmt.Sprintf("master_link_status: %s", status)

		switch status {
		case "up":
			return checkers.Ok(msg)
		case "down":
			return checkers.Critical(msg)
		default:
			return checkers.Unknown(msg)
		}

	} else {
		// it may be a master!
		return checkers.Unknown("couldn't get master_link_status")
	}
}
