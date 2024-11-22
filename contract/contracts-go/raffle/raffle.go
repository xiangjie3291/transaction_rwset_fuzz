/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

const queryErrorMsg = "get peoples data from store failed"

// RaffleContract demo contract for raffle
type RaffleContract struct {
}

// People raffle people information
type People struct {
	Num  int    `json:"num"`
	Name string `json:"name"`
}

// Peoples raffle peoples
type Peoples struct {
	//Peoples map[int]string `json:"peoples"`
	Peoples []*People `json:"peoples"`
}

// InitContract install contract func
func (f *RaffleContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract upgrade contract func
func (f *RaffleContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

// InvokeContract the entry func of invoke contract func
func (f *RaffleContract) InvokeContract(method string) protogo.Response {
	switch method {
	case "querySlow":
		return f.queryAndInit()
	case "evilPlayer":
		return f.evilPlayer()
	case "registerAll":
		return f.registerAll()
	case "query":
		return f.query()
	case "raffle":
		return f.raffle()
	default:
		return sdk.Error("invalid method")
	}
}

func (f *RaffleContract) queryAndInit() protogo.Response {
	peoplesData, err := sdk.Instance.GetStateByte("peoples", "init")
	if err != nil {
		sdk.Instance.Errorf(queryErrorMsg)
		return sdk.Error(queryErrorMsg)
	}
	time.Sleep(100 * time.Millisecond)

	if peoplesData == nil {
		peoplesData = []byte{}
	}

	err = sdk.Instance.PutStateByte("peoples", "init", peoplesData)

	if err != nil {
		sdk.Instance.Errorf(queryErrorMsg)
		return sdk.Error(queryErrorMsg)
	}
	return sdk.Success(peoplesData)
}

func (f *RaffleContract) evilPlayer() protogo.Response {
	peoplesData, err := sdk.Instance.GetStateByte("peoples", "init")

	if err != nil {
		return sdk.Error("failed to call get_state")
	}

	if peoplesData == nil {
		peoplesData = []byte{}
	}

	number := len(peoplesData)

	for i := 0; i < 1000000; i++ {
		for j := 0; j < 1000; j++ {
			number += int(math.Sqrt(float64(i * j)))
		}
	}

	sdk.Instance.Infof("[Current number] number=%d", number)

	return sdk.Success([]byte(fmt.Sprintf("[evilPlayer excute]")))
}

func (f *RaffleContract) registerAll() protogo.Response {
	params := sdk.Instance.GetArgs()

	// 获取参数
	value := params["peoples"]
	var errMsg string
	if len(value) == 0 {
		errMsg = "value should not be empty!"
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	sdk.Instance.Debugf("registerAll receive peoples value:%s", string(value))
	var peoples Peoples
	err := json.Unmarshal(value, &peoples)
	if err != nil {
		return sdk.Error(err.Error())
	}
	for i := 0; i < len(peoples.Peoples); i++ {
		if people := peoples.Peoples[i]; len(people.Name) == 0 {
			errMsg = fmt.Sprintf("[registerAll] name should not be empty for number %d", i)
			sdk.Instance.Errorf(errMsg)
			return sdk.Error(errMsg)
		}
	}
	err = sdk.Instance.PutStateByte("peoples", "", value)
	if err != nil {
		errMsg = fmt.Sprintf("[registerAll] put state bytes failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	// 返回结果
	return sdk.Success([]byte("ok"))
}

func (f *RaffleContract) raffle() protogo.Response {
	params := sdk.Instance.GetArgs()

	var errMsg string
	level := string(params["level"])
	if len(level) == 0 {
		errMsg = "level should not be empty!"
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	argTimestamp := string(params["timestamp"])
	if len(argTimestamp) == 0 {
		errMsg = "argTimestamp should not be empty!"
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}

	peoplesData, err := sdk.Instance.GetStateByte("peoples", "")
	if err != nil {
		sdk.Instance.Errorf(queryErrorMsg)
		return sdk.Error(queryErrorMsg)
	}

	var peoples Peoples
	err = json.Unmarshal(peoplesData, &peoples)
	if err != nil {
		errMsg = fmt.Sprintf("unmarshal peoples data failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	num := f.bkdrHash(argTimestamp)
	num = num % len(peoples.Peoples)

	sdk.Instance.Infof(fmt.Sprintf("raffle pos: %d", num))

	resultPeople := peoples.Peoples[num]
	result := fmt.Sprintf("num: %d, name: %s, level: %s", resultPeople.Num, resultPeople.Name, level)
	var newPeoples Peoples
	newPeoples.Peoples = append(newPeoples.Peoples, peoples.Peoples[0:num]...)
	if num+1 < len(peoples.Peoples) {
		newPeoples.Peoples = append(newPeoples.Peoples, peoples.Peoples[num+1:]...)
	}
	//delete(peoples.Peoples, num)
	newPeoplesData, err := json.Marshal(newPeoples)
	if err != nil {
		errMsg = fmt.Sprintf("marshal new peoples data failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	err = sdk.Instance.PutStateByte("peoples", "", newPeoplesData)
	if err != nil {
		errMsg = fmt.Sprintf("put new peoples data failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}

	return sdk.Success([]byte(result))
}

func (f *RaffleContract) bkdrHash(timestamp string) int {
	hash := 0
	seed := 131
	for x := range timestamp {
		hash = hash*seed + x
	}
	return hash & 0x7FFFFFFF
}

func (f *RaffleContract) query() protogo.Response {
	peoplesData, err := sdk.Instance.GetStateByte("peoples", "")
	if err != nil {
		sdk.Instance.Errorf(queryErrorMsg)
		return sdk.Error(queryErrorMsg)
	}

	return sdk.Success(peoplesData)
}

func main() {
	err := sandbox.Start(new(RaffleContract))
	if err != nil {
		sdk.Instance.Errorf(err.Error())
	}
}
