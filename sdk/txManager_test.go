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
	// contract owner
	sk1   = "6d6078a1f348b1c7f93b2b5dd1ac93a4d40e73c02b5f724b32dc5911daef34f8"
	addr1 = "f5426ae9197698ed77c04c4eca00b2ea3e1df00c"

	sk2   = "7766d545d9c1a22aabc3524990936ab1aa8e3269920b59f87eeb4d331e3b8b65"
	addr2 = "cd449a0cdb1c9b95a2bc2b531c565333e0c0bb0a"

	// CHANGE ME
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
	t.Logf("balance: %v", balance.String())
}

type ByteCode struct {
	Object    string `json:"object"`
	Opcodes   string `json:"opcodes"`
	SourceMap string `json:"sourceMap"`
}

func TestTransferEth(t *testing.T) {
	txMan, err := New(rpcURL2, 0, 21000, 0, 0)
	if err != nil {
		t.Fatalf("new txman error: %v", err)
	}
	fromSk := sk1
	toAddr := addr2
	v1 := "12"

	value := big.NewInt(0)
	_, ok := value.SetString(v1, 10)
	if !ok {
		t.Fatalf("set v1 error: %v", err)
	}

	gasPrice := uint64(10)

	txHash, err := txMan.TransferEth(fromSk, toAddr, value, gasPrice, 0)
	if err != nil {
		t.Fatalf("transfer eth error: %v", err)
	}
	t.Logf("tx hash: %v", txHash)
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
	hash, gasUsed, err := txMan.WriteContractSync(sk1, contractAddress, nil, abiStr, "transfer", args, 0, 0, 0)
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
	output, err := txMan.ReadContract(contractAddress, abiStr, "balanceOf", args, nil)
	if err != nil {
		t.Fatalf("read contract error: %s", err.Error())
	}
	t.Logf("result: %s", hex.EncodeToString(output))

}

