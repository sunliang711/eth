## eth

### Pack

> Pack can encode contract method arguments to abi format

```go
package main

import (
	"fmt"
	"os"

	ethSdk "github.com/sunliang711/eth/sdk"
)

type methodAndArgs struct {
	Method string
	Args   string
}

func main() {
	abiStr := `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_amount","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"totalSupply","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[],"name":"destroy","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function","stateMutability":"view"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function","stateMutability":"nonpayable"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"type":"function","stateMutability":"view"},{"inputs":[],"payable":false,"type":"constructor","stateMutability":"nonpayable"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"_spender","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Approval","type":"event"}]`
	methods := []methodAndArgs{
		{"totalSupply", ""},
		{"balanceOf", "address:0xd69cfc58b5a8b3b7866d2c2682ba971074a946a0"},
		{"transfer", "address:0xd69cfc58b5a8b3b7866d2c2682ba971074a946a0;uint256:3;"},
	}
	for _, method := range methods {
		packedBytes, err := ethSdk.Pack(abiStr, method.Method, method.Args)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Pack method: %v error: %v", method.Method, err)
			break
		}
		fmt.Printf("method: '%v', packedBytes: '%x'", method.Method, packedBytes)
	}
	// not, packedBytes can used as data of transaction
}
```

### create contract

```go
package main

import (
	"fmt"
	"os"

	ethSdk "github.com/sunliang711/eth/sdk"
)

func main() {
	var (
		rpc      string
		price    uint64
		limit    uint64
		timeout  uint64
		interval uint64
	)
	txManager, err := ethSdk.New(rpc, price, limit, timeout, interval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new transaction manager error: %v", err)
		os.Exit(1)
	}

	var (
		sk       string
		bytecode []byte
	)

	address, hash, gasUsed, err := txManager.CreateContractSync(sk, bytecode, 0, 0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create contract error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("contract created at %v with hash: %v gas used: %v\n", address, hash, gasUsed)
}
```

### call contract method (write and read)

```go
package main

import (
	"fmt"
	"math/big"
	"os"

	ethSdk "github.com/sunliang711/eth/sdk"
)

func main() {
	var (
		rpc      string
		price    uint64
		limit    uint64
		timeout  uint64
		interval uint64
	)
	txManager, err := ethSdk.New(rpc, price, limit, timeout, interval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new transaction manager error: %v", err)
		os.Exit(1)
	}

	var (
		sk              string
		contractAddress string
		v               *big.Int
		abi             string
		methodName      string
		readMethodName  string
		args            string
	)

	hash, gasUsed, err := txManager.WriteContractSync(sk, contractAddress, v, abi, methodName, args, 0, 0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "call method: %v error: %v", methodName, err)
		os.Exit(1)
	}
	fmt.Printf("hash: %v gasUsed: %v\n", hash, gasUsed)

	// read contract
	output, err := txManager.ReadContract(contractAddress, abi, readMethodName, args, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "call method: %v error: %v", readMethodName, err)
		os.Exit(1)
	}
	fmt.Printf("read contract result: %v\n", output)

}
```

### ERC20 support

```go

package main

import (
	"fmt"
	"os"

	ethSdk "github.com/sunliang711/eth/sdk"
)

func main() {
	var (
		rpc      string
		timeout  uint64
		interval uint64
		price    uint64
		limit    uint64
	)
	txManager, err := ethSdk.New(rpc, price, limit, timeout, interval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new transaction manager error: %v", err)
		os.Exit(1)
	}

	var (
		contractAddress string
		sk0             string
		addr0           string
		spender         string
		spenderSk       string
		to              string
	)
	balance, err := txManager.BalanceOf(contractAddress, addr0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get balance of %v at contract %v error: %v\n", addr0, contractAddress, err)
		os.Exit(1)
	}
	fmt.Printf("balanceOf %v is: %v\n", addr0, balance)

	symbol, err := txManager.Symbol(contractAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get symbol of contract %v error: %v\n", contractAddress, err)
		os.Exit(1)
	}
	fmt.Printf("symbol of contract %v is: %v\n", contractAddress, symbol)

	hash, err := txManager.Approve(contractAddress, sk0, spender, "100", price, 0, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "approve error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("approve tx hash: %v\n", hash)

	hash, err = txManager.TransferFrom(contractAddress, spenderSk, addr0, to, "100", price, 0, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "TransferFrom error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("transfer from tx hash: %v\n", hash)

}
```