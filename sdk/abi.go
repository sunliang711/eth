package sdk

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

// Pack encodes contract arguments to abi format
/* Usage:
* args format: uint256:123;bytes:0x12345678;string:"hello world";uint256[2]:1,2;uint256[]:1,2,3;address:0x1234;address[]:0x1234,0x5678...;
 * NOTE: for constructor : set methodName to empty string
**/
func Pack(abiStr string, methodName string, args string) ([]byte, error) {
	abiObj, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, fmt.Errorf("abi.JSON error: %v", err)
	}

	resultArgs := []interface{}{}
	allArgs := []string{}
	if len(args) > 0 {
		if strings.HasSuffix(args, ";") {
			args = args[:len(args)-1]
		}
		allArgs = strings.Split(args, ";")
	}
	logrus.Debugf("allArgs: '%v'", allArgs)
	logrus.Debugf("len(allArgs): %v", len(allArgs))
	for _, arg := range allArgs {
		if len(arg) == 0 {
			continue
		}
		//arg format: type:value
		logrus.Debugf("arg: %s", arg)
		splitedArg := strings.Split(arg, ":")
		if len(splitedArg) != 2 {
			return nil, fmt.Errorf("args format error")
		}
		typ := splitedArg[0]
		val := splitedArg[1]

		switch typ {
		case "uint256":
			v, err := DecodeUint256String(val)
			if err != nil {
				return nil, err
			}

			resultArgs = append(resultArgs, v)
		case "bytes":
			hexVal, err := DecodeHexString(val)
			if err != nil {
				return nil, fmt.Errorf("bytes value format error")
			}
			resultArgs = append(resultArgs, []byte(hexVal))
		case "bytes32":
			v, err := DecodeBytes32String(val)
			if err != nil {
				return nil, err
			}
			resultArgs = append(resultArgs, v)
		case "string":
			val = strings.Trim(val, `"`)
			logrus.Debugf("string value: %v", val)
			resultArgs = append(resultArgs, val)
		case "uint256[]":
			v, err := DecodeUint256ArrayString(val)
			if err != nil {
				return nil, err
			}
			resultArgs = append(resultArgs, v)
		case "bytes32[]":
			arr, err := DecodeBytes32ArrayString(val)
			if err != nil {
				return nil, err
			}
			resultArgs = append(resultArgs, arr)

		case "bytes[]":
			arr, err := DecodeBytesArrayString(val)
			if err != nil {
				return nil, err
			}
			resultArgs = append(resultArgs, arr)
		case "address":
			hexVal, err := DecodeHexString(val)
			if err != nil {
				return nil, fmt.Errorf("address format error")
			}
			resultArgs = append(resultArgs, common.BytesToAddress(hexVal))
		//TODO not tested
		case "address[]":
			hexArray, err := DecodeAddressStringArray(val)
			if err != nil {
				return nil, fmt.Errorf("address[] format error")
			}
			resultArgs = append(resultArgs, hexArray)

		default:
			//TODO
			//uint256[2]...
			return nil, fmt.Errorf("Not support type: %v", typ)
		}
	}
	logrus.Debugf("resultArgs: %x", resultArgs)
	result, err := abiObj.Pack(methodName, resultArgs...)
	if err != nil {
		return nil, fmt.Errorf("abi.JSON() error: %v", err)
	}
	return result, nil

}

// Unpack decodes output
func Unpack(abiStr string, methodName string, returnData string) error {
	abiObj, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		logrus.Fatalf("abi.JSON error: %v", err)
	}

	var result []interface{}
	var resultType []string
	method := abiObj.Methods[methodName]
	for _, o := range method.Outputs {
		resultType = append(resultType, o.Type.String())
		switch o.Type.String() {
		case "bytes":
			result = append(result, &[]byte{})
		case "bytes32":
			result = append(result, &[32]byte{})
		case "uint256":
			result = append(result, big.NewInt(0))
		case "address":
			result = append(result, &common.Address{})
		case "string":
			var s string
			result = append(result, &s)
		default:
			return fmt.Errorf("Not support type: %s\n", o.Type.String())
		}
	}

	data, err := DecodeHexString(returnData)
	if err != nil {
		return err
	}

	abiObj.Unpack(&result, methodName, data)
	for i, r := range result {
		fmt.Printf("type: %s, value: %+v\n", resultType[i], r)
	}
	//TODO return value

	return nil
}

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

// val format: 0x1234...,0x4567...,0x9999...
func DecodeBytes32ArrayString(val string) ([][32]byte, error) {
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

func DecodeUint256ArrayString(val string) ([]*big.Int, error) {
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

func DecodeBytesArrayString(val string) ([][]byte, error) {
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
