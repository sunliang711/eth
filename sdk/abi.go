package sdk

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Pack encodes contract arguments to abi format
/* Usage:
 * args format: <type>:<value>;<type>:<value>;<array type>:<v1>,<v2>...
 * example: uint256:123;bytes:0x12345678;string:"hello world";uint256[]:1,2,3;address:0x1234...;address[]:0x1234,0x5678...;
 * NOTE: for constructor : set methodName to empty string
**/
func Pack(abiStr string, methodName string, args string) ([]byte, error) {
	abiObj, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, fmt.Errorf("abi.JSON error: %v", err)
	}

	var resultArgs []interface{}
	var allArgs []string
	if len(args) > 0 {
		if strings.HasSuffix(args, ";") {
			args = args[:len(args)-1]
		}
		allArgs = strings.Split(args, ";")
	}
	for _, arg := range allArgs {
		if len(arg) == 0 {
			continue
		}
		// arg format: type:value
		splitArg := strings.Split(arg, ":")
		if len(splitArg) != 2 {
			return nil, fmt.Errorf("args format error")
		}
		typ := splitArg[0]
		val := splitArg[1]

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
			resultArgs = append(resultArgs, val)

		case "uint256[]":
			v, err := DecodeUint256StringArray(val)
			if err != nil {
				return nil, err
			}
			resultArgs = append(resultArgs, v)

		case "bytes32[]":
			arr, err := DecodeBytes32StringArray(val)
			if err != nil {
				return nil, err
			}
			resultArgs = append(resultArgs, arr)

		case "bytes[]":
			arr, err := DecodeBytesStringArray(val)
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

		case "address[]":
			hexArray, err := DecodeAddressStringArray(val)
			if err != nil {
				return nil, fmt.Errorf("address[] format error")
			}
			resultArgs = append(resultArgs, hexArray)

		default:
			// TODO
			// uint256[2]...
			return nil, fmt.Errorf("Not support type: %v", typ)
		}
	}
	result, err := abiObj.Pack(methodName, resultArgs...)
	if err != nil {
		return nil, fmt.Errorf("abi.JSON() error: %v", err)
	}
	return result, nil

}

// Unpack decodes output
func Unpack(abiStr string, methodName string, returnData []byte) ([]interface{}, error) {
	abiObj, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, err
	}

	return abiObj.Unpack(methodName, returnData)
}
