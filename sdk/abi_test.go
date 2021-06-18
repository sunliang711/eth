package sdk

import (
	"testing"
)

type methodAndArgs struct {
	Method string
	Args   string
}

func TestPack(t *testing.T) {
	abiStr := `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_amount","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"totalSupply","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[],"name":"destroy","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"inputs":[],"payable":false,"type":"constructor","stateMutability":"nonpayable"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"_spender","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Approval","type":"event"}]`

	methods := []methodAndArgs{
		{"totalSupply", ""},
		{"balanceOf", "address:0xd69cfc58b5a8b3b7866d2c2682ba971074a946a0"},
		{"transfer", "address:0xd69cfc58b5a8b3b7866d2c2682ba971074a946a0;uint256:3;"},
		// // function swapETHForExactTokens(uint256 amountOut, address[] calldata path, address to, uint256 deadline)
		// {"swapETHForExactTokens", "uint256:582842921510003747;address[]:0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c,0xb1035523a844371c2877f8a3b2f2f8d337403b6f;address:0x5411bbb30d8ad59b63057f4816761544811b66bf;uint256:1623824109;"},
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
