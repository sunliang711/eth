// Package sdk
// @Project:       eth
// @File:          erc721.go.go
// @Author:        eagle
// @Create:        2021/08/16 14:48:38
// @Description:
package sdk

import (
	"fmt"
)

const (
	ERC721_ABI = `[ { "constant": true, "inputs": [ { "name": "_tokenId", "type": "uint256" } ], "name": "getApproved", "outputs": [ { "name": "", "type": "address" } ], "payable": false, "stateMutability": "view", "type": "function" }, { "constant": false, "inputs": [ { "name": "_approved", "type": "address" }, { "name": "_tokenId", "type": "uint256" } ], "name": "approve", "outputs": [], "payable": true, "stateMutability": "payable", "type": "function" }, { "constant": false, "inputs": [ { "name": "_from", "type": "address" }, { "name": "_to", "type": "address" }, { "name": "_tokenId", "type": "uint256" } ], "name": "transferFrom", "outputs": [], "payable": true, "stateMutability": "payable", "type": "function" }, { "constant": false, "inputs": [ { "name": "_from", "type": "address" }, { "name": "_to", "type": "address" }, { "name": "_tokenId", "type": "uint256" } ], "name": "safeTransferFrom", "outputs": [], "payable": true, "stateMutability": "payable", "type": "function" }, { "constant": true, "inputs": [ { "name": "_tokenId", "type": "uint256" } ], "name": "ownerOf", "outputs": [ { "name": "", "type": "address" } ], "payable": false, "stateMutability": "view", "type": "function" }, { "constant": true, "inputs": [ { "name": "_owner", "type": "address" } ], "name": "balanceOf", "outputs": [ { "name": "", "type": "uint256" } ], "payable": false, "stateMutability": "view", "type": "function" }, { "constant": false, "inputs": [ { "name": "_operator", "type": "address" }, { "name": "_approved", "type": "bool" } ], "name": "setApprovalForAll", "outputs": [], "payable": false, "stateMutability": "nonpayable", "type": "function" }, { "constant": false, "inputs": [ { "name": "_from", "type": "address" }, { "name": "_to", "type": "address" }, { "name": "_tokenId", "type": "uint256" }, { "name": "data", "type": "bytes" } ], "name": "safeTransferFrom", "outputs": [], "payable": true, "stateMutability": "payable", "type": "function" }, { "constant": true, "inputs": [ { "name": "_owner", "type": "address" }, { "name": "_operator", "type": "address" } ], "name": "isApprovedForAll", "outputs": [ { "name": "", "type": "bool" } ], "payable": false, "stateMutability": "view", "type": "function" } ]`

	MethodBalanceOf721                = "balanceOf"
	MethodOwnerOf721                  = "ownerOf"
	MethodSafeTransferFrom721         = "safeTransferFrom"
	MethodSafeTransferFromWithData721 = "safeTransferFrom"
	MethodTransferFrom721             = "transferFrom"
	MethodApprove721                  = "approve"
	MethodSetApprovalFroAll721        = "setApprovalForAll"
	MethodGetApproved721              = "getApproved"
	MethodIsApprovedForAll            = "isApprovedForAll"
)

// TransferFrom721 send erc721 transferFrom interface
func (tm *TransactionManager) TransferFrom721(contractAddress string, sk string, from string, to string, tokenId string, price uint64, nonce uint64, limit uint64) (string, error) {
	args := fmt.Sprintf("address:%v;address:%v;uint256:%v", from, to, tokenId)
	return tm.WriteContract(sk, contractAddress, nil, ERC721_ABI, MethodTransferFrom721, args, price, nonce, limit)
}

// TransferFromSync721 send erc721 transferFrom interface
func (tm *TransactionManager) TransferFromSync721(contractAddress string, sk string, from string, to string, tokenId string, price uint64, nonce uint64, limit uint64) (string, uint64, error) {
	args := fmt.Sprintf("address:%v;address:%v;uint256:%v", from, to, tokenId)
	return tm.WriteContractSync(sk, contractAddress, nil, ERC721_ABI, MethodTransferFrom721, args, price, nonce, limit)
}

// SafeTransferFrom721 send erc721 transferFrom interface
func (tm *TransactionManager) SafeTransferFrom721(contractAddress string, sk string, from string, to string, tokenId string, price uint64, nonce uint64, limit uint64) (string, error) {
	args := fmt.Sprintf("address:%v;address:%v;uint256:%v", from, to, tokenId)
	return tm.WriteContract(sk, contractAddress, nil, ERC721_ABI, MethodSafeTransferFrom721, args, price, nonce, limit)
}

// SafeTransferFromSync721 send erc721 transferFrom interface
func (tm *TransactionManager) SafeTransferFromSync721(contractAddress string, sk string, from string, to string, tokenId string, price uint64, nonce uint64, limit uint64) (string, uint64, error) {
	args := fmt.Sprintf("address:%v;address:%v;uint256:%v", from, to, tokenId)
	return tm.WriteContractSync(sk, contractAddress, nil, ERC721_ABI, MethodSafeTransferFrom721, args, price, nonce, limit)
}
