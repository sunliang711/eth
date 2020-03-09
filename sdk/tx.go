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
	"github.com/sirupsen/logrus"
)

const (
	// 5 seconds timeout with dial
	dialTimeout = 5
	// gasLimit    = 2.3e4
)

func dial(rpcURL string) (*ethclient.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(dialTimeout))
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// SendTx sends generic transaction to chain, including 'transfer eth' 'create contract'(when toAddr is empty string) 'call contract',not including 'call msg tx'
// set gasPrice to 0 to use suggest gas price
func SendTx(rpcURL string, fromSK string, toAddr string, value uint64, data []byte, gasPrice uint64, gasLimit uint64) (common.Hash, error) {
	client, err := dial(rpcURL)
	if err != nil {
		return common.Hash{}, err
	}
	defer client.Close()

	privK, _, fromAddress, err := HexToAccount(fromSK)
	if err != nil {
		logrus.Errorf("convert hex sk to ECDSA error: %v", err)
		return common.Hash{}, err
	}
	var toAddress common.Address
	if toAddr != "" {
		toAddress = common.HexToAddress(toAddr)
	}
	//nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("PendingNonceAt() error: %s", err.Error())
	}

	//gas price
	if gasPrice == 0 {
		suggestGasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			return common.Hash{}, fmt.Errorf("SuggestGasPrice() error: %s", err.Error())
		}
		gasPrice = suggestGasPrice.Uint64()
	}
	price := new(big.Int)
	price.SetUint64(gasPrice)

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
	return signedTx.Hash(), client.SendTransaction(context.Background(), signedTx)
}

// SendCallMsgTx calls readonly contract functions on local node
// set gasPrice to 0 to use suggest gas price
func SendCallMsgTx(rpcURL string, fromAddr string, toAddr string, data []byte, gasPrice uint64, gasLimit uint64) ([]byte, error) {
	client, err := dial(rpcURL)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	from := common.HexToAddress(fromAddr)
	to := common.HexToAddress(toAddr)

	v := new(big.Int)
	v.SetUint64(uint64(0))

	//gas price
	if gasPrice == 0 {
		suggestGasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("SuggestGasPrice() error: %s", err.Error())
		}
		gasPrice = suggestGasPrice.Uint64()
	}
	price := new(big.Int)
	price.SetUint64(gasPrice)

	msg := ethereum.CallMsg{
		From:     from,
		To:       &to,
		Gas:      gasLimit,
		GasPrice: price,
		Value:    v,
		Data:     data,
	}
	return client.CallContract(context.Background(), msg, nil)
}

// GetBalance query balance of address
func GetBalance(rpcURL string, address string) (*big.Int, error) {
	client, err := dial(rpcURL)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	addr := common.HexToAddress(address)
	return client.BalanceAt(context.Background(), addr, nil)
}

// CreateContract create smart contract with sk and data
// set gasPrice to 0 to use suggest gas price
func CreateContract(rpcURL string, sk string, data []byte, gasPrice uint64, gasLimit uint64) (common.Hash, error) {
	return SendTx(rpcURL, sk, "", 0, data, gasPrice, gasLimit)
}

// GetTransactionUsedGas get used gas of a transaction
func GetTransactionUsedGas(rpcURL string, hash string) (uint64, error) {
	client, err := dial(rpcURL)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	h := common.HexToHash(hash)
	receipt, err := client.TransactionReceipt(context.Background(), h)
	if err != nil {
		return 0, err
	}
	return receipt.GasUsed, nil
}

// GetContractAddress query address of contract by it's hash
func GetContractAddress(rpcURL string, hash string) (string, error) {
	client, err := dial(rpcURL)
	if err != nil {
		return "", err
	}
	defer client.Close()

	h := common.HexToHash(hash)
	receipt, err := client.TransactionReceipt(context.Background(), h)
	if err != nil {
		return "", err
	}
	return receipt.ContractAddress.String(), nil
}

// WriteContract calls writable function of contract
// set gasPrice to 0 to use suggest gas price
func WriteContract(rpcURL string, sk string, contractAddress string, abi string, methodName, args string, gasPrice uint64, gasLimit uint64) (string, error) {
	payload, err := Pack(abi, methodName, args)
	if err != nil {
		return "", err
	}
	hash, err := SendTx(rpcURL, sk, contractAddress, 0, payload, gasPrice, gasLimit)
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}

// ReadContract calls readonly function of contract
// set gasPrice to 0 to use suggest gas price
func ReadContract(rpcURL string, fromAddr string, contractAddress string, abi string, methodName, args string, gasPrice uint64, gasLimit uint64) ([]byte, error) {
	payload, err := Pack(abi, methodName, args)
	if err != nil {
		return nil, err
	}
	output, err := SendCallMsgTx(rpcURL, fromAddr, contractAddress, payload, gasPrice, gasLimit)
	if err != nil {
		return nil, err
	}
	return output, nil
}

//TODO decode output: use Unpack
