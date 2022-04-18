package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/CortexFoundation/CortexTheseus/log"
	"github.com/CortexFoundation/torrentfs"
	cli "gopkg.in/urfave/cli.v1"
)

type Config struct {
	wg sync.WaitGroup
}

func main() {
	var conf Config
	app := cli.NewApp()
	app.Flags = []cli.Flag{}

	app.Action = func(c *cli.Context) error {
		err := run(&conf)
		return err
	}
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

func run(conf *Config) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	go http.ListenAndServe("127.0.0.1:8080", mux)

	config := &torrentfs.DefaultConfig
	config.DataDir = ".data"
	fs, err := torrentfs.New(config, true, false, false)
	if err != nil {
		log.Error("err", "e", err)
	}

	fmt.Println(err)

	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	conf.wg.Add(1)
	go func() {
		defer conf.wg.Done()
		fs.Start(nil)
	}()

	conf.wg.Wait()
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v, %v, %v\n", r.URL, r.Method, r.URL.Path)
	res := "OK"
	//uri := r.URL.Path
	switch r.Method {
	case "GET":
		//res = Get(uri)
	case "POST":
		// TODO
	default:
		res = "method not found"
	}
	fmt.Fprintf(w, res)
}
