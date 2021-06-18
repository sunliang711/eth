// Package sdk
// @Project:       eth
// @File:          codec.go
// @Author:        eagle
// @Create:        2021/06/18 13:46:06
// @Description:
package sdk

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// DecodeHexString decodes hex string with "0x" or without "0x"
func DecodeHexString(input string) ([]byte, error) {
	if strings.HasPrefix(input, "0x") {
		input = input[2:]
	}
	if len(input)&1 != 0 {
		return nil, fmt.Errorf("input length is not even.")
	}
	hexVal := make([]byte, len(input)/2)
	n, err := hex.Decode(hexVal, []byte(input))
	if err != nil {
		return nil, err
	}
	if n != len(hexVal) {
		return nil, fmt.Errorf("input error")
	}

	return hexVal, nil
}

func DecodeAddressStringArray(val string) ([]common.Address, error) {
	var ret []common.Address
	if strings.HasSuffix(val, ",") {
		val = val[:len(val)-1]
	}
	allHexStrings := strings.Split(val, ",")
	for _, s := range allHexStrings {
		v, err := DecodeHexString(s)
		if err != nil {
			return nil, err
		}
		ret = append(ret, common.BytesToAddress(v))
	}
	return ret, nil
}

func DecodeBytes32String(val string) ([32]byte, error) {

	if strings.HasPrefix(val, "0x") {
		val = val[2:]
	}
	if len(val) > 64 {
		return [32]byte{}, fmt.Errorf("bytes32 value greater than 32 bytes")
	}
	hexVal, err := DecodeHexString(val)
	hexVal = common.LeftPadBytes(hexVal, 32)
	if err != nil {
		return [32]byte{}, fmt.Errorf("bytes32 value format error")
	}
	var v [32]byte
	copy(v[:], hexVal[:32])
	return v, nil
}

// DecodeBytes32StringArray val format: 0x1234...,0x4567...,0x9999...
func DecodeBytes32StringArray(val string) ([][32]byte, error) {
	var ret [][32]byte
	if strings.HasSuffix(val, ",") {
		val = val[:len(val)-1]
	}
	allBytes32 := strings.Split(val, ",")
	for _, b := range allBytes32 {
		v, err := DecodeBytes32String(b)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}

	return ret, nil
}

func DecodeUint256String(val string) (*big.Int, error) {
	v, err := strconv.Atoi(val)
	if err != nil {
		return nil, fmt.Errorf("uint256 value format error")
	}
	return big.NewInt(int64(v)), nil
}

func DecodeUint256StringArray(val string) ([]*big.Int, error) {
	var ret []*big.Int
	if strings.HasSuffix(val, ",") {
		val = val[:len(val)-1]
	}
	allUint256 := strings.Split(val, ",")
	for _, u := range allUint256 {
		v, err := DecodeUint256String(u)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func DecodeBytesStringArray(val string) ([][]byte, error) {
	var ret [][]byte
	if strings.HasSuffix(val, ",") {
		val = val[:len(val)-1]
	}
	bytesArr := strings.Split(val, ",")
	for _, b := range bytesArr {
		v, err := DecodeHexString(b)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}
