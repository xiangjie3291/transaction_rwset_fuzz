/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"chainmaker.org/chainmaker-cryptogen/command"
	"chainmaker.org/chainmaker-cryptogen/config"
	"github.com/spf13/cobra"
)

func main() {
	mainCmd := &cobra.Command{
		Use: "chainmaker-cryptogen",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.LoadCryptoGenConfig(command.ConfigPath)
			if config.PKCS11Enabled() {
				if err := config.LoadPKCS11KeysConfig(command.P11KeysPath); err != nil {
					log.Fatalln(err)
					return
				}
			}
		},
	}
	mainFlags := mainCmd.PersistentFlags()
	mainFlags.StringVarP(&command.ConfigPath, "config", "c", "../config/crypto_config_template.yml", "specify config file path")
	mainFlags.StringVarP(&command.P11KeysPath, "pkcs11_keys", "p", "../config/pkcs11_keys.yml", "specify pkcs11 keys file path")

	mainCmd.AddCommand(command.ShowConfigCmd())
	mainCmd.AddCommand(command.GenerateCmd())
	mainCmd.AddCommand(command.ExtendCmd())
	mainCmd.AddCommand(command.GeneratePWKCmd())
	mainCmd.AddCommand(command.ExtendPWKCmd())
	mainCmd.AddCommand(command.GeneratePKCmd())
	mainCmd.AddCommand(command.ExtendPKCmd())

	if err := mainCmd.Execute(); err != nil {
		log.Fatalf("failed to execute, err = %s", err)
	}
}
