/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"fmt"
	"strings"

	"chainmaker.org/chainmaker/common/v2/cert"
	"chainmaker.org/chainmaker/common/v2/crypto"
	"github.com/pkg/errors"

	"chainmaker.org/chainmaker/common/v2/crypto/pkcs11"
	"chainmaker.org/chainmaker/common/v2/crypto/sdf"
	"github.com/spf13/viper"
)

const (
	PKCS11 = "pkcs11"
	SDF    = "sdf"
)

type PKCS11KeysConfig struct {
	P11KeysMap map[string]OrgKeys `mapstructure:"pkcs11_keys"`
}

type OrgKeys struct {
	CA       []string `mapstructure:"ca"`
	UserKeys UserKeys `mapstructure:"user"`
	NodeKeys NodeKeys `mapstructure:"node"`
}

type NodeKeys struct {
	Consensus []string `mapstructure:"consensus"`
	Common    []string `mapstructure:"common"`
}
type UserKeys struct {
	Admin  []string `mapstructure:"admin"`
	Client []string `mapstructure:"client"`
	Light  []string `mapstructure:"light"`
}

var p11KeysConfig *PKCS11KeysConfig

const (
	defaultPKCS11KeysPath = "../config/pkcs11_keys.yml"
)

func LoadPKCS11KeysConfig(path string) error {
	p11KeysConfig = &PKCS11KeysConfig{}

	if path == "" {
		path = defaultPKCS11KeysPath
	}

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&p11KeysConfig); err != nil {
		return err
	}

	return initCertPKCS11()
}

func initCertPKCS11() error {
	if cryptoGenConfig == nil || len(cryptoGenConfig.Item) <= 0 {
		return errors.New("cryptoGenConfig not initialized")
	}
	p11Config := cryptoGenConfig.Item[0].P11Config
	var handle interface{}
	var err error
	if p11Config.Enabled {
		if p11Config.Type == PKCS11 {
			handle, err = pkcs11.New(p11Config.Library, p11Config.Label, p11Config.Password, p11Config.SessionCacheSize, p11Config.Hash)
			if err != nil {
				return errors.WithMessage(err, "failed to initialize pkcs11 handle")
			}
		} else if p11Config.Type == SDF {
			handle, err = sdf.New(p11Config.Library, p11Config.SessionCacheSize)
			if err != nil {
				return errors.WithMessage(err, "failed to initialize pkcs11 handle")
			}
		}
	}
	cert.InitP11Handle(handle)
	return nil
}

func SetPrivKeyContext(keyType crypto.KeyType, orgName string, j int, usage string) error {
	if !PKCS11Enabled() {
		return nil
	}
	if _, exist := p11KeysConfig.P11KeysMap[orgName]; !exist {
		return fmt.Errorf("pkcs11 org keys not set, orgName = %s", orgName)
	}
	var keyLabel string
	switch usage {
	case "ca":
		if j >= len(p11KeysConfig.P11KeysMap[orgName].CA) {
			return fmt.Errorf("pkcs11 key not set, orgName = %s, caId = %d", orgName, j)
		}
		keyLabel = p11KeysConfig.P11KeysMap[orgName].CA[j]
	case "admin":
		if j >= len(p11KeysConfig.P11KeysMap[orgName].UserKeys.Admin) {
			return fmt.Errorf("pkcs11 key not set, orgName = %s, adminId = %d", orgName, j)
		}
		keyLabel = p11KeysConfig.P11KeysMap[orgName].UserKeys.Admin[j]
	case "light":
		if j >= len(p11KeysConfig.P11KeysMap[orgName].UserKeys.Light) {
			return fmt.Errorf("pkcs11 key not set, orgName = %s, lightId = %d", orgName, j)
		}
		keyLabel = p11KeysConfig.P11KeysMap[orgName].UserKeys.Light[j]
	case "client":
		if j >= len(p11KeysConfig.P11KeysMap[orgName].UserKeys.Client) {
			return fmt.Errorf("pkcs11 key not set, orgName = %s, clientId = %d", orgName, j)
		}
		keyLabel = p11KeysConfig.P11KeysMap[orgName].UserKeys.Client[j]
	case "consensus":
		if j >= len(p11KeysConfig.P11KeysMap[orgName].NodeKeys.Consensus) {
			return fmt.Errorf("pkcs11 key not set, orgName = %s, consensusId = %d", orgName, j)
		}
		keyLabel = p11KeysConfig.P11KeysMap[orgName].NodeKeys.Consensus[j]
	case "common":
		if j >= len(p11KeysConfig.P11KeysMap[orgName].NodeKeys.Common) {
			return fmt.Errorf("pkcs11 key not set, orgName = %s, commonId = %d", orgName, j)
		}
		keyLabel = p11KeysConfig.P11KeysMap[orgName].NodeKeys.Common[j]
	default:
		return fmt.Errorf("pkcs11 key not set, orgName = %s, id = %d, usage = %s", orgName, j, usage)
	}
	keyId, keyPwd := getKeyIdPwd(keyLabel)
	cert.P11Context.WithPrivKeyId(keyId).WithPrivKeyPwd(keyPwd).WithPrivKeyType(keyType)
	return nil
}

func getKeyIdPwd(label string) (keyId, keyPwd string) {
	key := strings.Split(label, ",")
	if len(key) == 0 {
		return "", ""
	} else if len(key) == 1 {
		return strings.TrimSpace(key[0]), ""
	}
	return strings.TrimSpace(key[0]), strings.TrimSpace(key[1])
}
