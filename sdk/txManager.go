package sdk

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	defaultTimeout   = 120
	defaultInterval  = 2
	transferEthLimit = 2.1e4
)

// TransactionManager store info to operate tx
type TransactionManager struct {
	rpcURL   string
	gasPrice uint64
	timeout  uint64
	interval uint64
	*ethclient.Client
}

// New makes a new TransactionManager
/*
** if gasPrice is 0, use suggest gas price
** if timeout is 0, use default timeout
** if interval is 0, use default interval
 */
func New(rpcURL string, gasPrice, timeout, interval uint64) (*TransactionManager, error) {
	var err error
	tm := &TransactionManager{
		rpcURL:   rpcURL,
		gasPrice: gasPrice,
		timeout:  timeout,
		interval: interval,
	}
	tm.Client, err = dial(tm.rpcURL)
	if err != nil {
		return nil, err
	}
	if gasPrice == 0 {
		sgp, err := tm.Client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("Get suggest gas price error: %s", err.Error())
		}
		tm.gasPrice = sgp.Uint64()
	}

	if timeout == 0 {
		tm.timeout = defaultTimeout
	}
	if interval == 0 {
		tm.interval = defaultInterval
	}

	return tm, nil
}

func (tm *TransactionManager) Close() {
	tm.Client.Close()
}

func (tm *TransactionManager) GasPrice() uint64 {
	return tm.gasPrice
}

// SendTx sends an async tx, and return tx's hash
func (tm *TransactionManager) SendTx(fromSK string, toAddr string, value uint64, data []byte, gasLimit uint64) (string, error) {
	privK, _, fromAddress, err := HexToAccount(fromSK)
	if err != nil {
		return "", fmt.Errorf("convert hex sk to ECDSA error: %s", err.Error())
	}
	var toAddress common.Address
	if toAddr != "" {
		toAddress = common.HexToAddress(toAddr)
	}
	//nonce
	nonce, err := tm.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("PendingNonceAt() error: %s", err.Error())
	}

	price := new(big.Int)
	price.SetUint64(tm.gasPrice)

	// value
	val := new(big.Int)
	val.SetUint64(value)

	var tx *types.Transaction
	if toAddr != "" {
		tx = types.NewTransaction(nonce, toAddress, val, gasLimit, price, data)
	} else {
		tx = types.NewContractCreation(nonce, val, gasLimit, price, data)
	}
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privK)
	return signedTx.Hash().String(), tm.Client.SendTransaction(context.Background(), signedTx)
}

// SendTxSync sends an sync tx
// 调用者应该比较参数gasLimit和返回值的第二个gasUsed
// 如果gasUsed 等于 gasLimit
//    1. 如果这是个智能合约相关的操作(创建合约、写合约)，那么这个交易可能是部分完成，执行了部分指令, 用掉了gasLimit等量的gas，应该提高gasLimit上限重新调用一次
//    2. 如果这是个转账操作，那么执行时成功的（转账的gasLimit为固定值21000） TODO 转账时gasLimit小于21000会发生啥
func (tm *TransactionManager) SendTxSync(fromSK string, toAddr string, value uint64, data []byte, gasLimit uint64) (string, uint64, error) {
	hash, err := tm.SendTx(fromSK, toAddr, value, data, gasLimit)
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
				//ignore
			} else {
				return hash, receipt.GasUsed, nil
			}
		}
	}
}

func (tm *TransactionManager) SendCallMsgTx(fromAddr string, toAddr string, data []byte, gasLimit uint64) ([]byte, error) {
	from := common.HexToAddress(fromAddr)
	to := common.HexToAddress(toAddr)

	v := new(big.Int)
	v.SetUint64(uint64(0))

	price := new(big.Int)
	price.SetUint64(tm.gasPrice)

	msg := ethereum.CallMsg{
		From:     from,
		To:       &to,
		Gas:      gasLimit,
		GasPrice: price,
		Value:    v,
		Data:     data,
	}
	return tm.Client.CallContract(context.Background(), msg, nil)
}

// CreateContract creates a contract,return tx's hash,use it to query contract address
func (tm *TransactionManager) CreateContract(sk string, data []byte, gasLimit uint64) (string, error) {
	return tm.SendTx(sk, "", 0, data, gasLimit)
}

// CreateContractSync creates a contract syncly, return contract address ,tx hash ,gas used, error
// set timeout to 0 to use default timeout value
func (tm *TransactionManager) CreateContractSync(sk string, data []byte, gasLimit uint64) (string, string, uint64, error) {
	hash, err := tm.CreateContract(sk, data, gasLimit)
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
				//skip
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

// WriteContract sends an async write contract,return hash,error
func (tm *TransactionManager) WriteContract(sk string, contractAddress string, abi string, methodName, args string, gasLimit uint64) (string, error) {
	payload, err := Pack(abi, methodName, args)
	if err != nil {
		return "", err
	}
	hash, err := tm.SendTx(sk, contractAddress, 0, payload, gasLimit)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// WriteContract sends an sync write contract,return hash, gas used, error
func (tm *TransactionManager) WriteContractSync(sk string, contractAddress string, abi string, methodName, args string, gasLimit uint64) (string, uint64, error) {
	hash, err := tm.WriteContract(sk, contractAddress, abi, methodName, args, gasLimit)
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
				//skip
			} else {
				return hash, receipt.GasUsed, nil
			}
		}
	}
}

func (tm *TransactionManager) ReadContract(fromAddr string, contractAddress string, abi string, methodName, args string, gasLimit uint64) ([]byte, error) {
	payload, err := Pack(abi, methodName, args)
	if err != nil {
		return nil, err
	}
	output, err := tm.SendCallMsgTx(fromAddr, contractAddress, payload, gasLimit)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// TransferEth send an async eth-transfer tx
// return tx hash,error
func (tm *TransactionManager) TransferEth(fromSK string, toAddr string, value uint64) (string, error) {
	return tm.SendTx(fromSK, toAddr, value, nil, transferEthLimit)
}

// GetBalance query balance of 'address'
func (tm *TransactionManager) GetBalance(address string) (uint64, error) {
	balance, err := tm.Client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}
