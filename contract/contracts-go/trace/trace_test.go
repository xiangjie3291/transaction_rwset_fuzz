package main

import (
	"reflect"
	"testing"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"github.com/golang/mock/gomock"
)

func TestTraceContract_InitContract(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TraceContract{}
			if got := c.InitContract(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTraceContract_InvokeContract(t *testing.T) {
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
			c := &TraceContract{}
			if got := c.InvokeContract(tt.args.method); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvokeContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTraceContract_UpgradeContract(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TraceContract{}
			if got := c.UpgradeContract(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpgradeContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTraceContract_new(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstance := sdk.NewMockSDKInterface(ctrl)

	args := make(map[string][]byte)
	args[goodsIdArgKey] = []byte("{\"goodsId\":\"1000\"}")
	args[nameArgKey] = []byte("{\"name\":\"apple\"}")
	args[factoryArgKey] = []byte("{\"factory\":\"chainMakerFC\"}")
	mockInstance.EXPECT().Sender()
	mockInstance.EXPECT().GetTxTimeStamp()
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
			name: "newGoods",
			want: protogo.Response{
				Status:  0,
				Message: "",
				Payload: []byte("newGoods success"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TraceContract{}
			if got := c.newGoods(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newGoods() = %v, want %v", got, tt.want)
			}
		})
	}
}
