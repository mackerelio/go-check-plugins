package checkgateway

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts struct {
	Host      string `long:"host" short:"H" description:"check target IP Address" required:"true"`
	Gateway   string `long:"gateway" short:"g" description:"gateway MAC or IP or name" required:"true"`
	Count     int    `long:"count" short:"n" default:"1" description:"sending (and receiving) count ping packets"`
	Interface string `long:"interface" short:"I" description:"the interface on which the packets must be sent"`

	Thresholds struct {
		WarningLoss  int `long:"warning" short:"w" default:"30" description:"maximum percentage of allowed data loss"`
		CriticalLoss int `long:"critical" short:"c" default:"80"`
	} `group:"Tuning thresholds"`

	Paths struct {
		Arping string `long:"arping-path" default:"arping"`
		Arp    string `long:"arp-path" default:"arp"`
	} `group:"Binaries paths"`
}

func resolveIP(ip string) (string, error) {
	out, err := exec.Command(opts.Paths.Arp, ip).Output()
	if err != nil {
		return "", err
	}
	macre := regexp.MustCompile("..:..:..:..:..:..")
	if !macre.Match(out) {
		return "", fmt.Errorf("cannot find mac address for %s", ip)
	}
	return macre.FindString(string(out)), nil
}

func run(args []string) *checkers.Checker {
	var parser = flags.NewParser(&opts, flags.Default)
	_, err := parser.ParseArgs(args)
	if err != nil {
		os.Exit(1)
	}
	if opts.Host == "" {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if !strings.Contains(opts.Gateway, ":") {
		ra, err := net.ResolveIPAddr("ip4", opts.Gateway)
		if err != nil {
			return checkers.Critical("Could not resolve host")
		}
		ip := ra.IP.String()
		mac, err := resolveIP(ip)
		if err != nil {
			return checkers.Critical("Could not resolve IP")
		}
		opts.Gateway = mac
	}
	cmdArgs := []string{"-c", strconv.Itoa(opts.Count),
		"-T",
		opts.Host}
	if len(opts.Interface) > 0 {
		cmdArgs = append(cmdArgs, "-I", opts.Interface)
	}
	cmdArgs = append(cmdArgs, opts.Gateway)

	cmd := exec.Command(opts.Paths.Arping, cmdArgs...)
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(string(out))
		return checkers.Critical(err.Error())
	}
	outs := strings.Split(string(out), "\n")
	stats := outs[len(outs)-3]
	// TODO: times := outs[len(outs)-2]

	re := regexp.MustCompile(`(\d+)%`)
	if !re.MatchString(stats) || len(re.FindStringSubmatch(stats)) < 2 {
		return checkers.Critical("error parsing arping output")
	}
	perc, err := strconv.Atoi(re.FindStringSubmatch(stats)[1])
	if err != nil {
		return checkers.Critical("error parsing arping output: non-integer percentage found")
	}
	explaination := fmt.Sprintf("lost=%d%%", perc)
	if perc > opts.Thresholds.CriticalLoss {
		return checkers.Critical(explaination)
	}
	if perc > opts.Thresholds.WarningLoss {
		return checkers.Warning(explaination)
	}
	return checkers.Ok(explaination)
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Gateway"
	ckr.Exit()
}
