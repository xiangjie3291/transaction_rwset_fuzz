/*
  Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

  SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"github.com/golang/mock/gomock"
)

func TestVoteContract_InitContract(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &VoteContract{}
			if got := c.InitContract(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVoteContract_InvokeContract(t *testing.T) {
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
			c := &VoteContract{}
			if got := c.InvokeContract(tt.args.method); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvokeContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVoteContract_UpgradeContract(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &VoteContract{}
			if got := c.UpgradeContract(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpgradeContract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVoteContract_issueProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstance := sdk.NewMockSDKInterface(ctrl)

	args := make(map[string][]byte)
	args[projectInfoArgKey] = []byte("{\"Id\":\"projectId1\",\"PicUrl\":\"www.sina.com\",\"Title\":\"wonderful\"," +
		"\"StartTime\":\"1666234349\",\"EndTime\":\"1667098349\",\"Desc\":\"the 1\",\"Items\":[{\"Id\":\"item1\"," +
		"\"PicUrl\":\"www.baidu.com\",\"Desc\":\"beautiful\",\"Url\":\"www.qq.com\"},{\"Id\":\"item2\"," +
		"\"PicUrl\":\"www.baidu.com\",\"Desc\":\"beautiful\",\"Url\":\"www.qq.com\"}]}")
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
			name: "issueProject",
			want: protogo.Response{
				Status:  0,
				Message: "",
				Payload: []byte("issue project success"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &VoteContract{}
			if got := c.issueProject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("issueProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVoteContract_queryItemVoters(t *testing.T) {
	type args struct {
		projectId string
		itemId    string
	}
	tests := []struct {
		name    string
		args    args
		want    *ItemVotesInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &VoteContract{}
			got, err := c.queryItemVoters(tt.args.projectId, tt.args.itemId)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryItemVoters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryItemVoters() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVoteContract_queryProjectItemVoters(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &VoteContract{}
			if got := c.queryProjectItemVoters(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryProjectItemVoters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVoteContract_queryProjectVoters(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &VoteContract{}
			if got := c.queryProjectVoters(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryProjectVoters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVoteContract_vote(t *testing.T) {
	tests := []struct {
		name string
		want protogo.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &VoteContract{}
			if got := c.vote(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("vote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newItemVotesInfo(t *testing.T) {
	type args struct {
		itemId string
	}
	tests := []struct {
		name string
		args args
		want *ItemVotesInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newItemVotesInfo(tt.args.itemId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newItemVotesInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newProjectVotesInfo(t *testing.T) {
	type args struct {
		projectId string
	}
	tests := []struct {
		name string
		args args
		want *ProjectVotesInfo
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newProjectVotesInfo(tt.args.projectId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newProjectVotesInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_marshalAndUnmarshal(t *testing.T) {
	pi := &ProjectInfo{
		Id:        "projectId1",
		PicUrl:    "www.sina.com",
		Title:     "wonderful",
		StartTime: strconv.FormatInt(time.Now().Unix(), 10),
		EndTime:   strconv.FormatInt(time.Now().Add(time.Hour*240).Unix(), 10),
		Desc:      "the 1",
		Items: []*ItemInfo{
			{
				Id:     "item1",
				PicUrl: "www.baidu.com",
				Desc:   "beautiful",
				Url:    "www.qq.com",
			},
			{
				Id:     "item2",
				PicUrl: "www.baidu.com",
				Desc:   "beautiful",
				Url:    "www.qq.com",
			},
		},
	}
	piBytes, err := json.Marshal(pi)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(piBytes))
	var piUm ProjectInfo
	err = json.Unmarshal(piBytes, &piUm)
	if err != nil {
		fmt.Println(err)
	}
	if !reflect.DeepEqual(pi, &piUm) {
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
