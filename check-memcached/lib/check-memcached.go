package checkmemcached

import (
	"os"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts struct {
	Host    string `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Port    string `short:"p" long:"port" default:"11211" description:"Port"`
	Timeout uint64 `short:"t" long:"timeout" default:"3" description:"Dial Timeout in sec"`
	Key     string `short:"k" long:"key" required:"true" description:"Cache key used within set and get test"`
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Memcached"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	mc := memcache.New(opts.Host + ":" + opts.Port)
	mc.Timeout = time.Duration(opts.Timeout) * time.Second

	err = mc.Set(&memcache.Item{Key: opts.Key, Value: []byte("Check key"), Expiration: 240})
	if err != nil {
		return checkers.Critical("couldn't set a key: " + err.Error())
	}

	item, err := mc.Get(opts.Key)
	if err != nil {
		return checkers.Critical("couldn't get a key: " + err.Error())
	}
	if string(item.Value) != "Check key" {
		return checkers.Critical("not correct value")
	}
	return checkers.Ok("Get,Set OK")
}
