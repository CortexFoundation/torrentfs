package merkletree

import (
//"golang.org/x/crypto/sha3"
)

type BlockContent struct {
	X string
	N uint64
}

func (t BlockContent) CalculateHash() ([]byte, error) {
	//h := sha3.NewLegacyKeccak256()
	//if _, err := h.Write([]byte(t.x)); err != nil {
	//	return nil, err
	//}

	//return h.Sum(nil), nil
	h := newHasher()
	defer returnHasherToPool(h)
	return h.sum([]byte(t.X)), nil
}

//Equals tests for equality of two Contents
func (t BlockContent) Equals(other Content) (bool, error) {
	return t.X == other.(BlockContent).X, nil
}
