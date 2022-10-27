package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/CortexFoundation/CortexTheseus/log"
	t "github.com/CortexFoundation/torrentfs"
	"github.com/CortexFoundation/torrentfs/params"
	//"github.com/CortexFoundation/torrentfs/wormhole"
	xprometheus "github.com/anacrolix/missinggo/v2/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ucwong/golang-kv"
	cli "github.com/urfave/cli/v2"
)

type Config struct {
	tfs  *t.TorrentFS
	dir  string
	port string
	db   kv.Bucket
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
		Value: "7882",
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
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(true)))
	glogger.Verbosity(log.LvlDebug)
	glogger.Vmodule("")
	log.Root().SetHandler(glogger)
	conf.db = kv.Badger("")
	//conf.db = kv.Bolt("")
	//conf.db = kv.LevelDB("")
	if conf.db != nil {
		defer conf.db.Close()
	}

	config := &params.DefaultConfig
	config.DataDir = conf.dir
	config.Mode = params.FULL
	config.Port = 0
	config.Server = true
	config.Wormhole = false

	//config.DisableUTP = false

	fs, err := t.New(config, false, false, false)
	if err != nil {
		log.Error("err", "e", err)
		return err
	}
	//fs.Start(nil)
	defer fs.Stop()

	conf.tfs = fs

	prometheus.MustRegister(xprometheus.NewExpvarCollector())

	mux := http.NewServeMux()
	mux.HandleFunc("/download", conf.DownloadHandler)
	mux.HandleFunc("/tunnel", conf.TunnelHandler)
	mux.HandleFunc("/seed", conf.SeedHandler)
	mux.HandleFunc("/list", conf.ListHandler)
	mux.HandleFunc("/drop", conf.DropHandler)

	fileServer := http.FileServer(http.Dir("./.storage/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/metrics", promhttp.Handler())

	log.Info("Server started", "port", conf.port)
	//ret := conf.db.Prefix([]byte("s:"))
	//ret := conf.db.Scan()
	//log.Info("db length", "len", len(ret))
	//for _, v := range ret {
	//	log.Info("Seeding file", "ih", string(v))
	//}

	//wormhole.BestTrackers()

	if err := http.ListenAndServe("127.0.0.1:"+conf.port, mux); err != nil {
		log.Error("Failed to start server", "err", err)
		return err
	}

	var c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}
