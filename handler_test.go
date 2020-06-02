package torrentfs

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestGetFile(t *testing.T) {
	DefaultConfig.DataDir = "data"
	ih := "aea5584d0cd3865e90c80eace3bfcb062473d966"
	fmt.Println(DefaultConfig)
	tm, _ := NewTorrentManager(&DefaultConfig, 1, false, false)
	tm.Search(ih, 0)
	tm.Start()
	defer tm.Close()
	time.Sleep(3 * time.Second)
	a, _ := tm.Available(ih, 100000000)
	fmt.Println("available", a)
	file, _ := tm.GetFile(ih, "data")
	if file == nil {
		log.Fatal("failed to get file")
	}
	fmt.Println("file", file[:20])
}
