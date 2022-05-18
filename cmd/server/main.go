package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

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

func (conf *Config) DownloadHandler(w http.ResponseWriter, r *http.Request) {
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

func (conf *Config) SeedHandler(w http.ResponseWriter, r *http.Request) {
	res := "OK"

	q := r.URL.Query()
	switch r.Method {
	case "GET":
		res = "GET NOT SUPPORT"
	case "POST":
		path, err := os.Getwd()
		if err != nil {
			res = "seeding path failed"
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			name := q.Get("file")
			//name = strings.Replace(name, "/", "", -1)
			match, _ := regexp.MatchString(`^[0-9A-Za-z._]*$`, name)
			if !match || strings.Count(name, ".") > 1 || strings.Contains(name, "/") || strings.Contains(name, "\\") {
				log.Error("invalid file name", "name", name)
				res = "invalid file name pattern"
			} else {
				file := filepath.Join(path, name)
				log.Info("Seeding path", "root", path, "file", file)
				_, err := conf.tfs.SeedingLocal(ctx, file, false)
				if err != nil {
					log.Error("err", "e", err)
					res = err.Error()
				}
			}
		}
	default:
		res = "method not found"
	}
	fmt.Fprintf(w, res)
}
