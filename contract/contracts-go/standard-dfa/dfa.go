/*
 Copyright (C) BABEC. All rights reserved.
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

 SPDX-License-Identifier: Apache-2.0
*/

/*
ContractStandardNameCMDFA ChainMaker - Contract Standard - Digital Fungible Assets
https://git.chainmaker.org.cn/contracts/standard/-/blob/master/draft/CM-CS-221221-DFA.md
*/
package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"chainmaker.org/chainmaker/contract-utils/address"
	"chainmaker.org/chainmaker/contract-utils/standard"

	"chainmaker.org/chainmaker/contract-utils/safemath"
)

const (
	//db key
	balanceKey     = "b"
	allowanceKey   = "a"
	totalSupplyKey = "totalSupplyKey"
	nameKey        = "name"
	symbolKey      = "symbol"
	decimalKey     = "decimal"
	adminKey       = "admin"
)

var (
	defaultName        = "TestToken"
	defaultSymbol      = "TT"
	defaultDecimals    = 18
	defaultTotalSupply = safemath.SafeUintZero
)

var _ standard.CMDFA = (*CmdfaContract)(nil)
var _ standard.CMBC = (*CmdfaContract)(nil)
var _ standard.CMDFAOption = (*CmdfaContract)(nil)

// CmdfaContract erc20 contract
type CmdfaContract struct {
}

// Standards 合约支持的标准
// @return []string
func (c *CmdfaContract) Standards() []string {
	return []string{standard.ContractStandardNameCMBC, standard.ContractStandardNameCMDFA}
}

// SupportStandard 查询本合约是否支持某合约标准
// @param standardName
// @return bool
func (c *CmdfaContract) SupportStandard(standardName string) bool {
	return standardName == standard.ContractStandardNameCMDFA || standardName == standard.ContractStandardNameCMBC
}

// TotalSupply 发行总量
// @return *safemath.SafeUint256
// @return error
func (c *CmdfaContract) TotalSupply() (*safemath.SafeUint256, error) {
	return c.GetTotalSupply()
}

// BalanceOf 账户余额查询
// @param account
// @return *safemath.SafeUint256
// @return error
func (c *CmdfaContract) BalanceOf(account string) (*safemath.SafeUint256, error) {
	return c.GetBalance(account)
}

// Transfer 转账操作
// @param to
// @param amount
// @return error
func (c *CmdfaContract) Transfer(to string, amount *safemath.SafeUint256) error {
	from, err := sdk.Instance.Sender()

	if err != nil {
		return fmt.Errorf("Get sender address failed, err:%s", err)
	}
	err = c.baseTransfer(from, to, amount)
	if err != nil {
		return err
	}
	return nil
}

// TransferFrom 代为转账操作
// @param from
// @param to
// @param amount
// @return error
func (c *CmdfaContract) TransferFrom(from, to string, amount *safemath.SafeUint256) error {
	sender, err := sdk.Instance.Sender()
	if err != nil {
		return fmt.Errorf("Get sender address failed, err:%s", err)
	}
	err = c.baseSpendAllowance(from, sender, amount)
	if err != nil {
		return fmt.Errorf("spend allowance failed, err:%s", err)
	}
	err = c.baseTransfer(from, to, amount)
	if err != nil {
		return err
	}
	return nil
}

// Approve 授权额度
// @param spender
// @param amount
// @return error
func (c *CmdfaContract) Approve(spender string, amount *safemath.SafeUint256) error {
	sender, err := sdk.Instance.Sender()
	if err != nil {
		return fmt.Errorf("Get sender address failed, err:%s", err)
	}
	err = c.baseApprove(sender, spender, amount)
	if err != nil {
		return err
	}
	return nil
}

// Allowance 查询授权额度
// @param owner
// @param spender
// @return *safemath.SafeUint256
// @return error
func (c *CmdfaContract) Allowance(owner, spender string) (*safemath.SafeUint256, error) {
	return c.GetAllowance(owner, spender)
}

// Name Token的名字
// @return string
// @return error
func (c *CmdfaContract) Name() (string, error) {
	return c.GetName()
}

// Symbol Token的符号
// @return string
// @return error
func (c *CmdfaContract) Symbol() (string, error) {
	return c.GetSymbol()
}

