// Package sdk
// @Project:       eth
// @File:          contract.go.go
// @Author:        eagle
// @Create:        2021/07/15 14:00:44
// @Description:
package sdk

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func (tm *TransactionManager) SendCallMsgTx(toAddr string, data []byte, blockNumber *big.Int) ([]byte, error) {
	to := common.HexToAddress(toAddr)

	msg := ethereum.CallMsg{
		To:   &to,
		Data: data,
	}
	return tm.Client.CallContract(context.Background(), msg, blockNumber)
}

// CreateContract creates a contract,return tx's hash,use it to query contract address
func (tm *TransactionManager) CreateContract(sk string, data []byte, gasPrice uint64, nonce uint64, gasLimit uint64) (string, error) {
	return tm.SendTx(sk, "", nil, data, gasPrice, nonce, gasLimit)
}

// CreateContractSync creates a contract syncly, return contract address ,tx hash ,gas used, error
// set timeout to 0 to use default timeout value
func (tm *TransactionManager) CreateContractSync(sk string, data []byte, gasPrice uint64, nonce uint64, gasLimit uint64) (string, string, uint64, error) {
	hash, err := tm.CreateContract(sk, data, gasPrice, nonce, gasLimit)
	if err != nil {
		return "", "", 0, err
	}

	deadline := time.After(time.Second * time.Duration(tm.timeout))
	tick := time.Tick(time.Second * time.Duration(tm.interval))
	for {
		select {
		case <-deadline:
			return "", "", 0, fmt.Errorf("timeout")
		case <-tick:
			receipt, err := tm.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
			if err != nil {
				// skip
			} else {
				return receipt.ContractAddress.String(), hash, receipt.GasUsed, nil
			}
		}
	}

}

func (tm *TransactionManager) GetContractAddress(hash string) (string, error) {
	receipt, err := tm.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		return "", err
	}
	return receipt.ContractAddress.String(), nil
}

func (tm *TransactionManager) GetContractAddressSync(hash string) (string, error) {
	deadline := time.After(time.Second * time.Duration(tm.timeout))
	tick := time.Tick(time.Second * time.Duration(tm.interval))

	for {
		select {
		case <-deadline:
			return "", fmt.Errorf("timeout")
		case <-tick:
			receipt, err := tm.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
			if err != nil {
				// skip
			} else {
				return receipt.ContractAddress.String(), nil
			}
		}
	}

}

// WriteContract sends an async write contract,return hash,error
func (tm *TransactionManager) WriteContract(sk string, contractAddress string, v *big.Int, abi string, methodName, args string, gasPrice uint64, nonce uint64, gasLimit uint64) (string, error) {
	payload, err := Pack(abi, methodName, args)
	if err != nil {
		return "", err
	}
	hash, err := tm.SendTx(sk, contractAddress, v, payload, gasPrice, nonce, gasLimit)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// WriteContractSync sends an sync write contract,return hash, gas used, error
func (tm *TransactionManager) WriteContractSync(sk string, contractAddress string, v *big.Int, abi string, methodName, args string, gasPrice uint64, nonce uint64, gasLimit uint64) (string, uint64, error) {
	hash, err := tm.WriteContract(sk, contractAddress, v, abi, methodName, args, gasPrice, nonce, gasLimit)
	if err != nil {
		return "", 0, err
	}

	deadline := time.After(time.Second * time.Duration(tm.timeout))
	tick := time.Tick(time.Second * time.Duration(tm.interval))

	for {
		select {
		case <-deadline:
			return "", 0, fmt.Errorf("timeout")
		case <-tick:
			receipt, err := tm.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
			if err != nil {
				// skip
			} else {
				return hash, receipt.GasUsed, nil
			}
		}
	}
}

// ReadContract send a call msg tx to contract, set blockNumber to nil for latest block
func (tm *TransactionManager) ReadContract(contractAddress string, abi string, methodName, args string, blockNumber *big.Int) ([]byte, error) {
	payload, err := Pack(abi, methodName, args)
	if err != nil {
		return nil, err
	}
	output, err := tm.SendCallMsgTx(contractAddress, payload, blockNumber)
	if err != nil {
		return nil, err
	}
	return output, nil
}

