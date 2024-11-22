/*
  Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

  SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"github.com/golang/mock/gomock"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
)

func TestRaffleContract_InitContract(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &RaffleContract{}
			if got := f.InitContract(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRaffleContract_InvokeContract(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &RaffleContract{}
			if got := f.InvokeContract(tt.args.method); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvokeContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRaffleContract_UpgradeContract(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &RaffleContract{}
			if got := f.UpgradeContract(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpgradeContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRaffleContract_bkdrHash(t *testing.T) {
	type args struct {
		timestamp string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &RaffleContract{}
			if got := f.bkdrHash(tt.args.timestamp); got != tt.want {
				t.Errorf("bkdrHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRaffleContract_query(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &RaffleContract{}
			if got := f.query(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRaffleContract_raffle(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &RaffleContract{}
			if got := f.raffle(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("raffle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRaffleContract_registerAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstance := sdk.NewMockSDKInterface(ctrl)

	args := make(map[string][]byte)
	args["peoples"] = []byte("{\"peoples\":[{\"num\":1,\"name\":\"Chris\"},{\"num\":2,\"name\":\"Linus\"}]}")
	mockInstance.EXPECT().GetArgs().AnyTimes().Return(args)
	mockInstance.EXPECT().GetStateByte(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
	mockInstance.EXPECT().PutStateByte(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	mockInstance.EXPECT().EmitEvent(gomock.Any(), gomock.Any()).AnyTimes().Return()
	sdk.Instance = mockInstance
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
		{
			name: "right",
			want: protogo.Response{
				Status:  0,
				Message: "",
				Payload: []byte("ok"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &RaffleContract{}
			if got := f.registerAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("registerAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_marshalAndUnmarshal(t *testing.T) {
	peoples := &Peoples{
		Peoples: []*People{
			{
				Num:  1,
				Name: "Chris",
			},
			{
				Num:  2,
				Name: "Linus",
			},
		},
	}
	peoplesBytes, err := json.Marshal(peoples)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(peoplesBytes))
	var peoplesUm Peoples
	err = json.Unmarshal(peoplesBytes, &peoplesUm)
	if err != nil {
		fmt.Println(err)
	}
	if !reflect.DeepEqual(peoples, &peoplesUm) {
		t.Errorf("projectInfo not equal between beform and after marshalled")
	}

	//var bytesSlice []byte
	//fmt.Printf("%v", bytesSlice == nil)
	//if piUm.Id != pi.Id || len(piUm.Items) != len(pi.Items) {
	//	t.Errorf("projectInfo not equal between beform and after marshalled")
	//}
	//for i := 0; i < len(pi.Items); i++ {
	//	if !deepEqual(pi.Items[i] != piUm.Items[i]) {
	//		t.Errorf("projectInfo not equal between beform and after marshalled")
	//	}
	//}
}
