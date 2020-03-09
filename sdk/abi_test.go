package sdk

import (
	"io/ioutil"
	"testing"
)

type methodAndArgs struct {
	Method string
	Args   string
}

func TestPack(t *testing.T) {
	bs, err := ioutil.ReadFile("abi.txt")
	if err != nil {
		t.Errorf("read abi.txt error: %v", err)
	}
	abiStr := string(bs)

	methods := []methodAndArgs{
		{"totalSupply", ""},
		{"balanceOf", "address:0xd69cfc58b5a8b3b7866d2c2682ba971074a946a0"},
		{"transfer", "address:0xd69cfc58b5a8b3b7866d2c2682ba971074a946a0;uint256:3;"},
	}
	for _, method := range methods {
		packedBytes, err := Pack(abiStr, method.Method, method.Args)
		if err != nil {
			t.Errorf("Pack method: %v error: %v", method.Method, err)
			break
		}
		t.Logf("method: '%v', packedBytes: '%x'", method.Method, packedBytes)
	}
}
