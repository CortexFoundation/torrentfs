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

package compress

import (
	"fmt"
	"log"
	"testing"
)

func TestCountValues(t *testing.T) {
	// define original data
	data := []byte(`MkjdflskfjslkfjalkfjalfkjalfjafyzYrIyMLyNqwDSTBqSwM2D6KD9sA8S/d3Vyy6ldE+oRVdWyqNQrjTxQ6uG3XBOS0P4GGaIMJEPQ/gYZogwkQ+A0/gSU03fRJvdhIGQ1AMARVdWyqNQrjRFV1bKo1CuNEVXVsqjUK40RVdWyqNQrjRFV1bKo1CuNPmQF870PPsnSNeKI1U/MrOA0/gSU03fRb2A3OsnORNIruhCUYTIrsdfsfsdfsfOMTNU7JuGb5dfsfsfsfdsfsfsfRSYJxa6PiMHdiRmFtsdfsfsdfXLNoY+GVmTD7aOV/K1yo4ydfsfsf0dR7Q=ilsdjflskfjlskdfjldsf`)
	fmt.Println("original data:", data)
	fmt.Println("original data len:", len(data))

	// compress data
	compressedData, compressedDataErr := ZipData(data)
	if compressedDataErr != nil {
		log.Fatal(compressedDataErr)
	}

	fmt.Println("compressed data:", compressedData)
	fmt.Println("compressed data len:", len(compressedData))

	// uncompress data
	uncompressedData, uncompressedDataErr := UnzipData(compressedData)
	if uncompressedDataErr != nil {
		log.Fatal(uncompressedDataErr)
	}
	fmt.Println("uncompressed data:", uncompressedData)
	fmt.Println("uncompressed data len:", len(uncompressedData))

	compressedData, compressedDataErr = SnappyEncode(data)
	if compressedDataErr != nil {
		log.Fatal(compressedDataErr)
	}
	fmt.Println("compressed data:", compressedData)
	fmt.Println("compressed data len:", len(compressedData))

	uncompressedData, uncompressedDataErr = SnappyDecode(compressedData)
	if uncompressedDataErr != nil {
		log.Fatal(uncompressedDataErr)
	}
	fmt.Println("uncompressed data:", uncompressedData)
	fmt.Println("uncompressed data len:", len(uncompressedData))
}
