/*
  Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

  SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"chainmaker.org/chainmaker/contract-utils/standard"
	"github.com/golang/mock/gomock"
)

func TestCMNFAContract_MintBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstance := sdk.NewMockSDKInterface(ctrl)

	params := make(map[string][]byte)
	mockInstance.EXPECT().GetArgs().AnyTimes().Return(params)
	mockInstance.EXPECT().GetStateByte(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
	mockInstance.EXPECT().GetStateFromKey(gomock.Any()).AnyTimes().Return("", nil)
	mockInstance.EXPECT().PutStateFromKey(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	mockInstance.EXPECT().PutStateByte(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	mockInstance.EXPECT().EmitEvent(gomock.Any(), gomock.Any()).AnyTimes().Return()
	sdk.Instance = mockInstance

	type args struct {
		Tokens []standard.NFA `json:"tokens"`
	}
	tests := []struct {
		Name    string `json:"name"`
		Args    args   `json:"args"`
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			Name: "test1",
			Args: args{
				Tokens: []standard.NFA{
					{
						TokenId:      "xxxxx",
						CategoryName: "111",
						To:           "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
						Metadata:     []byte("aaa"),
					},
					{
						TokenId:      "xxxxx1",
						CategoryName: "111",
						To:           "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
						Metadata:     []byte("aaa"),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			c := &CMNFAContract{}
			marshalBytes, _ := json.Marshal(tt.Args.Tokens)
			fmt.Println(string(marshalBytes))
			tokens := make([]standard.NFA, 0)
			err := json.Unmarshal(marshalBytes, &tokens)
			if err != nil {
				fmt.Println(err)
			}
			if err := c.MintBatch(tokens); (err != nil) != tt.wantErr {
				t.Errorf("MintBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCMNFAContract_CreateOrSetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstance := sdk.NewMockSDKInterface(ctrl)

	params := make(map[string][]byte)
	mockInstance.EXPECT().GetArgs().AnyTimes().Return(params)
	mockInstance.EXPECT().GetStateByte(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
	mockInstance.EXPECT().GetStateFromKey(gomock.Any()).AnyTimes().Return("", nil)
	mockInstance.EXPECT().PutStateFromKey(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	mockInstance.EXPECT().PutStateByte(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	mockInstance.EXPECT().EmitEvent(gomock.Any(), gomock.Any()).AnyTimes().Return()
	sdk.Instance = mockInstance
	type args struct {
		category *standard.Category
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				category: &standard.Category{
					CategoryName: "1111",
					CategoryURI:  "www.chainmaker.org.cn",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CMNFAContract{}
			marshalBytes, _ := json.Marshal(tt.args.category)
			fmt.Println(string(marshalBytes))
			if err := c.CreateOrSetCategory(tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrSetCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
