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

	"github.com/CortexFoundation/torrentfs/backend"
	"github.com/CortexFoundation/torrentfs/params"
)

func TestLocal(t *testing.T) {
	params.DefaultConfig.DataDir = "testdata"
	params.DefaultConfig.Port = 0
	params.DefaultConfig.Mode = "LAZY"
	fs, err := New(&params.DefaultConfig, true, false, false)
	if err != nil {
		log.Fatal(err)
	}
	fs.Start()
	//fs.Simulate()
	defer fs.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if _, err := fs.SeedingLocal(ctx, "fs.go", false); err != nil {
		log.Fatal("failed to seed local file")
	}
	time.Sleep(10 * time.Second)
}

func TestInfoHash(t *testing.T) {
	hash, err := backend.Hash("testdata/data")
	if len(hash) == 0 || err != nil {
		log.Fatal(err)
	}
	fmt.Println(hash)
}