// Decimals Token的小数位数
// @return uint8
// @return error
func (c *CmdfaContract) Decimals() (uint8, error) {
	return c.GetDecimals()
}

// Mint 铸造Token
// @param account
// @param amount
// @return error
func (c *CmdfaContract) Mint(account string, amount *safemath.SafeUint256) error {
	//check is admin
	sender, err := sdk.Instance.Sender()
	if err != nil {
		return fmt.Errorf("Get sender address failed, err:%s", err)
	}
	admin, err := c.GetAdmin()
	if err != nil {
		return err
	}
	if sender != admin {
		return errors.New("only admin can mint tokens")
	}
	//call base mint
	err = c.baseMint(account, amount)
	if err != nil {
		return err
	}
	return nil
}

// Burn 销毁Token
// @param amount
// @return error
func (c *CmdfaContract) Burn(amount *safemath.SafeUint256) error {
	spender, err := sdk.Instance.Sender()
	if err != nil {
		return fmt.Errorf("Get sender address failed, err:%s", err)
	}
	//call base burn
	err = c.baseBurn(spender, amount)
	if err != nil {
		return err
	}
	return nil
}

// BurnFrom 代为销毁Token
// @param account
// @param amount
// @return error
func (c *CmdfaContract) BurnFrom(account string, amount *safemath.SafeUint256) error {
	spender, err := sdk.Instance.Sender()
	if err != nil {
		return fmt.Errorf("Get sender address failed, err:%s", err)
	}
	err = c.baseSpendAllowance(account, spender, amount)
	if err != nil {
		return err
	}
	//call base burn
	err = c.baseBurn(account, amount)
	if err != nil {
		return err
	}
	return nil
}

// EmitTransferEvent 触发转账事件
// @param spender
// @param to
// @param amount
func (c *CmdfaContract) EmitTransferEvent(spender, to string, amount *safemath.SafeUint256) {
	sdk.Instance.EmitEvent(standard.TopicTransfer, []string{spender, to, amount.ToString()})
}

// EmitApproveEvent 触发授权事件
// @param owner
// @param spender
// @param amount
func (c *CmdfaContract) EmitApproveEvent(owner, spender string, amount *safemath.SafeUint256) {
	sdk.Instance.EmitEvent(standard.TopicApprove, []string{owner, spender, amount.ToString()})
}

// EmitMintEvent 触发铸造事件
// @param account
// @param amount
func (c *CmdfaContract) EmitMintEvent(account string, amount *safemath.SafeUint256) {
	sdk.Instance.EmitEvent(standard.TopicMint, []string{account, amount.ToString()})
}

// EmitBurnEvent 触发销毁事件
// @param spender
// @param amount
func (c *CmdfaContract) EmitBurnEvent(spender string, amount *safemath.SafeUint256) {
	sdk.Instance.EmitEvent(standard.TopicBurn, []string{spender, amount.ToString()})

}

/////////////////////////Data Access Layer/////////////////////////////////

func createCompositeKey(prefix string, data ...string) string {
	return prefix + "_" + strings.Join(data, "_")
}

// GetUint256 获得DB中的SafeUint256
// @param key
// @param field
// @return *safemath.SafeUint256
// @return error
func (c *CmdfaContract) GetUint256(key, field string) (*safemath.SafeUint256, error) {
	fromBalStr, err := sdk.Instance.GetState(key, field)
	if err != nil {
		return nil, err
	}

	fromBalance, pass := safemath.ParseSafeUint256(string(fromBalStr))
	if !pass {
		return nil, errors.New("invalid uint256 data")
	}
	return fromBalance, nil
}

// GetBalance 获得DB中的账户余额
// @param account
// @return *safemath.SafeUint256
// @return error
func (c *CmdfaContract) GetBalance(account string) (*safemath.SafeUint256, error) {
	return c.GetUint256(balanceKey, account)
}

// SetBalance 设置DB中的账户余额
// @param account
// @param amount
// @return error
func (c *CmdfaContract) SetBalance(account string, amount *safemath.SafeUint256) error {
	return sdk.Instance.PutState(balanceKey, account, amount.ToString())
}

