// Package sdk
// @Project:       eth
// @File:          erc20_test.go
// @Author:        eagle
// @Create:        2021/06/17 10:13:37
// @Description:
package sdk

import (
	"math/big"
	"testing"
)

func TestTransactionManager_TotalSupply(t *testing.T) {
}

func TestBalanceOf(t *testing.T) {
	abiStr := `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_amount","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"totalSupply","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[],"name":"destroy","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"inputs":[],"payable":false,"type":"constructor","stateMutability":"nonpayable"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"_spender","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Approval","type":"event"}]`
	var (
		price uint64 = 0
		limit uint64 = 0

		// rpc string = "https://bsc-dataseed1.defibit.io/"
		// contractAddress string = "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c" // WBNB
		// addr0        string = "0x14bc30855e76Ba7e83d73BAb362C5cdc79EF2AF3"

		rpc             string = "http://sh.gitez.cc:8547"
		contractAddress string = "0x5ab9F95fDe4A43689461241485C8eE55F9DC85DE" // TRT
		sk0             string = "6a139aa3de139e7b744fb49f684d77144d4d3476368dce463895c596645c423b"
		addr0           string = "0xF884c247f1EeD69f8DFa618fB4CAcE8EEb47C91F"
		spender         string = "0x96a8fc39cea5e5f1a1ea2090bd40de70ffa88747"
		spenderSk       string = "019a583104ce1f0bcf2f20e647c0fc268ff878ea110d0374a7c351b4a4ca54f1"
		to              string = "0x14bc30855e76Ba7e83d73BAb362C5cdc79EF2AF3"
	)
	tm, err := New(rpc, price, limit, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	ret, err := tm.BalanceOf(contractAddress, addr0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ret: %x", ret)
	ret2, err := Unpack(abiStr, "balanceOf", ret)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("balanceOf ret2: %v", ret2)

	bi := big.NewInt(0)
	bi = bi.SetBytes(ret)
	t.Logf("ret: %v", bi.String())

	ret, err = tm.Symbol(contractAddress)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("symbol: %v", string(ret))
	ret3, err := Unpack(abiStr, "symbol", ret)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("symbol: %v", ret3)

	tm.SetChainID("20")
	hash, _, err := tm.ApproveSync(contractAddress, sk0, spender, "100", 10000000, 0, 1000000)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("approve hash: %v", hash)

	allowance, err := tm.Allowance(contractAddress, addr0, spender)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("allowlance: %x", allowance)

	hash, _, err = tm.TransferFromSync(contractAddress, spenderSk, addr0, to, "100", 10000000, 0, 1000000)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("transfer from hash: %v", hash)

	allowance, err = tm.Allowance(contractAddress, addr0, spender)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("allowlance: %x", allowance)
}
