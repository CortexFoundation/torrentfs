package merkletree

import (
	"bytes"
	"testing"
)

func TestContent(t *testing.T) {
	con1 := NewContent("qwerqwer", 2)
	p1,_ := con1.CalculateHash()
	con2 := NewContent("qwerqwer",4)
	p3, _ := con2.CalculateHash()
	if !bytes.Equal(p1, p3) {
		t.Errorf("should be he same")
	}

	if ok, err:=con1.Equals(con2);ok {
		t.Errorf("should not equel")
	}else if err != nil {
		t.Errorf("%v",err)
	}

}