// SetAllowance 设置DB中的授权额度
// @param owner
// @param spender
// @param amount
// @return error
func (c *CmdfaContract) SetAllowance(owner string, spender string, amount *safemath.SafeUint256) error {
	key := createCompositeKey(allowanceKey, owner, spender)
	return sdk.Instance.PutState(key, "", amount.ToString())
}

// GetAllowance 获得DB中的授权额度
// @param owner
// @param spender
// @return *safemath.SafeUint256
// @return error
func (c *CmdfaContract) GetAllowance(owner string, spender string) (*safemath.SafeUint256, error) {
	key := createCompositeKey(allowanceKey, owner, spender)
	return c.GetUint256(key, "")
}

// GetTotalSupply 获得发行总量
// @return *safemath.SafeUint256
// @return error
func (c *CmdfaContract) GetTotalSupply() (*safemath.SafeUint256, error) {
	return c.GetUint256(totalSupplyKey, "")
}

// SetTotalSupply 设置发行总量
// @param amount
// @return error
func (c *CmdfaContract) SetTotalSupply(amount *safemath.SafeUint256) error {
	return sdk.Instance.PutState(totalSupplyKey, "", amount.ToString())
}

// GetName 获得DB中的Name
// @return string
// @return error
func (c *CmdfaContract) GetName() (string, error) {
	return sdk.Instance.GetState(nameKey, "")
}

// SetName 设置DB中的Name
// @param name
// @return error
func (c *CmdfaContract) SetName(name string) error {
	return sdk.Instance.PutState(nameKey, "", name)
}

// GetSymbol 获得DB中的符号
// @return string
// @return error
func (c *CmdfaContract) GetSymbol() (string, error) {
	return sdk.Instance.GetState(symbolKey, "")
}

// SetSymbol 设置DB中的符号
// @param symbol
// @return error
func (c *CmdfaContract) SetSymbol(symbol string) error {
	return sdk.Instance.PutState(symbolKey, "", symbol)
}

