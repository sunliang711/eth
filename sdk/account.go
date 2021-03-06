package sdk

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// GenAccount makes new etherum account
func GenAccount() (sk, pk, address []byte, err error) {
	skStruct, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	sk = crypto.FromECDSA(skStruct)
	//or
	//sk = skStruct.D.Bytes()
	pk = crypto.FromECDSAPub(&skStruct.PublicKey)
	address = crypto.PubkeyToAddress(skStruct.PublicKey).Bytes()
	return
}

// ExportAccount exports encrypted key file as uncrypted output
func ExportAccount(utcFile string, auth string) ([]byte, error) {
	keyJSON, err := ioutil.ReadFile(utcFile)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keyJSON, auth)
	if err != nil {
		return nil, err
	}

	ret, err := key.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type Account struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privatekey"`
	ID         string `json:"id"`
}

func ExportAccountObject(utcFile, auth string) (*Account, error) {
	bs, err := ExportAccount(utcFile, auth)
	if err != nil {
		return nil, err
	}
	var account Account
	err = json.NewDecoder(bytes.NewReader(bs)).Decode(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// HexToAccount convert hex string of private key to ECDSA account
func HexToAccount(hexPrivKey string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address, error) {
	priK, err := crypto.HexToECDSA(hexPrivKey)
	if err != nil {
		return nil, nil, common.Address{}, err
	}
	pubK := priK.Public()
	pubKEcdsa := pubK.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(priK.PublicKey)

	return priK, pubKEcdsa, address, nil
}
