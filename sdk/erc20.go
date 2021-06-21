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

	MethodSymbol       = "symbol"
	MethodTotalSupply  = "totalSupply"
	MethodBalanceOf    = "balanceOf"
	MethodTransfer     = "transfer"
	MethodApprove      = "approve"
	MethodTransferFrom = "transferFrom"
	MethodAllowance    = "allowance"
)

// Symbol ERC20 symbol
func (tm *TransactionManager) Symbol(contractAddress string) ([]byte, error) {
	return tm.ReadContract(contractAddress, ERC20_ABI, MethodSymbol, "", nil)
}

// TotalSupply ERC20 totalSupply
func (tm *TransactionManager) TotalSupply(contractAddress string) ([]byte, error) {
	return tm.ReadContract(contractAddress, ERC20_ABI, MethodTotalSupply, "", nil)
}

// BalanceOf ERC20 balanceOf
func (tm *TransactionManager) BalanceOf(contractAddress string, owner string) ([]byte, error) {
	args := fmt.Sprintf("address:%v", owner)
	return tm.ReadContract(contractAddress, ERC20_ABI, MethodBalanceOf, args, nil)
}

// Transfer ERC20 transfer
func (tm *TransactionManager) Transfer(contractAddress string, sk string, to string, value string, price uint64, nonce uint64, limit uint64) (string, error) {
	args := fmt.Sprintf("address:%v;uint256:%v", to, value)
	return tm.WriteContract(sk, contractAddress, nil, ERC20_ABI, MethodTransfer, args, price, nonce, limit)
}

// Transfer ERC20 transfer sync
func (tm *TransactionManager) TransferSync(contractAddress string, sk string, to string, value string, price uint64, nonce uint64, limit uint64) (string, uint64, error) {
	args := fmt.Sprintf("address:%v;uint256:%v", to, value)
	return tm.WriteContractSync(sk, contractAddress, nil, ERC20_ABI, MethodTransfer, args, price, nonce, limit)
}

// Approve ERC20 approve
func (tm *TransactionManager) Approve(contractAddress string, sk string, spender string, value string, price uint64, nonce uint64, limit uint64) (string, error) {
	args := fmt.Sprintf("address:%v;uint256:%v", spender, value)
	return tm.WriteContract(sk, contractAddress, nil, ERC20_ABI, MethodApprove, args, price, nonce, limit)
}

// Approve ERC20 approve sync
func (tm *TransactionManager) ApproveSync(contractAddress string, sk string, spender string, value string, price uint64, nonce uint64, limit uint64) (string, uint64, error) {
	args := fmt.Sprintf("address:%v;uint256:%v", spender, value)
	return tm.WriteContractSync(sk, contractAddress, nil, ERC20_ABI, MethodApprove, args, price, nonce, limit)
}

// TransferFrom ERC20 transferFrom
func (tm *TransactionManager) TransferFrom(contractAddress string, sk string, from string, to string, value string, price uint64, nonce uint64, limit uint64) (string, error) {
	args := fmt.Sprintf("address:%v;address:%v;uint256:%v", from, to, value)
	return tm.WriteContract(sk, contractAddress, nil, ERC20_ABI, MethodTransferFrom, args, price, nonce, limit)
}

// TransferFrom ERC20 transferFrom sync
func (tm *TransactionManager) TransferFromSync(contractAddress string, sk string, from string, to string, value string, price uint64, nonce uint64, limit uint64) (string, uint64, error) {
	args := fmt.Sprintf("address:%v;address:%v;uint256:%v", from, to, value)
	return tm.WriteContractSync(sk, contractAddress, nil, ERC20_ABI, MethodTransferFrom, args, price, nonce, limit)
}

// Allowance ERC20 allowance
func (tm *TransactionManager) Allowance(contractAddress string, owner string, spender string) ([]byte, error) {
	args := fmt.Sprintf("address:%v;address:%v", owner, spender)
	return tm.ReadContract(contractAddress, ERC20_ABI, MethodAllowance, args, nil)
}

