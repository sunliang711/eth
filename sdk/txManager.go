package sdk

import (
	"context"
	"fmt"
	"math/big"
	"time"

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
	gasPrice uint64 // default gas price
	gasLimit uint64 // default gas limit
	timeout  uint64
	interval uint64
	*ethclient.Client
	chainID string
	eip155  bool
}

// New makes a new TransactionManager
/*
** if gasPrice is 0, use suggest gas price
** if timeout is 0, use default timeout
** if interval is 0, use default interval
 */
func New(rpcURL string, gasPrice, gasLimit, timeout, interval uint64) (*TransactionManager, error) {
	var err error
	tm := &TransactionManager{
		rpcURL:   rpcURL,
		gasPrice: gasPrice,
		gasLimit: gasLimit,
		timeout:  timeout,
		interval: interval,
		eip155:   true,
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

	chainID, err := tm.Client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	tm.chainID = chainID.String()

	if timeout == 0 {
		tm.timeout = defaultTimeout
	}
	if interval == 0 {
		tm.interval = defaultInterval
	}

	return tm, nil
}

// // set chain id for EIP155
// func (tm *TransactionManager) SetChainID(id string) {
// 	tm.chainID = id
// }

func (tm *TransactionManager) Close() {
	tm.Client.Close()
}

func (tm *TransactionManager) GasPrice() uint64 {
	return tm.gasPrice
}

func (tm *TransactionManager) DisableEIP155() {
	tm.eip155 = false
}

// SendTx sends an async tx, and return tx's hash
// pass gasPrice to 0 to use tm.gasPrice
// pass nonce to 0 to use pendingNonce
// pass gasLimit to 0 to use tm.gasLimit
func (tm *TransactionManager) SendTx(fromSK string, toAddr string, value *big.Int, data []byte, gasPrice uint64, nonce uint64, gasLimit uint64) (string, error) {
	privK, _, fromAddress, err := HexToAccount(fromSK)
	if err != nil {
		return "", fmt.Errorf("convert hex sk to ECDSA error: %s", err.Error())
	}
	var toAddress common.Address
	if toAddr != "" {
		toAddress = common.HexToAddress(toAddr)
	}
	if nonce == 0 {
		nonce, err = tm.Client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return "", fmt.Errorf("PendingNonceAt() error: %s", err.Error())
		}
	}

	price := new(big.Int)
	if gasPrice == 0 {
		gasPrice = tm.gasPrice
	}
	price.SetUint64(gasPrice)

	if gasLimit == 0 {
		gasLimit = tm.gasLimit
	}

	// TODO go-ethereum: blockchain_test.go
	// use NewTx()
	var tx *types.Transaction
	if toAddr != "" {
		tx = types.NewTransaction(nonce, toAddress, value, gasLimit, price, data)
	} else {
		tx = types.NewContractCreation(nonce, value, gasLimit, price, data)
	}
	var signedTx *types.Transaction
	if tm.eip155 {
		chainId, success := big.NewInt(0).SetString(tm.chainID, 10)
		if !success {
			return "", fmt.Errorf("invalid chain id")
		}
		signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainId), privK)
	} else {
		signedTx, err = types.SignTx(tx, types.HomesteadSigner{}, privK)
	}
	return signedTx.Hash().String(), tm.Client.SendTransaction(context.Background(), signedTx)
}

// SendTxSync sends an sync tx
// 调用者应该比较参数gasLimit和返回值的第二个gasUsed
// 如果gasUsed 等于 gasLimit
//    1. 如果这是个智能合约相关的操作(创建合约、写合约)，那么这个交易可能是部分完成，执行了部分指令, 用掉了gasLimit等量的gas，应该提高gasLimit上限重新调用一次
//    2. 如果这是个转账操作，那么执行时成功的（转账的gasLimit为固定值21000） TODO 转账时gasLimit小于21000会发生啥
func (tm *TransactionManager) SendTxSync(fromSK string, toAddr string, value *big.Int, data []byte, gasPrice uint64, nonce uint64, gasLimit uint64) (string, uint64, error) {
	hash, err := tm.SendTx(fromSK, toAddr, value, data, gasPrice, nonce, gasLimit)
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
				// ignore
			} else {
				return hash, receipt.GasUsed, nil
			}
		}
	}
}

// TransferEth send an async eth-transfer tx
// return tx hash,error
func (tm *TransactionManager) TransferEth(fromSK string, toAddr string, value *big.Int, gasPrice uint64, nonce uint64) (string, error) {
	return tm.SendTx(fromSK, toAddr, value, nil, gasPrice, nonce, transferEthLimit)
}

// TransferEthSync send an sync eth-transfer tx
// return tx hash,error
func (tm *TransactionManager) TransferEthSync(fromSK string, toAddr string, value *big.Int, gasPrice uint64, nonce uint64) (string, error) {
	hash, err := tm.TransferEth(fromSK, toAddr, value, gasPrice, nonce)
	if err != nil {
		return "", err
	}

	deadline := time.After(time.Second * time.Duration(tm.timeout))
	tick := time.Tick(time.Second * time.Duration(tm.interval))
	for {
		select {
		case <-deadline:
			return "", fmt.Errorf("timeout")
		case <-tick:
			_, err := tm.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
			if err != nil {
				// skip
			} else {
				return hash, nil
			}
		}
	}
}

// TransferEthWithData send an async eth-transfer tx
// return tx hash,error
func (tm *TransactionManager) TransferEthWithData(fromSK string, toAddr string, value *big.Int, data []byte, gasPrice uint64, nonce uint64) (string, error) {
	return tm.SendTx(fromSK, toAddr, value, data, gasPrice, nonce, transferEthLimit)
}

// TransferEthWithDataSync send an sync eth-transfer tx
// return tx hash,error
func (tm *TransactionManager) TransferEthWithDataSync(fromSK string, toAddr string, value *big.Int, data []byte, gasPrice uint64, nonce uint64) (string, error) {
	hash, err := tm.TransferEthWithData(fromSK, toAddr, value, data, gasPrice, nonce)
	if err != nil {
		return "", err
	}

	deadline := time.After(time.Second * time.Duration(tm.timeout))
	tick := time.Tick(time.Second * time.Duration(tm.interval))
	for {
		select {
		case <-deadline:
			return "", fmt.Errorf("timeout")
		case <-tick:
			_, err := tm.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
			if err != nil {
				// skip
			} else {
				return hash, nil
			}
		}
	}
}

// GetBalance query balance of 'address'
func (tm *TransactionManager) GetBalance(address string) (*big.Int, error) {
	return tm.Client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
}

func (tm *TransactionManager) Receipt(hash string) (*types.Receipt, error) {
	return tm.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
}
