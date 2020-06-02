package torrentfs

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestGetFile(t *testing.T) {
	DefaultConfig.DataDir = "/root/.cortex/storage"
	fmt.Println(DefaultConfig)
	tm, _ := NewTorrentManager(&DefaultConfig, 1, true, false)
	tm.Start()
	defer tm.Close()
	time.Sleep(10 * time.Second)
	a, _ := tm.Available("3f1f6c007e8da3e16f7c3378a20a746e70f1c2b0", 100000000)
	fmt.Println("available", a)
	file, _ := tm.GetFile("3f1f6c007e8da3e16f7c3378a20a746e70f1c2b0", "data")
	if file == nil {
		log.Fatal("failed to get file")
	}
	fmt.Println("file", file[:20])
}