func TestGetCollectionByID(t *testing.T) {
	rpc := "https://mainnet.infura.io/v3/f26e9265123241a4ba22cb9188089fe5"
	abi := `[{"constant":true,"inputs":[],"name":"currentStartingDigitalMediaId","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_interfaceID","type":"bytes4"}],"name":"supportsInterface","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_metadataPath","type":"string"}],"name":"createCollection","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_owner","type":"address"},{"name":"_totalSupply","type":"uint32"},{"name":"_digitalMediaMetadataPath","type":"string"},{"name":"_collectionMetadataPath","type":"string"},{"name":"_numReleases","type":"uint32"}],"name":"oboCreateDigitalMediaAndReleasesInNewCollection","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_tokenId","type":"uint256"}],"name":"getApproved","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_tokenId","type":"uint256"}],"name":"approve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"singleCreatorAddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"currentDigitalMediaStore","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_digitalMediaId","type":"uint256"}],"name":"burnDigitalMedia","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_index","type":"uint256"}],"name":"tokenOfOwnerByIndex","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"tokenIdToDigitalMediaRelease","outputs":[{"name":"printEdition","type":"uint32"},{"name":"digitalMediaId","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"unpause","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_tokenId","type":"uint256"}],"name":"safeTransferFrom","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"burn","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_creatorAddress","type":"address"}],"name":"removeApprovedTokenCreator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_tokenId","type":"uint256"}],"name":"exists","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_index","type":"uint256"}],"name":"tokenByIndex","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"approvedCreators","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_id","type":"uint256"}],"name":"getDigitalMedia","outputs":[{"name":"id","type":"uint256"},{"name":"totalSupply","type":"uint32"},{"name":"printIndex","type":"uint32"},{"name":"collectionId","type":"uint256"},{"name":"creator","type":"address"},{"name":"metadataPath","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_id","type":"uint256"}],"name":"getCollection","outputs":[{"name":"id","type":"uint256"},{"name":"creator","type":"address"},{"name":"metadataPath","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"paused","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_creator","type":"address"},{"name":"_newCreator","type":"address"}],"name":"changeCreator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_totalSupply","type":"uint32"},{"name":"_collectionId","type":"uint256"},{"name":"_metadataPath","type":"string"}],"name":"createDigitalMedia","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_owner","type":"address"},{"name":"_totalSupply","type":"uint32"},{"name":"_collectionId","type":"uint256"},{"name":"_metadataPath","type":"string"},{"name":"_numReleases","type":"uint32"}],"name":"oboCreateDigitalMediaAndReleases","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_approved","type":"bool"}],"name":"setOboApprovalForAll","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_tokenId","type":"uint256"}],"name":"burnToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_digitalMediaId","type":"uint256"},{"name":"_numReleases","type":"uint32"}],"name":"createDigitalMediaReleases","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"pause","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_id","type":"uint256"}],"name":"getDigitalMediaRelease","outputs":[{"name":"id","type":"uint256"},{"name":"printEdition","type":"uint32"},{"name":"digitalMediaId","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_totalSupply","type":"uint32"},{"name":"_digitalMediaMetadataPath","type":"string"},{"name":"_collectionMetadataPath","type":"string"},{"name":"_numReleases","type":"uint32"}],"name":"createDigitalMediaAndReleasesInNewCollection","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_dmsAddress","type":"address"}],"name":"setV1DigitalMediaStoreAddress","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_oboAddress","type":"address"}],"name":"disableOboAddress","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"v1DigitalMediaStore","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_newCreatorAddress","type":"address"}],"name":"changeSingleCreator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_owner","type":"address"},{"name":"_digitalMediaId","type":"uint256"},{"name":"_numReleases","type":"uint32"}],"name":"oboCreateDigitalMediaReleases","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"creatorRegistryStore","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_tokenId","type":"uint256"},{"name":"_data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_tokenId","type":"uint256"}],"name":"tokenURI","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_tokenId","type":"uint256"}],"name":"resetApproval","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"approvedTokenCreators","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"disabledOboOperators","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_totalSupply","type":"uint32"},{"name":"_collectionId","type":"uint256"},{"name":"_metadataPath","type":"string"},{"name":"_numReleases","type":"uint32"}],"name":"createDigitalMediaAndReleases","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_creatorAddress","type":"address"}],"name":"addApprovedTokenCreator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[{"name":"_tokenName","type":"string"},{"name":"_tokenSymbol","type":"string"},{"name":"_tokenIdStartingCounter","type":"uint256"},{"name":"_dmsAddress","type":"address"},{"name":"_crsAddress","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_owner","type":"address"},{"indexed":false,"name":"_operator","type":"address"},{"indexed":false,"name":"_approved","type":"bool"}],"name":"OboApprovalForAll","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_operator","type":"address"}],"name":"OboDisabledForAll","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"id","type":"uint256"},{"indexed":false,"name":"owner","type":"address"},{"indexed":false,"name":"printEdition","type":"uint32"},{"indexed":false,"name":"tokenURI","type":"string"},{"indexed":false,"name":"digitalMediaId","type":"uint256"}],"name":"DigitalMediaReleaseCreateEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"id","type":"uint256"},{"indexed":false,"name":"storeContractAddress","type":"address"},{"indexed":false,"name":"creator","type":"address"},{"indexed":false,"name":"totalSupply","type":"uint32"},{"indexed":false,"name":"printIndex","type":"uint32"},{"indexed":false,"name":"collectionId","type":"uint256"},{"indexed":false,"name":"metadataPath","type":"string"}],"name":"DigitalMediaCreateEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"id","type":"uint256"},{"indexed":false,"name":"storeContractAddress","type":"address"},{"indexed":false,"name":"creator","type":"address"},{"indexed":false,"name":"metadataPath","type":"string"}],"name":"DigitalMediaCollectionCreateEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"id","type":"uint256"},{"indexed":false,"name":"caller","type":"address"},{"indexed":false,"name":"storeContractAddress","type":"address"}],"name":"DigitalMediaBurnEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"tokenId","type":"uint256"},{"indexed":false,"name":"owner","type":"address"}],"name":"DigitalMediaReleaseBurnEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"digitalMediaId","type":"uint256"},{"indexed":false,"name":"printEdition","type":"uint32"}],"name":"UpdateDigitalMediaPrintIndexEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"creator","type":"address"},{"indexed":false,"name":"newCreator","type":"address"}],"name":"ChangedCreator","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"previousCreatorAddress","type":"address"},{"indexed":true,"name":"newCreatorAddress","type":"address"}],"name":"SingleCreatorChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_tokenId","type":"uint256"}],"name":"Transfer20","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"_approved","type":"address"},{"indexed":false,"name":"_tokenId","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"_operator","type":"address"},{"indexed":false,"name":"_approved","type":"bool"}],"name":"ApprovalForAll","type":"event"},{"anonymous":false,"inputs":[],"name":"Pause","type":"event"},{"anonymous":false,"inputs":[],"name":"Unpause","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"previousOwner","type":"address"},{"indexed":true,"name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"}]`

	man, err := New(rpc, 0, 21000, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	addr := "0x2a46f2ffd99e19a89476e2f62270e0a35bbf0756"
	// ret, err := man.ReadContract(addr, abi, "getCollection", "uint256:40913;", nil)
	// ret, err := man.ReadContract(addr, abi, "getDigitalMediaRelease", "uint256:40913;", nil)
	ret, err := man.ReadContract(addr, abi, "getDigitalMedia", "uint256:40913;", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ret: %x", ret)
}
