package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CortexFoundation/CortexTheseus/log"
	t "github.com/CortexFoundation/torrentfs"
	cli "github.com/urfave/cli/v2"
)

type Config struct {
	tfs *t.TorrentFS
	dir string
}

var (
	conf Config
)

func main() {

	app := cli.NewApp()

	StorageFlag := cli.StringFlag{
		Name:  "storage",
		Usage: "Data storage directory",
		Value: ".storage",
	}

	app.Flags = []cli.Flag{
		&StorageFlag,
	}

	app.Action = func(ctx *cli.Context) error {
		conf.dir = ctx.String(StorageFlag.Name)
		err := run(&conf)
		return err
	}
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

func run(conf *Config) error {
	config := &t.DefaultConfig
	config.DataDir = conf.dir
	config.Mode = t.LAZY
	fs, err := t.New(config, true, false, false)
	if err != nil {
		log.Error("err", "e", err)
		return err
	}
	defer fs.Stop()

	conf.tfs = fs

	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(true)))
	glogger.Verbosity(log.LvlDebug)
	glogger.Vmodule("")
	log.Root().SetHandler(glogger)

	mux := http.NewServeMux()
	mux.HandleFunc("/", conf.handler)
	http.ListenAndServe("127.0.0.1:8080", mux)

	var c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}

func (conf *Config) handler(w http.ResponseWriter, r *http.Request) {
	res := "OK"
	q := r.URL.Query()
	switch r.Method {
	case "GET":
		res = "GET NOT SUPPORT"
	case "POST":
		// TODO
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := conf.tfs.Download(ctx, q.Get("hash"), 1000000000)
		if err != nil {
			log.Error("err", "e", err)
			res = err.Error()
		}
	default:
		res = "method not found"
	}
	fmt.Fprintf(w, res)
}
