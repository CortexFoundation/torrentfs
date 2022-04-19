package main

import (
	"encoding/json"
	cli "github.com/urfave/cli/v2"
	"net/http"
	"os"
	"time"
	//"os/exec"
	"sync"
	//        "sync/atomic"
	"fmt"
)

type tracker_stats struct {
	Torrents              int `json:"torrents"`
	ActiveTorrents        int `json:"activeTorrents"`
	PeersAll              int `json:"peersAll"`
	PeersSeederOnly       int `json:"peersSeederOnly"`
	PeersLeecherOnly      int `json:"peersLeecherOnly"`
	PeersSeederAndLeecher int `json:"peersSeederAndLeecher"`
	PeersIPv4             int `json:"peersIPv4"`
	PeersIPv6             int `json:"peersIPv6"`
}

var ips = []string{
	"47.52.39.170",
	"3.17.249.140",
	"47.88.194.73",
	"47.91.207.20",
	"47.91.91.217",
	"47.88.7.24",
	"47.88.214.96",
	"47.89.178.175",
	"47.91.106.117",
	"47.91.43.70",
	"47.74.1.234",
}

var client http.Client

type Config struct {
	wg sync.WaitGroup
}

func http_healthy(ip string, port string) bool {
	url := "http://" + ip + ":" + port + "/stats.json"
	response, err := client.Get(url)
	if err != nil || response == nil || response.StatusCode != 200 {
		return false
	} else {
		var stats tracker_stats
		if jsErr := json.NewDecoder(response.Body).Decode(&stats); jsErr != nil {
			return false
		}
	}

	return true

}

// curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.34.0/install.sh | bash
// source ~/.bashrc
// nvm install node
// sudo npm install -g bittorrent-tracker
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
	client = http.Client{
		Timeout: time.Duration(3 * time.Second),
	}

	for _, ip := range ips {
		if http_healthy(ip, "5008") {
			fmt.Println(ip)
		} else {
			fmt.Println("ip failed : ", ip)
		}
	}
	conf.wg.Wait()
	//log.crit("cmd.Start", "err", err)
	return nil
}
