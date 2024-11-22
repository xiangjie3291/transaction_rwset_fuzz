/*
  Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

  SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"reflect"
	"testing"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"chainmaker.org/chainmaker/contract-utils/safemath"
	"github.com/golang/mock/gomock"
)

func TestERC20Contract_InitContract(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstance := sdk.NewMockSDKInterface(ctrl)

	args := make(map[string][]byte)
	mockInstance.EXPECT().GetArgs().AnyTimes().Return(args)
	mockGetPutState(mockInstance)
	mockInstance.EXPECT().EmitEvent(gomock.Any(), gomock.Any()).AnyTimes().Return()
	mockInstance.EXPECT().Sender().AnyTimes().Return("studyzy", nil)

	sdk.Instance = mockInstance

	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
		{
			name: "test init contract",
			want: protogo.Response{
				Status:  sdk.OK,
				Payload: []byte("Init contract success"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			if got := c.InitContract(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

var stateMap map[string]string

func mockGetPutState(mockInstance *sdk.MockSDKInterface) {
	stateMap = make(map[string]string)
	mockInstance.EXPECT().PutState(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(k string, f string, val string) error {
			stateMap[k+"#"+f] = val
			return nil
		})
	mockInstance.EXPECT().GetState(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(k string, f string) (string, error) {
			val, ok := stateMap[k+"#"+f]
			if !ok {
				return "", nil
			}
			return val, nil
		})
}

func TestERC20Contract_InvokeContract(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstance := sdk.NewMockSDKInterface(ctrl)

	params := make(map[string][]byte)
	mockInstance.EXPECT().GetArgs().AnyTimes().Return(params)
	mockGetPutState(mockInstance)
	mockInstance.EXPECT().EmitEvent(gomock.Any(), gomock.Any()).AnyTimes().Return()
	mockInstance.EXPECT().Sender().AnyTimes().Return("B250CEF819", nil)
	sdk.Instance = mockInstance

	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want protogo.Response
	}{
		// TODO: Add test cases.
		{
			name: "totalSupply",
			args: args{
				method: "TotalSupply",
			},
			want: protogo.Response{
				Status:  sdk.OK,
				Payload: []byte("0"),
			},
		},
		{
			name: "balanceOf",
			args: args{
				method: "BalanceOf",
			},
			want: protogo.Response{
				Status:  sdk.ERROR,
				Message: "CMDFA: require parameter:'account'",
			},
		},
		{
			name: "transfer",
			args: args{
				method: "Transfer",
			},
			want: protogo.Response{
				Status:  sdk.ERROR,
				Message: "CMDFA: require parameter:'to'",
			},
		},
		{
			name: "transferFrom",
			args: args{
				method: "TransferFrom",
			},
			want: protogo.Response{
				Status:  sdk.ERROR,
				Message: "CMDFA: require parameter:'from'",
			},
		},
		{
			name: "approve",
			args: args{
				method: "Approve",
			},
			want: protogo.Response{
				Status:  sdk.ERROR,
				Message: "CMDFA: require parameter:'spender'",
			},
		},
		{
			name: "allowance",
			args: args{
				method: "Allowance",
			},
			want: protogo.Response{
				Status:  sdk.ERROR,
				Message: "CMDFA: require parameter:'spender'",
			},
		},
		{
			name: "name",
			args: args{
				method: "Name",
			},
			want: protogo.Response{
				Status:  sdk.OK,
				Payload: []byte("TestToken"),
			},
		},
		{
			name: "symbol",
			args: args{
				method: "Symbol",
			},
			want: protogo.Response{
				Status:  sdk.OK,
				Payload: []byte("TT"),
			},
		},
		{
			name: "decimals",
			args: args{
				method: "Decimals",
			},
			want: protogo.Response{
				Status:  sdk.OK,
				Payload: []byte("18"),
			},
		},
		{
			name: "mint",
			args: args{
				method: "Mint",
			},
			want: protogo.Response{
				Status:  sdk.ERROR,
				Message: "CMDFA: require parameter:'account'",
			},
		},
		{
			name: "invalid",
			args: args{
				method: "invalid",
			},
			want: protogo.Response{
				Status:  sdk.ERROR,
				Message: "Invalid method",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			c.InitContract()
			if got := c.InvokeContract(tt.args.method); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvokeContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestERC20Contract_UpgradeContract(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			if got := c.UpgradeContract(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpgradeContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestERC20Contract_allowance(t *testing.T) {
	type args struct {
		owner   string
		spender string
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
			c := &CmdfaContract{}
			if got, _ := c.Allowance(tt.args.owner, tt.args.spender); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("allowance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestERC20Contract_approve(t *testing.T) {
	type args struct {
		//owner   string
		spender string
		amount  *safemath.SafeUint256
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
			c := &CmdfaContract{}
			if got := c.Approve(tt.args.spender, tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("approve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestERC20Contract_balanceOf(t *testing.T) {
	type args struct {
		account string
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
			c := &CmdfaContract{}
			if got, _ := c.BalanceOf(tt.args.account); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("balanceOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestERC20Contract_decimals(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			if got, _ := c.Decimals(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decimals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestERC20Contract_getAllowance(t *testing.T) {
	type args struct {
		//allowanceInfo *sdk.StoreMap
		owner   string
		spender string
	}
	tests := []struct {
		name          string
		args          args
		wantAllowance *safemath.SafeUint256
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			gotAllowance, err := c.GetAllowance(tt.args.owner, tt.args.spender)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAllowance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAllowance, tt.wantAllowance) {
				t.Errorf("getAllowance() gotAllowance = %v, want %v", gotAllowance, tt.wantAllowance)
			}
		})
	}
}

func TestERC20Contract_getBalance(t *testing.T) {
	type args struct {
		//balanceInfo *sdk.StoreMap
		account string
	}
	tests := []struct {
		name        string
		args        args
		wantBalance *safemath.SafeUint256
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			gotBalance, err := c.GetBalance(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBalance, tt.wantBalance) {
				t.Errorf("getBalance() gotBalance = %v, want %v", gotBalance, tt.wantBalance)
			}
		})
	}
}

func TestERC20Contract_mint(t *testing.T) {
	type args struct {
		account string
		amount  *safemath.SafeUint256
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
			c := &CmdfaContract{}
			if got := c.Mint(tt.args.account, tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestERC20Contract_name(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			if got, _ := c.Name(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestERC20Contract_setAllowance(t *testing.T) {
	type args struct {
		//allowanceInfo *sdk.StoreMap
		owner     string
		spender   string
		allowance *safemath.SafeUint256
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			if err := c.SetAllowance(tt.args.owner, tt.args.spender, tt.args.allowance); (err != nil) != tt.wantErr {
				t.Errorf("setAllowance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestERC20Contract_setBalance(t *testing.T) {
	type args struct {
		//balanceInfo *sdk.StoreMap
		account string
		value   *safemath.SafeUint256
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			if err := c.SetBalance(tt.args.account, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("setBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestERC20Contract_spendAllowance(t *testing.T) {
	type args struct {
		owner   string
		spender string
		amount  *safemath.SafeUint256
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			if err := c.baseSpendAllowance(tt.args.owner, tt.args.spender, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("spendAllowance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestERC20Contract_symbol(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			if got, _ := c.Symbol(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("symbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestERC20Contract_totalSupply(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdfaContract{}
			if got, _ := c.TotalSupply(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("totalSupply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestERC20Contract_transfer(t *testing.T) {
	type args struct {
		//spender string
		to     string
		amount *safemath.SafeUint256
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
			c := &CmdfaContract{}
			if got := c.Transfer(tt.args.to, tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transfer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestERC20Contract_transferFrom(t *testing.T) {
	type args struct {
		//owner   string
		spender string
		to      string
		amount  *safemath.SafeUint256
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
			c := &CmdfaContract{}
			if got := c.TransferFrom(tt.args.spender, tt.args.to, tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transferFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}
