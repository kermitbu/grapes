package core

import "testing"

func TestRpcUnpack(t *testing.T) {
	rpchead := &RpcHead{}

	buf := []byte{1, 9, 0, 10, 4, 97, 98, 99, 100}

	rpchead.Unpack(buf)

	if rpchead.Version != 1 {
		t.Errorf("version error", rpchead)
	}
	if rpchead.BodyLen != 10 {
		t.Errorf("body error", rpchead)
	}
	if rpchead.HeadLen != 9 {
		t.Errorf("body error", rpchead)
	}
	if rpchead.FuncName != "abcd" {
		t.Errorf("funcname error", rpchead)
	}

	newbuf := rpchead.Pack()
	if len(buf) != len(newbuf) {
		t.Errorf("funcname error", rpchead)
	}
}
