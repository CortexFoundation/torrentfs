// Copyright 2020 The CortexTheseus Authors
// This file is part of the CortexTheseus library.
//
// The CortexTheseus library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The CortexTheseus library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the CortexTheseus library. If not, see <http://www.gnu.org/licenses/>.

package torrentfs

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestLocal(t *testing.T) {
	DefaultConfig.DataDir = "data"
	DefaultConfig.Port = 0
	fs, err := New(&DefaultConfig, true, false, false)
	if err != nil {
		log.Fatal(err)
	}
	fs.Start(nil)
	defer fs.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if _, err := fs.SeedingLocal(ctx, "torrent.go", false); err != nil {
		log.Fatal("failed to get file")
	}

}

func TestInfoHash(t *testing.T) {
	hash, err := Hash("testdata/data")
	if len(hash) == 0 || err != nil {
		log.Fatal(err)
	}
	fmt.Println(hash)
}
