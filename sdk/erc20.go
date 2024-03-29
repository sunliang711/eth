// Package sdk
// @Project:       eth
// @File:          erc20.go
// @Author:        eagle
// @Create:        2021/06/17 09:52:56
// @Description:
package sdk

import (
	"errors"
	"fmt"
	"math/big"
)

const (
	ERC20_ABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_amount","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"totalSupply","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[],"name":"destroy","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"inputs":[],"payable":false,"type":"constructor","stateMutability":"nonpayable"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Transfer20","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"_spender","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Approval","type":"event"}]`

	MethodSymbol       = "symbol"
	MethodDecimals     = "decimals"
	MethodTotalSupply  = "totalSupply"
	MethodBalanceOf    = "balanceOf"
	MethodTransfer     = "transfer"
	MethodApprove      = "approve"
	MethodTransferFrom = "transferFrom"
	MethodAllowance    = "allowance"
)

// Symbol20 ERC20 symbol
func (tm *TransactionManager) Symbol20(contractAddress string) (string, error) {
	raw, err := tm.ReadContract(contractAddress, ERC20_ABI, MethodSymbol, "", nil)
	if err != nil {
		return "", err
	}
	data, err := Unpack(ERC20_ABI, MethodSymbol, raw)
	if err != nil {
		return "", err
	}
	if len(data) == 1 {
		if ret, ok := data[0].(string); ok {
			return ret, nil
		} else {
			return "", errors.New("returned data not string")
		}
	}
	return "", errors.New("symbol() returned data error")
}

// TotalSupply20 ERC20 totalSupply
func (tm *TransactionManager) TotalSupply20(contractAddress string) (*big.Int, error) {
	raw, err := tm.ReadContract(contractAddress, ERC20_ABI, MethodTotalSupply, "", nil)
	if err != nil {
		return nil, err
	}
	data, err := Unpack(ERC20_ABI, MethodTotalSupply, raw)
	if err != nil {
		return nil, err
	}
	if len(data) == 1 {
		if ret, ok := data[0].(*big.Int); ok {
			return ret, nil
		} else {
			return nil, errors.New("returned data not big.Int")
		}
	}
	return nil, errors.New("totalSupply() returned data error")
}

// BalanceOf20 ERC20 balanceOf
func (tm *TransactionManager) BalanceOf20(contractAddress string, owner string) (*big.Int, error) {
	args := fmt.Sprintf("address:%v", owner)
	raw, err := tm.ReadContract(contractAddress, ERC20_ABI, MethodBalanceOf, args, nil)
	if err != nil {
		return nil, err
	}

	data, err := Unpack(ERC20_ABI, MethodBalanceOf, raw)
	if err != nil {
		return nil, err
	}
	if len(data) == 1 {
		if ret, ok := data[0].(*big.Int); ok {
			return ret, nil
		} else {
			return nil, errors.New("returned data not big.Int")
		}
	}
	return nil, errors.New("balanceOf() returned data error")
}

// Transfer20 ERC20 transfer
func (tm *TransactionManager) Transfer20(contractAddress string, sk string, to string, value string, price uint64, nonce uint64, limit uint64) (string, error) {
	args := fmt.Sprintf("address:%v;uint256:%v", to, value)
	return tm.WriteContract(sk, contractAddress, nil, ERC20_ABI, MethodTransfer, args, price, nonce, limit)
}

// Transfer20 ERC20 transfer sync
func (tm *TransactionManager) TransferSync20(contractAddress string, sk string, to string, value string, price uint64, nonce uint64, limit uint64) (string, uint64, error) {
	args := fmt.Sprintf("address:%v;uint256:%v", to, value)
	return tm.WriteContractSync(sk, contractAddress, nil, ERC20_ABI, MethodTransfer, args, price, nonce, limit)
}

// Approve20 ERC20 approve
func (tm *TransactionManager) Approve20(contractAddress string, sk string, spender string, value string, price uint64, nonce uint64, limit uint64) (string, error) {
	args := fmt.Sprintf("address:%v;uint256:%v", spender, value)
	return tm.WriteContract(sk, contractAddress, nil, ERC20_ABI, MethodApprove, args, price, nonce, limit)
}

// Approve20 ERC20 approve sync
func (tm *TransactionManager) ApproveSync20(contractAddress string, sk string, spender string, value string, price uint64, nonce uint64, limit uint64) (string, uint64, error) {
	args := fmt.Sprintf("address:%v;uint256:%v", spender, value)
	return tm.WriteContractSync(sk, contractAddress, nil, ERC20_ABI, MethodApprove, args, price, nonce, limit)
}

// TransferFrom20 ERC20 transferFrom
func (tm *TransactionManager) TransferFrom20(contractAddress string, sk string, from string, to string, value string, price uint64, nonce uint64, limit uint64) (string, error) {
	args := fmt.Sprintf("address:%v;address:%v;uint256:%v", from, to, value)
	return tm.WriteContract(sk, contractAddress, nil, ERC20_ABI, MethodTransferFrom, args, price, nonce, limit)
}

// TransferFrom20 ERC20 transferFrom sync
func (tm *TransactionManager) TransferFromSync20(contractAddress string, sk string, from string, to string, value string, price uint64, nonce uint64, limit uint64) (string, uint64, error) {
	args := fmt.Sprintf("address:%v;address:%v;uint256:%v", from, to, value)
	return tm.WriteContractSync(sk, contractAddress, nil, ERC20_ABI, MethodTransferFrom, args, price, nonce, limit)
}

// Allowance20 ERC20 allowance
func (tm *TransactionManager) Allowance20(contractAddress string, owner string, spender string) (*big.Int, error) {
	args := fmt.Sprintf("address:%v;address:%v", owner, spender)
	raw, err := tm.ReadContract(contractAddress, ERC20_ABI, MethodAllowance, args, nil)
	if err != nil {
		return nil, err
	}

	data, err := Unpack(ERC20_ABI, MethodAllowance, raw)
	if err != nil {
		return nil, err
	}

	if len(data) == 1 {
		if ret, ok := data[0].(*big.Int); ok {
			return ret, nil
		} else {
			return nil, errors.New("returned data not big.Int")
		}
	}
	return nil, errors.New("allowance() returned data error")
}

// Decimals20 ERC20 decimals
func (tm *TransactionManager) Decimals20(contractAddress string) (uint8, error) {
	raw, err := tm.ReadContract(contractAddress, ERC20_ABI, MethodDecimals, "", nil)
	if err != nil {
		return 0, err
	}

	data, err := Unpack(ERC20_ABI, MethodDecimals, raw)
	if err != nil {
		return 0, err
	}
	if len(data) == 1 {
		if ret, ok := data[0].(uint8); ok {
			return ret, nil
		} else {
			return 0, errors.New("returned data not uint8")
		}
	}
	return 0, errors.New("decimals() returned data error")
}
