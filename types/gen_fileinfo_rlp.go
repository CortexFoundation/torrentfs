// Code generated by rlpgen. DO NOT EDIT.

package types

import "github.com/CortexFoundation/CortexTheseus/rlp"
import "io"

func (obj *FileInfo) EncodeRLP(_w io.Writer) error {
	w := rlp.NewEncoderBuffer(_w)
	_tmp0 := w.List()
	if err := obj.Meta.EncodeRLP(w); err != nil {
		return err
	}
	if obj.ContractAddr == nil {
		w.Write([]byte{0x80})
	} else {
		w.WriteBytes(obj.ContractAddr[:])
	}
	w.WriteUint64(obj.LeftSize)
	_tmp1 := w.List()
	for _, _tmp2 := range obj.Relate {
		w.WriteBytes(_tmp2[:])
	}
	w.ListEnd(_tmp1)
	w.ListEnd(_tmp0)
	return w.Flush()
}
