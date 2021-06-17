// Package sdk
// @Project:       eth
// @File:          erc20.go
// @Author:        eagle
// @Create:        2021/06/17 09:52:56
// @Description:
package sdk

import (
	"fmt"
)

const (
	ERC20_ABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_amount","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"totalSupply","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[],"name":"destroy","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"inputs":[],"payable":false,"type":"constructor","stateMutability":"nonpayable"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"_spender","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Approval","type":"event"}]`
)

func (tm *TransactionManager) BalanceOf(contractAddress string, fromAddr string, owner string) ([]byte, error) {
	args := fmt.Sprintf("address:%v", owner)
	return tm.ReadContract(fromAddr, contractAddress, ERC20_ABI, "balanceOf", args, 0, 0)
}

func (tm *TransactionManager) Transfer(contractAddress string, sk string, to string, value string, price uint64, nonce uint64, limit uint64) (string, error) {
	args := fmt.Sprintf("address:%v;uint256:%v", to, value)
	return tm.WriteContract(sk, contractAddress,nil, ERC20_ABI, "transfer", args, price, nonce, limit)
}

func (tm *TransactionManager) Approve() {
	panic("not implement")
}

func (tm *TransactionManager) TransferFrom() {
	panic("not implement")
}

func (tm *TransactionManager) TotalSupply(contractAddress string, fromAddr string) ([]byte, error) {
	return tm.ReadContract(fromAddr, contractAddress, ERC20_ABI, "totalSupply", "", 0, 0)
}
