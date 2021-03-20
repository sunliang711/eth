package sdk

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"
)

const (
	rpcURL2 = "http://localhost:8545"
	//contract owner
	sk1 = "6d6078a1f348b1c7f93b2b5dd1ac93a4d40e73c02b5f724b32dc5911daef34f8"
	addr1 = "f5426ae9197698ed77c04c4eca00b2ea3e1df00c"

	sk2   = "7766d545d9c1a22aabc3524990936ab1aa8e3269920b59f87eeb4d331e3b8b65"
	addr2 = "cd449a0cdb1c9b95a2bc2b531c565333e0c0bb0a"

	//CHANGE ME
	contractAddress = "0x102ac5F9A5362b572DbBa8d2C7C1175D125c3A43"

	createContractLimit = 2e6
	writeContractLimit  = 1e5
	readContractLimit   = 1e5
)

func TestCreateAccount(t *testing.T) {
	sk, pk, addr, err := GenAccount()
	if err != nil {
		t.Fatalf("gen account error: %s", err)
	}
	t.Logf("\nsk: %s\naddr: %s\npk: %s\n", hex.EncodeToString(sk), hex.EncodeToString(addr), hex.EncodeToString(pk))
}
func TestGetBalance2(t *testing.T) {
	txMan, err := New(rpcURL2, 0, createContractLimit, 0, 0)
	if err != nil {
		t.Fatalf("New txManager error: %s", err)
	}
	balance, err := txMan.GetBalance(addr2)
	if err != nil {
		t.Fatalf("get balance error: %s", err.Error())
	}
	t.Logf("balance: %d", balance)
}

type ByteCode struct {
	Object    string `json:"object"`
	Opcodes   string `json:"opcodes"`
	SourceMap string `json:"sourceMap"`
}

func TestTransferEth(t *testing.T){
	txMan,err := New(rpcURL2,0,21000,0,0)
	if err != nil{
		t.Fatalf("new txman error: %v",err)
	}
	fromSk :=sk1
	toAddr := addr2
	v1 := "12"

	value := big.NewInt(0)
	_,ok := value.SetString(v1,10)
	if !ok{
		t.Fatalf("set v1 error: %v",err)
	}
	
	
	gasPrice := uint64(10)

	txHash,err := txMan.TransferEth(fromSk,toAddr,value,gasPrice,0)
	if err != nil{
		t.Fatalf("transfer eth error: %v",err)
	}
	t.Logf("tx hash: %v",txHash)
}

func TestCreateContract(t *testing.T) {
	txMan, err := New(rpcURL2, 0, createContractLimit, 0, 0)
	if err != nil {
		t.Fatalf("New txManager error: %s", err)
	}
	defer txMan.Close()

	bytecodeBytes, err := ioutil.ReadFile("testData/bytecode.txt")
	if err != nil {
		t.Fatalf("read bytecode.txt error: %v", err)
	}
	var bc ByteCode
	err = json.Unmarshal(bytecodeBytes, &bc)
	if err != nil {
		t.Fatalf("Unmarshal bytecode.txt error: %v", err)
	}

	bytecode, err := hex.DecodeString(bc.Object)
	if err != nil {
		t.Fatalf("decode bytecode error: %v", err)
	}

	address, hash, gasUsed, err := txMan.CreateContractSync(sk1, bytecode, 0, 0, 0)
	if err != nil {
		t.Fatalf("create contract error: %s", err)
	}
	t.Logf("contract address: %s,hash: %s, gasUsed: %d", address, hash, gasUsed)

}
func TestWriteContract(t *testing.T) {
	txMan, err := New(rpcURL2, 0, writeContractLimit, 0, 0)
	if err != nil {
		t.Fatalf("New txManager error: %s", err)
	}
	defer txMan.Close()

	abiContent, err := ioutil.ReadFile("testData/abi.txt")
	if err != nil {
		t.Fatalf("read abi.txt error: %v", err)
	}
	abiStr := string(abiContent)
	args := fmt.Sprintf("address:%v;uint256:1;", addr2)
	t.Logf("args: %v", args)
	hash, gasUsed, err := txMan.WriteContractSync(sk1, contractAddress, abiStr, "transfer", args, 0, 0, 0)
	if err != nil {
		t.Fatalf("write contract error: %s", err.Error())
	}
	t.Logf("hash: %s\ngas used: %d\n", hash, gasUsed)
}

func TestReadContract(t *testing.T) {
	txMan, err := New(rpcURL2, 0, readContractLimit, 0, 0)
	if err != nil {
		t.Fatalf("New txManager error: %s", err)
	}
	defer txMan.Close()

	abiContent, err := ioutil.ReadFile("testData/abi.txt")
	if err != nil {
		t.Fatalf("read abi.txt error: %v", err)
	}
	abiStr := string(abiContent)
	args := fmt.Sprintf("address:%v;", addr2)
	t.Logf("args: %v", args)
	output, err := txMan.ReadContract(addr1, contractAddress, abiStr, "balanceOf", args, 0, 0)
	if err != nil {
		t.Fatalf("read contract error: %s", err.Error())
	}
	t.Logf("result: %s", hex.EncodeToString(output))

}
