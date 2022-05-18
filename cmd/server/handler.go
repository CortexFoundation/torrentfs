package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/CortexFoundation/CortexTheseus/log"
)

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
			match, _ := regexp.MatchString(`^[0-9A-Za-z._-]*$`, name)
			if name == "torrent" || !match || strings.Contains(name, "/") || strings.Contains(name, "\\") {
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
