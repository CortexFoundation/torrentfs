package compress

import (
	"fmt"
	clog "github.com/CortexFoundation/CortexTheseus/log"
	"github.com/golang/snappy"
	"log"
	"testing"
	"time"
)

func SnappyEncode(data []byte) ([]byte, error) {
	start := time.Now()
	defer clog.Info("Snappy encode", "cost", time.Since(start))

	return snappy.Encode(nil, data), nil
}

func SnappyDecode(data []byte) ([]byte, error) {
	start := time.Now()
	defer clog.Info("Snappy decode", "cost", time.Since(start))
	res, err := snappy.Decode(nil, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}
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
