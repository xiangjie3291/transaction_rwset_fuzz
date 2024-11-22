/*
Copyright (C) BABEC. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
/*
参照CMEVI合约标准实现：
https://git.chainmaker.org.cn/contracts/standard/-/blob/master/draft/CM-CS-221221-Evidence.md
*/

package main

import (
	"encoding/json"
	"errors"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"chainmaker.org/chainmaker/contract-utils/standard"
	"chainmaker.org/chainmaker/contract-utils/str"
)

const (
	paramId           = "id"
	paramHash         = "hash"
	paramMetadata     = "metadata"
	paramEvidences    = "evidences"
	paramStandardName = "standardName"

	methodPreEvidence         = "[evidence]"
	methodPreExistsOfHash     = "[existsOfHash]"
	methodPreExistsOfId       = "[existsOfId]"
	methodPreFindByHash       = "[findByHash]"
	methodPreFindById         = "[findById]"
	methodPreEvidenceBatch    = "[evidenceBatch]"
	methodPreStandards        = "[standards]"
	methodPreSupportsStandard = "[supportsStandard]"

	keyEvidenceHash = "h" // 此为存入数据库的世界状态key，故越短越好
	keyEvidenceId   = "i"
)

// 标记 EvidenceContract 结构体实现 CMEVI 接口
var _ standard.CMEVI = (*EvidenceContract)(nil)
var _ standard.CMBC = (*EvidenceContract)(nil)

// EvidenceContract 存证合约实现
type EvidenceContract struct {
}

// InitContract install contract func
func (e *EvidenceContract) InitContract() protogo.Response {
	return sdk.SuccessResponse
}

// UpgradeContract upgrade contract func
func (e *EvidenceContract) UpgradeContract() protogo.Response {
	return sdk.SuccessResponse
}

// InvokeContract the entry func of invoke contract func
func (e *EvidenceContract) InvokeContract(method string) (result protogo.Response) {
	// 记录异常结果日志
	defer func() {
		if result.Status != 0 {
			sdk.Instance.Warnf(result.Message)
		}
	}()

	switch method {
	case "Evidence":
		return e.evidenceCore()
	case "EvidenceBatch":
		return e.evidenceBatchCore()
	case "ExistsOfHash":
		return e.existsOfHashCore()
	case "ExistsOfId":
		return e.existsOfIdCore()
	case "FindByHash":
		return e.findByHashCore()
	case "FindById":
		return e.findByIdCore()
	case "Standards":
		return e.standardsCore()
	case "SupportStandard":
		return e.supportStandardCore()
	default:
		return sdk.Error("invalid method")
	}
}

func (e *EvidenceContract) evidenceCore() protogo.Response {
	params := sdk.Instance.GetArgs()

	// 获取参数
	id := string(params[paramId])
	hash := string(params[paramHash])
	metadata := string(params[paramMetadata])

	// 执行逻辑
	err := e.Evidence(id, hash, metadata)
	if err != nil {
		return e.error(methodPreEvidence, err.Error())
	}

	// 返回OK
	return sdk.SuccessResponse
}

func (e *EvidenceContract) evidenceBatchCore() protogo.Response {
	params := sdk.Instance.GetArgs()

	// 获取、解析参数
	evidencesParam := params[paramEvidences]
	if str.IsAnyBlank(evidencesParam) {
		return e.error(methodPreEvidenceBatch, "evidences is empty")
	}
	evidences := make([]standard.Evidence, 0)
	err := json.Unmarshal(evidencesParam, &evidences)
	if err != nil {
		return e.error(methodPreEvidenceBatch, err.Error())
	}

	// 存证
	err = e.EvidenceBatch(evidences)
	if err != nil {
		return e.error(methodPreEvidenceBatch, err.Error())
	}
	return sdk.SuccessResponse
}

func (e *EvidenceContract) existsOfHashCore() protogo.Response {
	params := sdk.Instance.GetArgs()

	hash := string(params[paramHash])
	if str.IsAnyBlank(hash) {
		return e.error(methodPreExistsOfHash, "hash is empty")
	}

	if exists, err := e.ExistsOfHash(hash); err != nil {
		return e.error(methodPreExistsOfHash, err.Error())
	} else if exists {
		return sdk.Success([]byte(standard.TrueString))
	} else {
		return sdk.Success([]byte(standard.FalseString))
	}
}

func (e *EvidenceContract) existsOfIdCore() protogo.Response {
	params := sdk.Instance.GetArgs()

	id := string(params[paramId])
	if str.IsAnyBlank(id) {
		return e.error(methodPreExistsOfId, "id is empty")
	}

	if exists, err := e.ExistsOfId(id); err != nil {
		return e.error(methodPreExistsOfId, err.Error())
	} else if exists {
		return sdk.Success([]byte(standard.TrueString))
	} else {
		return sdk.Success([]byte(standard.FalseString))
	}
}

func (e *EvidenceContract) findByHashCore() protogo.Response {
	params := sdk.Instance.GetArgs()

	hash := string(params[paramHash])
	if str.IsAnyBlank(hash) {
		return e.error(methodPreFindByHash, "hash is empty")
	}

	evidence, err := e.FindByHash(hash)
	if err != nil {
		return e.error(methodPreFindByHash, err.Error())
	}

	data, err := json.Marshal(evidence)
	if err != nil {
		return e.error(methodPreFindByHash, err.Error())
	}
	return sdk.Success(data)
}

