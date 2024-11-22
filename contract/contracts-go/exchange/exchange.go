/*
 Copyright (C) BABEC. All rights reserved.
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

const (
	paramToken  = "tokenId"
	paramFrom   = "from"
	paramTo     = "to"
	paramAmount = "amount"
	trueString  = "true"
)

type exchange interface {

	// 购买 event
	// tokenId uintStr, amount uintStr
	buyNow(tokenId, from, to, amount string) protogo.Response // return "true","false"
	InitContract() protogo.Response                           // return "Init contract success"
	UpgradeContract() protogo.Response                        // return "Upgrade contract success"
}

var _ exchange = (*ExchangeContract)(nil)

// ExchangeContract contract
type ExchangeContract struct {
}

// InitContract install contract func
func (e *ExchangeContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract upgrade contract func
func (e *ExchangeContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

// InvokeContract the entry func of invoke contract func
func (e *ExchangeContract) InvokeContract(method string) protogo.Response {
	args := sdk.Instance.GetArgs()
	if len(method) == 0 {
		return Error("method of param should not be empty")
	}

	switch method {
	case "buyNow":
		amount := string(args[paramAmount])
		token := string(args[paramToken])
		from := string(args[paramFrom])
		to := string(args[paramTo])
		return e.buyNow(token, from, to, amount)
	default:
		return Error("Invalid method" + method)
	}
}
func (e *ExchangeContract) buyNow(tokenId, from, to, amount string) protogo.Response {
	// 查看是否在白名单(注册)
	args := make(map[string][]byte)
	args["address"] = []byte(from)
	resp := sdk.Instance.CallContract("identity", "isApprovedUser", args)
	if string(resp.Payload) != trueString {
		return Error("address[" + from + "] not registered")
	}
	args["address"] = []byte(to)
	resp = sdk.Instance.CallContract("identity", "isApprovedUser", args)
	if string(resp.Payload) != trueString {
		return Error("address[" + to + "] not registered")
	}
	// erc721 转移
	args = make(map[string][]byte)
	args["from"] = []byte(from)
	args["to"] = []byte(to)
	args["tokenId"] = []byte(tokenId)
	resp = sdk.Instance.CallContract("erc721", "safeTransferFrom", args)
	if resp.Status != sdk.OK {
		return Error("erc721 safeTransferFrom error. " + resp.Message)
	}

	// erc20 转账
	args = make(map[string][]byte)
	args["owner"] = []byte(to)
	args["to"] = []byte(from)
	args["amount"] = []byte(amount)
	resp = sdk.Instance.CallContract("erc20", "transferFrom", args)
	if resp.Status != sdk.OK {
		return Error("erc20 transferFrom error. " + resp.Message)
	}

	spender, _ := sdk.Instance.Origin()
	sdk.Instance.EmitEvent("buyNow", []string{from, to, spender, tokenId, amount})

	return sdk.Success([]byte("true"))
}

// Error return error response with message
func Error(message string) protogo.Response {
	return protogo.Response{
		Status:  sdk.ERROR,
		Message: "[exchange] " + message,
	}
}

func main() {
	err := sandbox.Start(new(ExchangeContract))
	if err != nil {
		sdk.Instance.Errorf(err.Error())
	}
}
