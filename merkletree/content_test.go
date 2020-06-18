package merkletree

import (
	"github.com/status-im/keycard-go/hexutils"
	"testing"
)

func TestBlockContent_CalculateHash(t *testing.T)  {
	bc := BlockContent{
		X: "blockContent test",
		N: 0,
	}
	hash, err := bc.CalculateHash()
	if err != nil {
		t.Fatal("hash err: ", err)
	}

	want := "1350F92C4DD40B6CF5340C1FF1F0278A6558417EAE315C52BDE4E5A107E3B9B2"
	got := hexutils.BytesToHex(hash)
	if got != want {
		t.Fatalf("calculated hash mismatch: got %s, want %s", got, want)
	}
}

func TestBlockContent_Equals(t *testing.T) {
	bc1 := BlockContent{
		X: "blockContent1",
		N: 0,
	}
	bc2 := BlockContent{
		X: "blockContent2",
		N: 0,
	}
	bc3 := BlockContent{
		X: "blockContent2",
		N: 0,
	}
	equal12, err := bc1.Equals(bc2)
	if err != nil {
		t.Fatal("equal err: ", err)
	}
	if equal12 != false {
		t.Fatalf("blockContent1 and blockContent2 should be unequal")
	}
	equal23, err := bc2.Equals(bc3)
	if err != nil {
		t.Fatal("equal err: ", err)
	}
	if equal23 != true {
		t.Fatalf("blockContent2 and blockContent3 should be equal")
	}
}
