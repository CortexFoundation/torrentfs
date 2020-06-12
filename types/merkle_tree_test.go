package types

import (
	"bytes"
	"github.com/CortexFoundation/CortexTheseus/common"
	"testing"
)

var (
	testHash0 = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
	testHash1 = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
	testHash2 = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
	testHash3 = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
)

type testContent struct {
	common.Hash
}

func(tc testContent) CalculateHash() ([]byte, error) {
	return tc.Bytes(), nil
}

func(tc testContent) Equals(other Content) (bool, error)  {
	h1, err := tc.CalculateHash()
	if err != nil {
		return false, err
	}
	h2, err := other.CalculateHash()
	if err != nil {
		return false, err
	}
	return bytes.Equal(h1, h2), nil
}

func TestNewTree(t *testing.T) {
	testContents := []testContent{
		{testHash0},
		{testHash1},
		{testHash2},
		{testHash3},
	}
	root, err := NewTree([]Content{testContents[0], testContents[1]})
	if err != nil {
		t.Error("new tree error: ", err)
	}
	rootHash := common.BytesToHash(root.Root.Hash).Hex()
	want := "0x90f4b39548df55ad6187a1d20d731ecee78c545b94afd16f42ef7592d99cd365"
	if rootHash != want {
		t.Errorf("root unmatched. should be %s, got %s", want, rootHash)
	}
}
