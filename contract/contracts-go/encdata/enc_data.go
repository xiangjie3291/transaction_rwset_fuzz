/*
Copyright (C) BABEC. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"strconv"

	bccrypto "chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/hash"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

const (
	enc_data_func_name        = "enc_data"
	enc_get_data_func_name    = "get_enc_data"
	enc_auth_func_name        = "enc_auth"
	enc_get_auth_func_name    = "get_enc_auth"
	enc_update_auth_func_name = "update_enc_auth"
)

const (
	// ENC_KEY parameter of the data encryption contract -- enc_key
	ENC_KEY = "enc_key"

	// ENC_AUTHED_PERSON parameter of the data encryption contract -- authorized_person
	ENC_AUTHED_PERSON = "authorized_person"

	// ENC_AUTHOR parameter of the data encryption contract -- authorizer
	ENC_AUTHOR = "authorizer"

	// ENC_AUTH_SIGN parameter of the data encryption contract -- auth_sign
	ENC_AUTH_SIGN = "auth_sign"

	// ENC_AUTH_LEVEL parameter of the data encryption contract -- auth_level
	ENC_AUTH_LEVEL = "auth_level"

	// DATA_KEY parameter of the data encryption contract -- data_key
	DATA_KEY = "data_key"

	// DATA_VALUE parameter of the data encryption contract -- data_value
	DATA_VALUE = "data_value"

	enc_data_filed = "enc_data_contract_filed"
)

// EncAuth the encdata contract struct
type EncAuth struct {

	// DataKey 加密数据的key
	DataKey []byte `json:"dataKey"`

	// AuthorizedPerson 被授权人 （证书或者公钥）
	AuthorizedPerson []byte `json:"authorizedPerson"`

	// EncKey 加密后的 KEY
	EncKey []byte `json:"encAESKey"`

	// Authorizer 授权人（证书或者公钥）
	Authorizer []byte `json:"authorizer"`

	// AuthSignature 授权人签名
	AuthSignature []byte `json:"authSignature"`

	// AuthLevel 授权等级
	AuthLevel AuthLevel `json:"authLevel"`
}

// AuthMsg information of the authorized person
type AuthMsg struct {
	// AuthorizedPerson the identity of the authorized person
	AuthorizedPerson []byte `json:"authorizedPerson"`

	// AuthLevel the level of the authorizer
	AuthLevel []byte `json:"authLevel"`
}

// AuthLevel the Auth Level
type AuthLevel int32

const (
	// ROOT 合约创建人的等级，能授权ADMIN权限
	ROOT AuthLevel = iota + 1
	// ADMIN 可以向下授权COMMON
	ADMIN
	// COMMON 无法继续授权
	COMMON
)

// EncDataContract the data encryption contract
type EncDataContract struct {
}

func hashHex(data []byte) (string, error) {
	hashBytes, err := hash.GetByStrType("SHA256", data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hashBytes), nil
}

// InitContract install contract func
func (e *EncDataContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract upgrade contract func
func (e *EncDataContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

func (e *EncDataContract) enc_data() protogo.Response {
	args := sdk.Instance.GetArgs()
	var err error
	dataKey, ok := args[DATA_KEY]
	if !ok {
		err = fmt.Errorf("failed to store encrypted data, err: the parameter [%s] does not exist", DATA_KEY)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	dataValue, ok := args[DATA_VALUE]
	if !ok {
		err = fmt.Errorf("failed to store encrypted data, err: the parameter [%s] does not exist", DATA_VALUE)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	encKey, ok := args[ENC_KEY]
	if !ok {
		err = fmt.Errorf("failed to store encrypted data, err: the parameter [%s] does not exist", ENC_KEY)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	authorizedPerson, ok := args[ENC_AUTHED_PERSON]
	if !ok {
		err = fmt.Errorf("failed to store encrypted data, err: the parameter [%s] does not exist", ENC_AUTHED_PERSON)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	_, err = sdk.Instance.GetStateByte(string(dataKey), enc_data_filed)
	if err != nil {
		err = fmt.Errorf("failed to store encrypted data, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	err = sdk.Instance.PutStateByte(string(dataKey), enc_data_filed, dataValue)
	if err != nil {
		err = fmt.Errorf("failed to store encrypted data, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	encAuth := &EncAuth{
		DataKey:          dataKey,
		AuthorizedPerson: authorizedPerson,
		EncKey:           encKey,
		Authorizer:       authorizedPerson,
		AuthLevel:        ROOT,
	}

	encAuthBytes, err := json.Marshal(encAuth)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}
	apHashHex, err := hashHex(authorizedPerson)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 存储数据
	err = sdk.Instance.PutStateByte(apHashHex, string(dataKey), encAuthBytes)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	return sdk.Success([]byte("store encrypted data successfully"))
}

func (e *EncDataContract) get_enc_data() protogo.Response {
	args := sdk.Instance.GetArgs()
	var err error
	dataKey, ok := args[DATA_KEY]
	if !ok {
		err = fmt.Errorf("failed to get encrypted data, err: the parameter [%s] does not exist", DATA_KEY)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 查询结果
	result, err := sdk.Instance.GetStateByte(string(dataKey), enc_data_filed)
	if err != nil {
		err = fmt.Errorf("failed to get encrypted data, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	return sdk.Success(result)
}

// nolint:gocyclo,revive
func (e *EncDataContract) enc_auth() protogo.Response {
	args := sdk.Instance.GetArgs()
	var err error
	authPerson, ok := args[ENC_AUTHED_PERSON]
	if !ok {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTHED_PERSON)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	authorizer, ok := args[ENC_AUTHOR]
	if !ok {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTHOR)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	authSign, ok := args[ENC_AUTH_SIGN]
	if !ok {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTH_SIGN)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	authLevel, ok := args[ENC_AUTH_LEVEL]
	if !ok {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTH_LEVEL)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	encKey, ok := args[ENC_KEY]
	if !ok {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_KEY)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	dataKey, ok := args[DATA_KEY]
	if !ok {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", DATA_KEY)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	if bytes.Equal(authPerson, authorizer) {
		err = fmt.Errorf("fail to authorize, err: can not authorize for self")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 验证授权者是否在链上
	authorizerHashHex, err := hashHex(authorizer)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 查询结果
	result, err := sdk.Instance.GetStateByte(authorizerHashHex, string(dataKey))
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	if len(result) == 0 {
		err = fmt.Errorf("fail to authorize, err: the authorizer does not exist")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 验证被授权者是否已经存在
	authPersonHashHex, err := hashHex(authPerson)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 查询结果
	result2, err := sdk.Instance.GetStateByte(authPersonHashHex, string(dataKey))
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	if len(result2) != 0 {
		err = fmt.Errorf("fail to authorize, err: the authorized person already exist")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	var authorizerInfo EncAuth
	err = json.Unmarshal(result, &authorizerInfo)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 构建MSG，验签
	msg := &AuthMsg{
		AuthorizedPerson: authPerson,
		AuthLevel:        authLevel,
	}

	var authorizerCert *bcx509.Certificate
	certBlock, rest := pem.Decode(authorizer)
	if certBlock == nil {
		authorizerCert, err = bcx509.ParseCertificate(rest)
		if err != nil {
			err = fmt.Errorf("fail to authorize, err: %s", err.Error())
			sdk.Instance.Log(err.Error())
			return sdk.Error(err.Error())
		}
	} else {
		authorizerCert, err = bcx509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			err = fmt.Errorf("fail to authorize, err: %s", err.Error())
			sdk.Instance.Log(err.Error())
			return sdk.Error(err.Error())
		}
	}

	// msg 序列化
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 验签
	ok, err = authorizerCert.PublicKey.VerifyWithOpts(msgBytes, authSign, &bccrypto.SignOpts{
		Hash: bccrypto.HASH_TYPE_SM3,
		UID:  bccrypto.CRYPTO_DEFAULT_UID,
	})

	if !ok || err != nil {
		err = fmt.Errorf("fail to authorize, err: invalid signature")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	intAuthLevel, err := strconv.Atoi(string(authLevel))
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	switch AuthLevel(intAuthLevel) {
	case ADMIN:
		if authorizerInfo.AuthLevel != ROOT {
			err = fmt.Errorf("fail to authorize, err: the authorizer has no power to authorize")
			sdk.Instance.Log(err.Error())
			return sdk.Error(err.Error())
		}
	case COMMON:
		if authorizerInfo.AuthLevel != ROOT && authorizerInfo.AuthLevel != ADMIN {
			err = fmt.Errorf("fail to authorize, err: the authorizer has no power to authorize")
			sdk.Instance.Log(err.Error())
			return sdk.Error(err.Error())
		}
	default:
		err = fmt.Errorf("fail to authorize, err: invalid auth level")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 授权存储上链
	encAuth := &EncAuth{
		AuthorizedPerson: authPerson,
		EncKey:           encKey,
		Authorizer:       authorizer,
		AuthLevel:        AuthLevel(intAuthLevel),
		DataKey:          dataKey,
		AuthSignature:    authSign,
	}

	encAuthBytes, err := json.Marshal(encAuth)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}
	apHashHex, err := hashHex(authPerson)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 存储数据
	err = sdk.Instance.PutStateByte(apHashHex, string(dataKey), encAuthBytes)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	return sdk.Success([]byte("store auth info successfully"))

}

func (e *EncDataContract) get_enc_auth() protogo.Response {
	args := sdk.Instance.GetArgs()
	var err error
	authorizer, ok := args[ENC_AUTHOR]
	if !ok {
		err = fmt.Errorf("fail to get auth info, err: the parameter [%s] does not exist", ENC_AUTHOR)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	dataKey, ok := args[DATA_KEY]
	if !ok {
		err = fmt.Errorf("fail to get auth info, err: the parameter [%s] does not exist", DATA_KEY)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	authorizerHashHex, err := hashHex(authorizer)
	if err != nil {
		err = fmt.Errorf("fail to get auth info, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 查询结果
	result, err := sdk.Instance.GetStateByte(authorizerHashHex, string(dataKey))
	if err != nil {
		err = fmt.Errorf("fail to get auth info, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	if result == nil {
		err = fmt.Errorf("fail to get auth info, err: the authorized information does not exist")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	var authInfo EncAuth

	err = json.Unmarshal(result, &authInfo)
	if err != nil {
		err = fmt.Errorf("fail to get auth info, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	if len(authInfo.EncKey) == 0 {
		return sdk.Error("the encrypted key is empty")
	}

	return sdk.Success(authInfo.EncKey)
}

// nolint:gocyclo,revive
func (e *EncDataContract) update_enc_auth() protogo.Response {
	args := sdk.Instance.GetArgs()
	var err error
	authPerson, ok := args[ENC_AUTHED_PERSON]
	if !ok {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTHED_PERSON)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	authorizer, ok := args[ENC_AUTHOR]
	if !ok {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTHOR)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	authSign, ok := args[ENC_AUTH_SIGN]
	if !ok {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTH_SIGN)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	authLevel, ok := args[ENC_AUTH_LEVEL]
	if !ok {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTH_LEVEL)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	dataKey, ok := args[DATA_KEY]
	if !ok {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", DATA_KEY)
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	if bytes.Equal(authPerson, authorizer) {
		err = fmt.Errorf("fail to authorize, err: can not authorize for self")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 验证授权者是否在链上
	authorizerHashHex, err := hashHex(authorizer)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 查询结果
	result, err := sdk.Instance.GetStateByte(authorizerHashHex, string(dataKey))
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	if len(result) == 0 {
		err = fmt.Errorf("fail to authorize, err: the authorizer does not exist")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 验证被授权者是否已经存在
	authPersonHashHex, err := hashHex(authPerson)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 查询结果
	result2, err := sdk.Instance.GetStateByte(authPersonHashHex, string(dataKey))
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	if len(result2) == 0 {
		err = fmt.Errorf("fail to authorize, err: the authorized person does not exist")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	var authorizerInfo EncAuth
	err = json.Unmarshal(result, &authorizerInfo)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	var authPersonInfo EncAuth
	err = json.Unmarshal(result2, &authPersonInfo)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 构建MSG，验签
	msg := &AuthMsg{
		AuthorizedPerson: authPerson,
		AuthLevel:        authLevel,
	}

	var authorizerCert *bcx509.Certificate
	certBlock, rest := pem.Decode(authorizer)
	if certBlock == nil {
		authorizerCert, err = bcx509.ParseCertificate(rest)
		if err != nil {
			err = fmt.Errorf("fail to authorize, err: %s", err.Error())
			sdk.Instance.Log(err.Error())
			return sdk.Error(err.Error())
		}
	} else {
		authorizerCert, err = bcx509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			err = fmt.Errorf("fail to authorize, err: %s", err.Error())
			sdk.Instance.Log(err.Error())
			return sdk.Error(err.Error())
		}
	}

	// msg 序列化
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 验签
	ok, err = authorizerCert.PublicKey.VerifyWithOpts(msgBytes, authSign, &bccrypto.SignOpts{
		Hash: bccrypto.HASH_TYPE_SM3,
		UID:  bccrypto.CRYPTO_DEFAULT_UID,
	})

	if !ok || err != nil {
		err = fmt.Errorf("fail to authorize, err: invalid signature")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	intAuthLevel, err := strconv.Atoi(string(authLevel))
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	if authorizerInfo.AuthLevel != ROOT {
		err = fmt.Errorf("fail to authorize, only the root can change the user permissions")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	if authPersonInfo.AuthLevel == ROOT {
		err = fmt.Errorf("fail to authorize, can't update the user permissions who is the root")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	if AuthLevel(intAuthLevel) != ADMIN && AuthLevel(intAuthLevel) != COMMON {
		err = fmt.Errorf("fail to authorize, err: invalid auth level")
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 授权存储上链
	encAuth := &EncAuth{
		AuthorizedPerson: authPerson,
		EncKey:           authPersonInfo.EncKey,
		Authorizer:       authorizer,
		AuthLevel:        AuthLevel(intAuthLevel),
		DataKey:          dataKey,
		AuthSignature:    authSign,
	}

	encAuthBytes, err := json.Marshal(encAuth)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}
	apHashHex, err := hashHex(authPerson)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	// 存储数据
	err = sdk.Instance.PutStateByte(apHashHex, string(dataKey), encAuthBytes)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		sdk.Instance.Log(err.Error())
		return sdk.Error(err.Error())
	}

	return sdk.Success([]byte("store auth info successfully"))
}

// InvokeContract the entry func of invoke contract func
func (e *EncDataContract) InvokeContract(method string) protogo.Response {

	switch method {
	case enc_data_func_name:
		return e.enc_data()
	case enc_auth_func_name:
		return e.enc_auth()
	case enc_get_data_func_name:
		return e.get_enc_data()
	case enc_get_auth_func_name:
		return e.get_enc_auth()
	case enc_update_auth_func_name:
		return e.update_enc_auth()
	default:
		return sdk.Error("invalid method")
	}

}

func main() {
	err := sandbox.Start(new(EncDataContract))
	if err != nil {
		log.Fatal(err)
	}
}
