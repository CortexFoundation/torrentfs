package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/CortexFoundation/CortexTheseus/log"
	t "github.com/CortexFoundation/torrentfs"
	"github.com/CortexFoundation/torrentfs/params"
	xprometheus "github.com/anacrolix/missinggo/v2/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cli "github.com/urfave/cli/v2"
)

type Config struct {
	tfs  *t.TorrentFS
	dir  string
	port string
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

	PortFlag := cli.StringFlag{
		Name:  "port",
		Usage: "Listen port",
		Value: "8080",
	}

	app.Flags = []cli.Flag{
		&StorageFlag,
		&PortFlag,
	}

	app.Action = func(ctx *cli.Context) error {
		conf.dir = ctx.String(StorageFlag.Name)
		conf.port = ctx.String(PortFlag.Name)
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
	config.Mode = params.LAZY
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

	prometheus.MustRegister(xprometheus.NewExpvarCollector())

	mux := http.NewServeMux()
	mux.HandleFunc("/download", conf.DownloadHandler)
	mux.HandleFunc("/seed", conf.SeedHandler)
	mux.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe("127.0.0.1:"+conf.port, mux)

	var c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}
