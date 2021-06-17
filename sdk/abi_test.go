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
		t.Fatalf("read abi.txt error: %v", err)
	}
	abiStr := string(bs)

	methods := []methodAndArgs{
		// {"totalSupply", ""},
		// {"balanceOf", "address:0xd69cfc58b5a8b3b7866d2c2682ba971074a946a0"},
		// {"transfer", "address:0xd69cfc58b5a8b3b7866d2c2682ba971074a946a0;uint256:3;"},
		// function swapETHForExactTokens(uint256 amountOut, address[] calldata path, address to, uint256 deadline)
		{"swapETHForExactTokens","uint256:582842921510003747;address[]:0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c,0xb1035523a844371c2877f8a3b2f2f8d337403b6f;address:0x5411bbb30d8ad59b63057f4816761544811b66bf;uint256:1623824109;"},
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
