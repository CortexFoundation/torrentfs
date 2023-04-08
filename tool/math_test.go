// Copyright 2023 The CortexTheseus Authors
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

package tool

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	for i := 0; i < 10; i++ {
		r := Rand(100)
		fmt.Printf("%v ", r)
	}
	fmt.Println()
}

func TestRand(t *testing.T) {
	// Test case 1: s = 0
	if Rand(0) != 0 {
		t.Errorf("Rand(0) should return 0")
	}

	// Test cases 2-3: s = 1 or 2
	/*for i := int64(0); i <= 2; i++ {
		r := Rand(i)
		if i < 16 {
			if r != smallRandTable[i] {
				t.Errorf("Rand(%d) should return %d", i, smallRandTable[i])
			}
		} else {
			if r != 0 && r != 1 {
				t.Errorf("Rand(%d) should return 0 or 1", i)
			}
		}
	}*/

	// Test cases 4-6: s = 16, 100, 1000000
	for _, s := range []int64{16, 100, 1000000} {
		r := Rand(s)
		if r < 0 || r >= s {
			t.Errorf("Rand(%d) should return a value between 0 and %d", s, s-1)
		}
	}
}
