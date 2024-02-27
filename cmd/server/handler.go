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

	//"strconv"

	"github.com/CortexFoundation/CortexTheseus/common"
	"github.com/CortexFoundation/CortexTheseus/log"

	"github.com/CortexFoundation/torrentfs/backend"
)

const (
	WORKSPACE = "/share/"
)

func (conf *Config) DropHandler(w http.ResponseWriter, r *http.Request) {
	res := "OK"
	q := r.URL.Query()
	switch r.Method {
	case "GET":
		res = "GET NOT SUPPORT"
	case "POST":
		ih := q.Get("hash")
		if !common.IsHexAddress(ih) {
			res = "invalid hash"
		} else {
			err := conf.tfs.Drop(ih)
			if err != nil {
				log.Error("err", "e", err)
				res = err.Error()
			}
		}
	default:
		res = "method not found"
	}
	fmt.Fprintf(w, res)
}

func (conf *Config) TunnelHandler(w http.ResponseWriter, r *http.Request) {
	res := "OK"
	q := r.URL.Query()
	switch r.Method {
	case "GET":
		res = "GET NOT SUPPORT"
	case "POST":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ih := q.Get("hash")
		if !common.IsHexAddress(ih) {
			res = "invalid hash"
		} else {
			go func() {
				timer := time.NewTicker(900 * time.Second)
				defer timer.Stop()
				err := conf.tfs.Tunnel(ctx, ih)
				if err != nil {
					log.Error("err", "e", err)
					res = err.Error()
				}

				select {
				case <-timer.C:
					conf.tfs.Drop(ih)
				}
			}()
		}
	default:
		res = "method not found"
	}
	fmt.Fprintf(w, res)
}

func (conf *Config) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	res := "OK"
	q := r.URL.Query()
	switch r.Method {
	case "GET":
		res = "GET NOT SUPPORT"
	case "POST":
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		ih := q.Get("hash")
		if !common.IsHexAddress(ih) {
			res = "invalid hash"
		} else {
			err := conf.tfs.Tunnel(ctx, ih)
			if err != nil {
				log.Error("err", "e", err)
				res = err.Error()
			}
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
		//path, err := os.Getwd()
		//path := "/share"
		//if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		//	log.Error("Mkdir failed", "path", path, "err", err)
		//	res = err.Error()
		//	return
		//} else {
		//ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		//defer cancel()
		name := q.Get("file")
		match, _ := regexp.MatchString(`^[0-9A-Za-z._-]*$`, name)
		if len(name) == 0 || name == "torrent" || !match || strings.Contains(name, "/") || strings.Contains(name, "\\") {
			log.Error("invalid file name", "name", name)
			res = "invalid file name pattern"
		} else {
			path := WORKSPACE
			file := filepath.Join(path, name)

			log.Info("Seeding path", "root", path, "file", file, "name", name)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			if ih, err := conf.tfs.SeedingLocal(ctx, file, false); err != nil {
				log.Error("err", "e", err)
				res = err.Error()
			} else {
				res = ih
				conf.db.Set([]byte("s:"+ih), []byte(ih))
			}
		}
		//}
	default:
		res = "method not found"
	}
	fmt.Fprintf(w, res)
}

func (conf *Config) ListHandler(w http.ResponseWriter, r *http.Request) {
	res := ""
	//q := r.URL.Query()
	switch r.Method {
	case "GET":
		res = "GET NOT SUPPORT"
	case "POST":
		files, err := os.ReadDir(WORKSPACE)
		if err != nil {
			res = err.Error()
		}

		var path string
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if file.Type() == os.ModeSymlink {
				continue
			}
			//if ok, err := common.IsDirectory(WORKSPACE + file.Name()); ok || err == nil {
			//	log.Error("Dir failed", "err", err , "ok", ok, "path", WORKSPACE + file.Name())
			//	continue
			//}
			var h string
			path = WORKSPACE + file.Name()
			if b := conf.db.Get([]byte(file.Name())); b != nil {
				h = string(b)
			} else {
				h, err = backend.Hash(path)
				if len(h) == 0 || err != nil {
					log.Error("Hash failed", "path", path, "err", err)
					continue
				}
				conf.db.Set([]byte(file.Name()), []byte(h))
			}
			res += file.Name() + "     " + h + "\n"
		}
	default:
		res = "method not found"
	}
	fmt.Fprintf(w, res)
}

/*func (conf *Config) FetchHandler(w http.ResponseWriter, r *http.Request) {
	res := "OK"
	q := r.URL.Query()
	switch r.Method {
	case "GET":
		res = "GET NOT SUPPORT"
	case "POST":
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		n, err := strconv.ParseInt(q.Get("size"), 10, 64)
		if err != nil {
			fmt.Printf("%d of type %T", n, n)
			res = "size failed"
		} else {
			_, err := conf.tfs.GetFileWithSize(ctx, q.Get("hash"), uint64(n), q.Get("subpath"))
			if err != nil {
				log.Error("err", "e", err)
				res = err.Error()
			} else {
				//res = string(ret[:])
			}
		}
	default:
		res = "method not found"
	}
	fmt.Fprintf(w, res)
}*/
