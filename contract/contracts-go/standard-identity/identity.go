/*
Copyright (C) BABEC. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
/*
参照CMID合约标准实现：
https://git.chainmaker.org.cn/contracts/standard/-/blob/master/draft/CM-CS-221221-Identity.md
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"chainmaker.org/chainmaker/contract-utils/standard"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"chainmaker.org/chainmaker/contract-utils/address"
	"chainmaker.org/chainmaker/contract-utils/str"
)

const (
	paramAdminAddress = "adminAddress"
	paramAddress      = "address"
	paramIdentities   = "identities"
	paramLevel        = "level"
	paramPkPem        = "pkPem"
	paramMetadata     = "metadata"
	paramStandardName = "standardName"

	methodPreIdentities        = "[identities]"
	methodPreSetIdentity       = "[setIdentity]"
	methodPreSetIdentityBatch  = "[setIdentityBatch]"
	methodPreIdentityOf        = "[identityOf]"
	methodPreLevelOf           = "[levelOf]"
	methodPrePkPemOf           = "[pkPemOf]"
	methodPreAlterAdminAddress = "[alterAdminAddress]"
	methodPreStandards         = "[standards]"
	methodPreSupportsStandard  = "[supportsStandard]"

	keyAdminAddress    = "a"
	keyAddressLevel    = "l"
	keyAddressIdentity = "i"
	keyAddressPkPem    = "p"
)
const (
	// LevelIllegalAccount level 权限编号示例
	LevelIllegalAccount    int = iota // nolint 0 非法
	LevelPlatformAccount              // nolint 1 平台的用户
	LevelPhoneAccount                 // nolint 2 个人手机注册用户
	LevelPersonalAccount              // nolint 3 个人实名用户
	LevelEnterpriseAccount            // nolint 4 企业实名用户
	LevelMax
)

// 标记 IdentityContract 结构体实现 CMID 接口
var _ standard.CMID = (*IdentityContract)(nil)
var _ standard.CMBC = (*IdentityContract)(nil)

// IdentityContract 身份认证合约实现
type IdentityContract struct {
}

// InitContract install contract func
func (i *IdentityContract) InitContract() (result protogo.Response) {
	// 记录异常结果日志
	defer func() {
		if result.Status != 0 {
			sdk.Instance.Warnf(result.Message)
		}
	}()

	sender, err := sdk.Instance.Origin()
	if err != nil {
		return sdk.Error(err.Error())
	}

	err = i.AlterAdminAddress(sender)
	if err != nil {
		return sdk.Error(err.Error())
	}

	return sdk.SuccessResponse
}

// UpgradeContract upgrade contract func
func (i *IdentityContract) UpgradeContract() protogo.Response {
	return sdk.SuccessResponse
}

// InvokeContract the entry func of invoke contract func
func (i *IdentityContract) InvokeContract(method string) (result protogo.Response) {
	// 记录异常结果日志
	defer func() {
		if result.Status != 0 {
			sdk.Instance.Warnf(result.Message)
		}
	}()

	switch method {
	case "Identities":
		return i.identitiesCore()
	case "SetIdentity":
		return i.setIdentityCore()
	case "SetIdentityBatch":
		return i.setIdentityBatchCore()
	case "IdentityOf":
		return i.identityOfCore()
	case "LevelOf":
		return i.levelOfCore()
	case "PkPemOf":
		return i.pkPemOfCore()
	case "AlterAdminAddress":
		return i.alterAdminAddressCore()
	case "Standards":
		return i.standardsCore()
	case "SupportStandard":
		return i.supportStandardCore()
	default:
		return sdk.Error("invalid method")
	}
}

func (i *IdentityContract) identitiesCore() protogo.Response {
	identityMetas := i.Identities()

	data, err := json.Marshal(identityMetas)
	if err != nil {
		return i.error(methodPreIdentities, err.Error())
	}
	return sdk.Success(data)
}

func (i *IdentityContract) setIdentityCore() protogo.Response {
	args := sdk.Instance.GetArgs()

	// 获取参数
	addressParam := string(args[paramAddress])
	levelParam := string(args[paramLevel])
	pkPemParam := string(args[paramPkPem])
	metadataParam := string(args[paramMetadata])

	// 检查参数
	if !address.IsValidAddress(addressParam) {
		return i.error(methodPreSetIdentity, "address format error")
	}
	if str.IsAnyBlank(addressParam, levelParam) {
		return i.error(methodPreSetIdentity, "address or level of param is empty")
	}
	intLevel, err := strconv.Atoi(levelParam)
	if err != nil {
		return i.error(methodPreSetIdentity, err.Error())
	}
	if intLevel < 0 || intLevel >= LevelMax {
		return i.error(methodPreSetIdentity, "level of param is illegal, level="+levelParam)
	}

	// 执行逻辑
	err = i.SetIdentity(addressParam, pkPemParam, intLevel, metadataParam)

	// 返回响应
	if err != nil {
		return i.error(methodPreSetIdentity, err.Error())
	}
	return sdk.SuccessResponse
}

func (i *IdentityContract) setIdentityBatchCore() protogo.Response {
	args := sdk.Instance.GetArgs()

	// 获取参数
	identitiesParam := args[paramIdentities]
	if str.IsAnyBlank(identitiesParam) {
		return i.error(methodPreSetIdentityBatch, "identities of param  is empty")
	}

	// 解析参数
	identities := make([]standard.Identity, 0)
	err := json.Unmarshal(identitiesParam, &identities)
	if err != nil {
		return i.error(methodPreSetIdentityBatch, err.Error())
	}

	// 执行逻辑
	err = i.SetIdentityBatch(identities)

	// 返回响应
	if err != nil {
		return i.error(methodPreSetIdentityBatch, err.Error())
	}
	return sdk.SuccessResponse
}
func (i *IdentityContract) levelOfCore() protogo.Response {
	args := sdk.Instance.GetArgs()

	// 获取、校验参数
	addressParam := string(args[paramAddress])
	if str.IsAnyBlank(addressParam) {
		return i.error(methodPreLevelOf, "address of param is empty")
	}

	// 查询level
	level, err := i.LevelOf(addressParam)
	if err != nil {
		return i.error(methodPreLevelOf, err.Error())
	}

	return sdk.Success([]byte(strconv.Itoa(level)))
}

func (i *IdentityContract) identityOfCore() protogo.Response {
	args := sdk.Instance.GetArgs()

	// 获取、校验参数
	addressStr := string(args[paramAddress])
	if str.IsAnyBlank(addressStr) {
		return i.error(methodPreIdentityOf, "address of param is empty")
	}

	// 查询level
	identity, err := i.IdentityOf(addressStr)
	if err != nil {
		return i.error(methodPreIdentityOf, err.Error())
	}

	identityBytes, err := json.Marshal(identity)
	if err != nil {
		return i.error(methodPreIdentityOf, err.Error())
	}
	return sdk.Success(identityBytes)
}

func (i *IdentityContract) pkPemOfCore() protogo.Response {
	args := sdk.Instance.GetArgs()

	// 获取、校验参数
	addressParam := string(args[paramAddress])
	if str.IsAnyBlank(addressParam) {
		return i.error(methodPrePkPemOf, "address of param is empty")
	}

	// 查询level
	pkPem, err := i.PkPemOf(addressParam)
	if err != nil {
		return i.error(methodPrePkPemOf, err.Error())
	}

	return sdk.Success([]byte(pkPem))
}

func (i *IdentityContract) alterAdminAddressCore() protogo.Response {
	args := sdk.Instance.GetArgs()

	// 获取参数
	adminAddressParam := string(args[paramAdminAddress])

	// 解析参数
	if str.IsAnyBlank(adminAddressParam) {
		return i.error(methodPreAlterAdminAddress, "adminAddress of param is empty")
	}
	if !address.IsValidAddress(adminAddressParam) {
		return i.error(methodPreAlterAdminAddress, "address format error")
	}

	// 保存admin地址
	err := i.AlterAdminAddress(adminAddressParam)

	// 返回响应
	if err != nil {
		return i.error(methodPreAlterAdminAddress, err.Error())
	}
	return sdk.SuccessResponse
}

func (i *IdentityContract) standardsCore() protogo.Response {
	data, err := json.Marshal(i.Standards())
	if err != nil {
		return i.error(methodPreStandards, err.Error())
	}
	return sdk.Success(data)
}

func (i *IdentityContract) supportStandardCore() protogo.Response {
	params := sdk.Instance.GetArgs()

	standardName := string(params[paramStandardName])
	if str.IsAnyBlank(standardName) {
		return i.error(methodPreSupportsStandard, "standardName is empty")
	}

	if i.SupportStandard(standardName) {
		return sdk.Success([]byte(standard.TrueString))
	}
	return sdk.Success([]byte(standard.FalseString))
}

// Identities 获取该合约支持的所有认证类型
func (i *IdentityContract) Identities() (metas []standard.IdentityMeta) {
	metas = make([]standard.IdentityMeta, 5)
	metas[0] = standard.IdentityMeta{Level: LevelIllegalAccount, Description: "未认证"}
	metas[1] = standard.IdentityMeta{Level: LevelPhoneAccount, Description: "个人手机号注册用户"}
	metas[2] = standard.IdentityMeta{Level: LevelPersonalAccount, Description: "个人实名用户"}
	metas[3] = standard.IdentityMeta{Level: LevelPlatformAccount, Description: "企业的用户"}
	metas[4] = standard.IdentityMeta{Level: LevelEnterpriseAccount, Description: "企业实名用户"}
	return metas
}

// SetIdentity 为地址设置认证类型，管理员可调用
func (i *IdentityContract) SetIdentity(address, pkPem string, level int, metadata string) error {
	if !i.senderIsAdmin() {
		sender, _ := sdk.Instance.Sender()
		return errors.New("sender is not admin " + sender)
	}

	identity := standard.Identity{
		Address:  address,
		PkPem:    pkPem,
		Level:    level,
		Metadata: metadata,
	}
	identityBytes, err := json.Marshal(identity)
	if err != nil {
		return err
	}

	err = sdk.Instance.PutStateByte(keyAddressIdentity, address, identityBytes)
	if err != nil {
		return err
	}
	err = sdk.Instance.PutState(keyAddressLevel, address, strconv.Itoa(level))
	if err != nil {
		return err
	}
	if len(pkPem) > 0 {
		err = sdk.Instance.PutState(keyAddressPkPem, address, pkPem)
		if err != nil {
			return err
		}
	}

	// 你可以使用 metadata 在下方执行自身业务逻辑
	i.EmitSetIdentityEvent(address, pkPem, level)

	return nil
}

// EmitSetIdentityEvent 发送设置认证类型事件
func (i *IdentityContract) EmitSetIdentityEvent(address, pkPem string, level int) {
	sdk.Instance.EmitEvent("setIdentity", []string{address, strconv.Itoa(level), pkPem})
}

// SetIdentityBatch 设置多个认证类型，管理员可调用
func (i *IdentityContract) SetIdentityBatch(identities []standard.Identity) error {
	for j := range identities {
		id := identities[j]
		if !address.IsValidAddress(id.Address) {
			return errors.New("address format error")
		}
		if str.IsAnyBlank(id.PkPem) {
			return errors.New("address or level of param is empty")
		}
		if id.Level < 0 || id.Level >= LevelMax {
			return fmt.Errorf("level of param is illegal, level=%d", id.Level)
		}

		err := i.SetIdentity(id.Address, id.PkPem, id.Level, id.Metadata)
		if err != nil {
			return err
		}
	}
	return nil
}

// IdentityOf 获取认证信息
func (i *IdentityContract) IdentityOf(address string) (identity standard.Identity, err error) {
	identityByte, err := sdk.Instance.GetStateByte(keyAddressIdentity, address)
	if err != nil {
		return identity, err
	}
	if str.IsAnyBlank(identityByte) {
		return identity, errors.New("not found")
	}
	err = json.Unmarshal(identityByte, &identity)
	return identity, err
}

// LevelOf 获取认证编号
func (i *IdentityContract) LevelOf(address string) (int, error) {
	level, err := sdk.Instance.GetState(keyAddressLevel, address)
	if err != nil {
		return 0, err
	}
	if str.IsAnyBlank(level) {
		return 0, errors.New("not found")
	}
	return strconv.Atoi(level)
}

// PkPemOf 获取公钥
func (i *IdentityContract) PkPemOf(address string) (string, error) {
	pkPem, err := sdk.Instance.GetState(keyAddressPkPem, address)
	if err != nil {
		return "", err
	}

	if str.IsAnyBlank(pkPem) {
		return "0", errors.New("not found")
	}
	return pkPem, nil
}

// AlterAdminAddress 修改管理员，管理员可调用
func (i *IdentityContract) AlterAdminAddress(adminAddress string) error {
	if !i.senderIsAdmin() {
		sender, _ := sdk.Instance.Sender()
		return errors.New("sender is not admin " + sender)
	}

	adminAddressByte, err := json.Marshal([]string{adminAddress})
	if err != nil {
		return err
	}

	err = sdk.Instance.PutStateFromKeyByte(keyAdminAddress, adminAddressByte)
	if err != nil {
		return errors.New("alter admin address of identityInfo failed")
	}

	sdk.Instance.EmitEvent("alterAdminAddress", []string{adminAddress})

	return nil
}
func (i *IdentityContract) senderIsAdmin() bool {
	sender, _ := sdk.Instance.Sender()
	adminAddressByte, err := sdk.Instance.GetStateFromKeyByte(keyAdminAddress)

	if err != nil {
		sdk.Instance.Warnf("Get totalSupply failed, err:%s", err)
		return false
	}

	if len(adminAddressByte) == 0 {
		return true
	}
	var adminAddress []string
	_ = json.Unmarshal(adminAddressByte, &adminAddress)
	for j := range adminAddress {
		if adminAddress[j] == sender {
			return true
		}
	}
	return false
}

// Standards  获取当前合约支持的标准协议列表
func (i *IdentityContract) Standards() []string {
	return []string{standard.ContractStandardNameCMID, standard.ContractStandardNameCMBC}
}

// SupportStandard  获取当前合约是否支持某合约标准协议
func (i *IdentityContract) SupportStandard(standardName string) bool {
	return standardName == standard.ContractStandardNameCMID || standardName == standard.ContractStandardNameCMBC
}

func (i *IdentityContract) error(methodPre, error string) protogo.Response {
	return sdk.Error(methodPre + "method invoke fail, error: " + error)
}

func main() {
	err := sandbox.Start(new(IdentityContract))
	if err != nil {
		log.Fatal(err)
	}
}
