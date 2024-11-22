/*
  Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

  SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

const (
	projectInfoArgKey       = "projectInfo"
	projectIdArgKey         = "projectId"
	projectItemIdArgKey     = "itemId"
	projectInfoStoreKey     = "project"
	projectItemsStoreMapKey = "projectItems"
	projectVotesStoreMapKey = "projectVotes"
	trueString              = "1"
)

// VoteContract the vote contract
type VoteContract struct {
}

// ProjectInfo specific information about the project
type ProjectInfo struct {
	Id        string
	PicUrl    string
	Title     string
	StartTime string
	EndTime   string
	Desc      string
	Items     []*ItemInfo
}

// ItemInfo specific information about the item
type ItemInfo struct {
	Id     string
	PicUrl string
	Desc   string
	Url    string
}

// ProjectVotesInfo a collection of votes for projects
type ProjectVotesInfo struct {
	ProjectId string
	ItemVotes []*ItemVotesInfo
}

// ItemVotesInfo specific information about the vote
type ItemVotesInfo struct {
	ItemId     string
	VotesCount int
	Voters     []string
}

// InitContract install contract func
func (c *VoteContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract upgrade contract func
func (c *VoteContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

// InvokeContract the entry func of invoke contract func
func (c *VoteContract) InvokeContract(method string) protogo.Response {
	if len(method) == 0 {
		return sdk.Error("method of param should not be empty")
	}

	switch method {
	case "issueProject":
		return c.issueProject()
	case "vote":
		return c.vote()
	case "queryProjectVoters":
		return c.queryProjectVoters()
	case "queryProjectItemVoters":
		return c.queryProjectItemVoters()
	default:
		return sdk.Error("Invalid method")
	}
}

func (c *VoteContract) issueProject() protogo.Response {
	args := sdk.Instance.GetArgs()
	if len(args) == 0 {
		return sdk.Error("issueProject should have arg of " + projectInfoArgKey)
	}
	projectInfoBytes := args[projectInfoArgKey]
	var pi ProjectInfo
	if err := json.Unmarshal(projectInfoBytes, &pi); err != nil {
		return sdk.Error("unmarshal project info failed")
	}

	piBytes, err := sdk.Instance.GetStateByte(projectInfoStoreKey, pi.Id)
	if err != nil {
		return sdk.Error("get project info failed")
	}
	if len(piBytes) > 0 {
		return sdk.Error("project already exist")
	}
	err = sdk.Instance.PutStateByte(projectInfoStoreKey, pi.Id, projectInfoBytes)
	if err != nil {
		return sdk.Error("store project info failed")
	}

	sm, err := sdk.NewStoreMap(projectItemsStoreMapKey, 2, crypto.HASH_TYPE_SHA256)
	if err != nil {
		return sdk.Error(fmt.Sprintf("new store map of vote project failed, err: %s", err))
	}
	for _, id := range pi.Items {
		if err = sm.Set([]string{pi.Id, id.Id}, []byte(trueString)); err != nil {
			return sdk.Error(fmt.Sprintf("set project item to storemap failed, err:%s", err))
		}
	}

	sdk.Instance.EmitEvent("issue project", []string{string(projectInfoBytes)})

	return sdk.Success([]byte("issue project success"))
}

func (c *VoteContract) vote() protogo.Response {
	args := sdk.Instance.GetArgs()
	if len(args) < 2 {
		return sdk.Error("vote should have args of " + projectIdArgKey + " and " + projectItemIdArgKey)
	}
	projectId := string(args[projectIdArgKey])
	projectItemId := string(args[projectItemIdArgKey])
	if len(projectId) == 0 {
		return sdk.Error("invalid project id")
	}
	if len(projectItemId) == 0 {
		return sdk.Error("invalid project item id")
	}

	projectInfoBytes, err := sdk.Instance.GetStateByte(projectInfoStoreKey, projectId)
	if err != nil {
		return sdk.Error("get project info from store failed")
	}

	var pi ProjectInfo
	if err = json.Unmarshal(projectInfoBytes, &pi); err != nil {
		return sdk.Error("unmarshal project info from store failed")
	}

	txTime, err := sdk.Instance.GetTxTimeStamp()
	if err != nil {
		return sdk.Error("get timestamp of tx failed")
	}
	if txTime < pi.StartTime || txTime > pi.EndTime {
		return sdk.Error("vote time is not in valid scope")
	}

	sender, err := sdk.Instance.Sender()
	if err != nil {
		return sdk.Error(fmt.Sprintf("get sender failed, err: %s", err))
	}
	isVoted, err := sdk.Instance.GetState(projectId, sender)
	if err != nil {
		return sdk.Error(fmt.Sprintf("get is voted from store failed, err: %s", err))
	}
	if isVoted == trueString {
		return sdk.Error("you had already voted this project, please don't vote same project multi times")
	}
	if err = sdk.Instance.PutState(projectId, sender, trueString); err != nil {
		return sdk.Error(fmt.Sprintf("store voted info failed, err: %s", err))
	}

	smOfPi, err := sdk.NewStoreMap(projectItemsStoreMapKey, 2, crypto.HASH_TYPE_SHA256)
	if err != nil {
		return sdk.Error(fmt.Sprintf("new store map of vote project failed, err: %s", err))
	}
	val, err := smOfPi.Get([]string{projectId, projectItemId})
	if err != nil || string(val) != trueString {
		return sdk.Error("can't found project or item")
	}
	smOfPiVotes, err := sdk.NewStoreMap(projectVotesStoreMapKey, 3, crypto.HASH_TYPE_SHA256)
	if err != nil {
		return sdk.Error(fmt.Sprintf("new store map of project vote info failed, err: %s", err))
	}
	if err = smOfPiVotes.Set([]string{projectId, projectItemId, sender}, []byte(trueString)); err != nil {
		return sdk.Error(fmt.Sprintf("store map set vote failed, err: %s", err))
	}
	sdk.Instance.EmitEvent("vote", []string{projectId, projectItemId, sender})
	return sdk.Success([]byte("vote success"))
}

func (c *VoteContract) queryProjectVoters() protogo.Response {
	args := sdk.Instance.GetArgs()
	if len(args) == 0 {
		return sdk.Error("vote should have arg of " + projectIdArgKey)
	}
	pi := string(args[projectIdArgKey])
	if len(pi) == 0 {
		return sdk.Error("invalid project id")
	}
	smOfPiVotes, err := sdk.NewStoreMap(projectItemsStoreMapKey, 2, crypto.HASH_TYPE_SHA256)
	if err != nil {
		return sdk.Error(fmt.Sprintf("new store map of project vote info failed, err: %s", err))
	}
	rs, err := smOfPiVotes.NewStoreMapIteratorPrefixWithKey([]string{pi})
	if err != nil {
		return sdk.Error(fmt.Sprintf("new store map iterator of project info failed, err: %s", err))
	}
	projectVotesInfo := newProjectVotesInfo(pi)
	for {
		if !rs.HasNext() {
			break
		}
		var item string
		item, _, _, err = rs.Next()
		if err != nil {
			return sdk.Error(fmt.Sprintf("iterator next failed, err: %s", err))
		}
		itemId := strings.TrimPrefix(strings.TrimPrefix(item, projectItemsStoreMapKey), pi)
		if len(itemId) == 0 {
			return sdk.Error("invalid itemId")
		}
		var itemVotesInfo *ItemVotesInfo
		itemVotesInfo, err = c.queryItemVoters(pi, itemId)
		if err != nil {
			return sdk.Error(err.Error())
		}
		projectVotesInfo.ItemVotes = append(projectVotesInfo.ItemVotes, itemVotesInfo)
	}

	var projectVotesInfoBytes []byte
	projectVotesInfoBytes, err = json.Marshal(projectVotesInfo)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success(projectVotesInfoBytes)
}

func (c *VoteContract) queryProjectItemVoters() protogo.Response {
	args := sdk.Instance.GetArgs()
	if len(args) < 2 {
		return sdk.Error("queryProjectItemVoters should have args of " + projectIdArgKey + " and " + projectItemIdArgKey)
	}
	pi := string(args[projectIdArgKey])
	if len(pi) == 0 {
		return sdk.Error("invalid project id")
	}
	projectItemId := string(args[projectItemIdArgKey])
	if len(projectItemId) == 0 {
		return sdk.Error("invalid project item id")
	}
	itemVotesInfo, err := c.queryItemVoters(pi, projectItemId)
	if err != nil {
		return sdk.Error(err.Error())
	}
	itemVotesInfoBytes, err := json.Marshal(itemVotesInfo)
	if err != nil {
		return sdk.Error("marshal itemVotesInfo failed")
	}
	return sdk.Success(itemVotesInfoBytes)
}

func (c *VoteContract) queryItemVoters(projectId, itemId string) (*ItemVotesInfo, error) {

	smOfPiVotes, err := sdk.NewStoreMap(projectVotesStoreMapKey, 3, crypto.HASH_TYPE_SHA256)
	if err != nil {
		return nil, fmt.Errorf("new store map of project vote info failed, err: %s", err)
	}
	rs, err := smOfPiVotes.NewStoreMapIteratorPrefixWithKey([]string{projectId, itemId})
	if err != nil {
		return nil, fmt.Errorf("new store map iterator of project info failed, err: %s", err)
	}
	itemVotesInfo := newItemVotesInfo(itemId)

	for {
		if !rs.HasNext() {
			break
		}
		item, _, _, err := rs.Next()
		if err != nil {
			return nil, fmt.Errorf("iterator next failed, err: %s", err)
		}
		voter := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(item, projectVotesStoreMapKey), projectId), itemId)
		if len(voter) == 0 {
			return nil, fmt.Errorf("found invalid voter")
		}
		itemVotesInfo.Voters = append(itemVotesInfo.Voters, voter)
		itemVotesInfo.VotesCount++
	}
	return itemVotesInfo, nil
}

func newItemVotesInfo(itemId string) *ItemVotesInfo {
	return &ItemVotesInfo{
		ItemId: itemId,
		Voters: make([]string, 0, 5),
	}
}

func newProjectVotesInfo(projectId string) *ProjectVotesInfo {
	return &ProjectVotesInfo{
		ProjectId: projectId,
		ItemVotes: make([]*ItemVotesInfo, 0, 5),
	}
}

func main() {
	err := sandbox.Start(new(VoteContract))
	if err != nil {
		sdk.Instance.Errorf(err.Error())
	}
}