func (e *EvidenceContract) findByIdCore() protogo.Response {
	params := sdk.Instance.GetArgs()

	id := string(params[paramId])
	if str.IsAnyBlank(id) {
		return e.error(methodPreFindById, "id is empty")
	}

	evidence, err := e.FindById(id)
	if err != nil {
		return e.error(methodPreFindById, err.Error())
	}

	data, err := json.Marshal(evidence)
	if err != nil {
		return e.error(methodPreFindById, err.Error())
	}
	return sdk.Success(data)
}

func (e *EvidenceContract) standardsCore() protogo.Response {
	data, err := json.Marshal(e.Standards())
	if err != nil {
		return e.error(methodPreStandards, err.Error())
	}
	return sdk.Success(data)
}

func (e *EvidenceContract) supportStandardCore() protogo.Response {
	params := sdk.Instance.GetArgs()

	standardName := string(params[paramStandardName])
	if str.IsAnyBlank(standardName) {
		return e.error(methodPreSupportsStandard, "standardName is empty")
	}

	if e.SupportStandard(standardName) {
		return sdk.Success([]byte(standard.TrueString))
	}
	return sdk.Success([]byte(standard.FalseString))
}

// Evidence 存证
func (e *EvidenceContract) Evidence(id string, hash string, metadata string) error {
	// 校验是否存在
	existsHash, err1 := sdk.Instance.GetStateByte(keyEvidenceHash, hash)
	if err1 != nil {
		return err1
	}
	if len(existsHash) > 0 {
		return errors.New("hash already exists: " + hash)
	}

	existsId, err2 := sdk.Instance.GetStateByte(keyEvidenceId, id)
	if err2 != nil {
		return err2
	}
	if len(existsId) > 0 {
		return errors.New("id already exists: " + id)
	}

	// 获取区块信息, 构建存证对象
	txId, _ := sdk.Instance.GetTxId()
	blockHeight, _ := sdk.Instance.GetBlockHeight()
	txTimeStamp, _ := sdk.Instance.GetTxTimeStamp()

	evidence := standard.Evidence{
		Id:          id,
		TxId:        txId,
		Hash:        hash,
		BlockHeight: blockHeight,
		Timestamp:   txTimeStamp,
		Metadata:    metadata,
	}

	data, err3 := json.Marshal(&evidence)
	if err3 != nil {
		return err3
	}

	err := sdk.Instance.PutStateByte(keyEvidenceHash, hash, data)
	if err != nil {
		return err
	}
	err = sdk.Instance.PutStateByte(keyEvidenceId, id, []byte(hash))
	if err != nil {
		return err
	}

	return nil
}

// ExistsOfHash 哈希是否存在
func (e *EvidenceContract) ExistsOfHash(hash string) (bool, error) {
	existsHash, err := sdk.Instance.GetStateByte(keyEvidenceHash, hash)
	return len(existsHash) > 0, err
}

// ExistsOfId ID是否存在
func (e *EvidenceContract) ExistsOfId(id string) (bool, error) {
	existsId, err := sdk.Instance.GetStateByte(keyEvidenceId, id)
	return len(existsId) > 0, err
}

// FindByHash 根据哈希查找
func (e *EvidenceContract) FindByHash(hash string) (*standard.Evidence, error) {
	data, err := sdk.Instance.GetStateByte(keyEvidenceHash, hash)
	if err != nil {
		return nil, err
	}
	evi := &standard.Evidence{}
	err = json.Unmarshal(data, evi)
	return evi, err
}

// FindById 根据id查找
func (e *EvidenceContract) FindById(id string) (*standard.Evidence, error) {
	hash, err := sdk.Instance.GetState(keyEvidenceId, id)
	if err != nil {
		return nil, err
	}
	evidenceByte, err := sdk.Instance.GetStateByte(keyEvidenceHash, hash)
	if err != nil {
		return nil, err
	}

	evi := &standard.Evidence{}
	err = json.Unmarshal(evidenceByte, evi)
	return evi, err
}

// EvidenceBatch 批量存证
func (e *EvidenceContract) EvidenceBatch(evidences []standard.Evidence) error {
	for i := range evidences {
		err := e.Evidence(evidences[i].Id, evidences[i].Hash, evidences[i].Metadata)
		if err != nil {
			return err
		}
	}
	return nil
}

// Standards  获取当前合约支持的标准协议列表
func (e *EvidenceContract) Standards() []string {
	return []string{standard.ContractStandardNameCMEVI, standard.ContractStandardNameCMBC}
}

// SupportStandard  获取当前合约是否支持某合约标准协议
func (e *EvidenceContract) SupportStandard(standardName string) bool {
	return standardName == standard.ContractStandardNameCMEVI || standardName == standard.ContractStandardNameCMBC
}

func (e *EvidenceContract) error(methodPre, error string) protogo.Response {
	return sdk.Error(methodPre + "method invoke fail, error: " + error)
}

func main() {
	err := sandbox.Start(new(EvidenceContract))
	if err != nil {
		sdk.Instance.Errorf(err.Error())
	}
}
