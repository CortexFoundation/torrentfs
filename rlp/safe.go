// Copyright 2021 the CortexTheseu Authors
// This file is part of the CortexTheseu library.
//
// the CortexTheseu library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the CortexTheseu library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the CortexTheseu library. If not, see <http://www.gnu.org/licenses/>.

//go:build nacl || js || !cgo
// +build nacl js !cgo

package rlp

import "reflect"

// byteArrayBytes returns a slice of the byte array v.
func byteArrayBytes(v reflect.Value, length int) []byte {
	return v.Slice(0, length).Bytes()
}
