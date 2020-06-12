package types

import (
	"bytes"
	"github.com/CortexFoundation/CortexTheseus/common"
	"testing"
)

var (
	testContents = []testContent{
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000006")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000007")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000008")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000009")},
	}
)

type testContent struct {
	common.Hash
}

func (tc testContent) CalculateHash() ([]byte, error) {
	return tc.Bytes(), nil
}

func (tc testContent) Equals(other Content) (bool, error) {
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

func TestMerkleTree_AddNode(t *testing.T) {
	//root_rebuild, err := NewTree([]Content{testContents[0]})
	//if err != nil {
	//	t.Fatal("new tree error: ", err)
	//}

	root_add, err := NewTree([]Content{testContents[0]})
	if err != nil {
		t.Fatal("new tree error: ", err)
	}
	//rebuildHash := common.BytesToHash(root_rebuild.Root.Hash).Hex()
	//addHash := common.BytesToHash(root_add.Root.Hash).Hex()
	//if rebuildHash != addHash {
	//	t.Fatalf("root unmatched at %d. rebuild hash is %s, add hash is %s", 0, rebuildHash, addHash)
	//}
	//t.Log(root_rebuild.String())
	//t.Log(root_add.String())

	for i := 1; i < 2; i += 1 {
		c := testContents[i]
		//root_rebuild.Leafs = append(root_rebuild.Leafs, &Node{
		//	C:      c,
		//})
		//root_rebuild.RebuildTree()

		root_add.AddNode(c)

		if v, err := root_add.VerifyTree(); !v || err != nil {
			t.Errorf("root add not verified, at %d", i)
		}
		//rebuildHash := common.BytesToHash(root_rebuild.Root.Hash).Hex()
		//addHash := common.BytesToHash(root_add.Root.Hash).Hex()
		//if rebuildHash != addHash {
		//t.Log(root_rebuild.String())
		//t.Log(root_add.String())
		//t.Fatalf("root unmatched at %d. rebuild hash is %s, add hash is %s", i, rebuildHash, addHash)
		//}
	}
}