// GetDecimals 获得DB中的小数位数
// @return uint8
// @return error
func (c *CmdfaContract) GetDecimals() (uint8, error) {
	d, err := sdk.Instance.GetState(decimalKey, "")
	if err != nil {
		return 0, err
	}
	decimal, err := strconv.ParseUint(d, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(decimal), nil
}

// SetDecimals 设置DB中的小数位数
// @param decimal
// @return error
func (c *CmdfaContract) SetDecimals(decimal uint8) error {
	return sdk.Instance.PutState(decimalKey, "", strconv.Itoa(int(decimal)))
}

// GetAdmin 获得DB中的Admin
// @return string
// @return error
func (c *CmdfaContract) GetAdmin() (string, error) {
	return sdk.Instance.GetState(adminKey, "")

}

// SetAdmin 设置DB中的Admin
// @param admin
// @return error
func (c *CmdfaContract) SetAdmin(admin string) error {
	return sdk.Instance.PutState(adminKey, "", admin)
}

////////////////////////////CMDFA Core/////////////////////////////

// baseTransfer
// @param from
// @param to
// @param amount
// @return error
func (c *CmdfaContract) baseTransfer(from string, to string, amount *safemath.SafeUint256) error {
	//检查from和to的合法性
	if !address.IsValidAddress(from) {
		return errors.New("CMDFA: transfer from the invalid address")
	}
	if address.IsZeroAddress(from) {
		return errors.New("CMDFA: transfer from the zero address")
	}
	if !address.IsValidAddress(to) {
		return errors.New("CMDFA: transfer to the invalid address")
	}
	if address.IsZeroAddress(to) {
		return errors.New("CMDFA: transfer to the zero address")
	}

	//检查from余额充足
	fromBalance, err := c.GetBalance(from)
	if err != nil {
		return err
	}
	if !fromBalance.GTE(amount) {
		return errors.New("CMDFA: transfer amount exceeds balance")
	}
	//更新from和to的余额
	fromNewBalance, _ := safemath.SafeSub(fromBalance, amount)
	err = c.SetBalance(from, fromNewBalance)
	if err != nil {
		return err
	}
	toBalance, err := c.GetBalance(to)
	if err != nil {
		return err
	}
	toNewBalance, ok := safemath.SafeAdd(toBalance, amount)
	if !ok {
		return errors.New("calculate new to balance error")
	}
	err = c.SetBalance(to, toNewBalance)
	if err != nil {
		return err
	}
	//触发事件
	c.EmitTransferEvent(from, to, amount)

	return nil
}

func (c *CmdfaContract) baseApprove(owner string, spender string, amount *safemath.SafeUint256) error {
	//检查from和to的合法性
	if !address.IsValidAddress(owner) {
		return errors.New("CMDFA: approve from the invalid address")
	}
	if address.IsZeroAddress(owner) {
		return errors.New("CMDFA: approve from the zero address")
	}
	if !address.IsValidAddress(spender) {
		return errors.New("CMDFA: approve to the invalid address")
	}
	if address.IsZeroAddress(spender) {
		return errors.New("CMDFA: approve to the zero address")
	}
	//设置Allowance
	err := c.SetAllowance(owner, spender, amount)
	if err != nil {
		return err
	}
	//触发事件Approval
	c.EmitApproveEvent(owner, spender, amount)
	return nil
}

func (c *CmdfaContract) baseSpendAllowance(owner string, spender string, amount *safemath.SafeUint256) error {
	//获得授权的额度
	currentAllowance, err := c.GetAllowance(owner, spender)
	if err != nil {
		return err
	}
	// Does not update the allowance amount in case of infinite allowance.
	if currentAllowance.IsMaxSafeUint256() {
		return nil
	}
	//计算额度是否够用
	if !currentAllowance.GTE(amount) {
		return errors.New("CMDFA: insufficient allowance")
	}
	//扣减授权额度
	newCurrentAllowance, ok := safemath.SafeSub(currentAllowance, amount)
	if !ok {
		return errors.New("spend allowance error")
	}
	return c.baseApprove(owner, spender, newCurrentAllowance)
}
func (c *CmdfaContract) baseMint(account string, amount *safemath.SafeUint256) error {
	//检查account的合法性
	if !address.IsValidAddress(account) {
		return errors.New("CMDFA: mint to the invalid address")
	}
	if address.IsZeroAddress(account) {
		return errors.New("CMDFA: mint to the zero address")
	}

	//更新TotalSupply
	totalSupply, err := c.GetTotalSupply()
	if err != nil {
		return err
	}
	newTotal, ok := safemath.SafeAdd(totalSupply, amount)
	if !ok {
		return errors.New("calculate totalSupply failed")
	}
	err = c.SetTotalSupply(newTotal)
	if err != nil {
		return err
	}
	//更新余额
	toBalance, err := c.GetBalance(account)
	if err != nil {
		return err
	}
	toNewBalance, ok := safemath.SafeAdd(toBalance, amount)
	if !ok {
		return errors.New("calculate new to balance error")
	}
	err = c.SetBalance(account, toNewBalance)
	if err != nil {
		return err
	}
	//触发事件
	c.EmitMintEvent(account, amount)
	return nil
}

func (c *CmdfaContract) baseBurn(account string, amount *safemath.SafeUint256) error {
	//检查account的合法性
	if !address.IsValidAddress(account) {
		return errors.New("CMDFA: burn from the invalid address")
	}
	if address.IsZeroAddress(account) {
		return errors.New("CMDFA: burn from the zero address")
	}
	//检查用户余额充足
	fromBalance, err := c.GetBalance(account)
	if err != nil {
		return err
	}
	if !fromBalance.GTE(amount) {
		return errors.New("CMDFA: burn amount exceeds balance")
	}
	//更新TotalSupply
	totalSupply, err := c.GetTotalSupply()
	if err != nil {
		return err
	}
	newTotal, ok := safemath.SafeSub(totalSupply, amount)
	if !ok {
		return errors.New("calculate totalSupply failed")
	}
	err = c.SetTotalSupply(newTotal)
	if err != nil {
		return err
	}
	//更新余额
	fromNewBalance, ok := safemath.SafeSub(fromBalance, amount)
	if !ok {
		return errors.New("calculate new to balance error")
	}
	err = c.SetBalance(account, fromNewBalance)
	if err != nil {
		return err
	}
	//触发事件
	c.EmitBurnEvent(account, amount)

	return nil
}
