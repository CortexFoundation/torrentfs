package torrentfs

import (
	"fmt"
	"log"
	"testing"
)

func TestGetFile(t *testing.T) {
	tm := NewTorrentManager(DefaultConfig{}, 1, true, false)
	tm.Start()
	defer tm.Close()
	file := tm.GetFile("3f1f6c007e8da3e16f7c3378a20a746e70f1c2b0", "data")
	if file == nil {
		log.Fatal("failed to get file")
	}
	fmt.Println("file", file)
}
