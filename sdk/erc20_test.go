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
	var (
		price uint64 = 0
		limit uint64 = 0

		// rpc             string = "https://bsc-dataseed1.defibit.io/"
		// contractAddress string = "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c" // WBNB
		// fromAddr        string = "0x14bc30855e76Ba7e83d73BAb362C5cdc79EF2AF3"
		// owner           string = "0x676c44240bab1f23f31c591fc422d48a9ec9de30"

		rpc             string = "http://sh.gitez.cc:8547"
		contractAddress string = "0x5ab9F95fDe4A43689461241485C8eE55F9DC85DE" // TRT
		fromAddr        string = "0xF884c247f1EeD69f8DFa618fB4CAcE8EEb47C91F"
		owner           string = "0xF884c247f1EeD69f8DFa618fB4CAcE8EEb47C91F"
	)
	tm, err := New(rpc, price, limit, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	ret, err := tm.BalanceOf(contractAddress, fromAddr, owner)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ret: %x", ret)

	bi := big.NewInt(0)
	bi = bi.SetBytes(ret)
	t.Logf("ret: %v",bi.String())

}
