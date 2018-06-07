package checkredis

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fzzy/radix/redis"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type redisSetting struct {
	Host    string `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Socket  string `short:"s" long:"socket" default:"" description:"Server socket"`
	Port    string `short:"p" long:"port" default:"6379" description:"Port"`
	Timeout uint64 `short:"t" long:"timeout" default:"5" description:"Dial Timeout in sec"`
}

var commands = map[string](func([]string) *checkers.Checker){
	"reachable": checkReachable,
	"slave":     checkSlave,
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

func connectRedis(m redisSetting) (*redis.Client, error) {
	network := "tcp"
	target := fmt.Sprintf("%s:%s", m.Host, m.Port)
	if m.Socket != "" {
		target = m.Socket
		network = "unix"
	}
	c, err := redis.DialTimeout(network, target, time.Duration(m.Timeout)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect: %s", err)
	}
	return c, nil
}

func getRedisInfo(c *redis.Client) (*map[string]string, error) {
	info := make(map[string]string)

	r := c.Cmd("info")
	if r.Err != nil {
		return nil, errors.New("couldn't execute query")
	}
	str, err := r.Str()
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

func connectRedisGetInfo(opts redisSetting) (*redis.Client, *map[string]string, error) {
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

func checkSlave(args []string) *checkers.Checker {
	opts := redisSetting{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Usage = "slave [OPTIONS]"
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
