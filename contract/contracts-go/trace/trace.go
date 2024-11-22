package main

import (
	"encoding/json"
	"fmt"
	"log"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

const (
	goodsIdArgKey    = "goodsId"
	nameArgKey       = "name"
	factoryArgKey    = "factory"
	fromArgKey       = "from"
	toArgKey         = "to"
	uploaderArgKey   = "uploader"
	sellerArgKey     = "seller"
	goodsStoreMapKey = "goodsList"
	statusCreate     = iota
	statusTransfer
	statusUpload
	statusSelled
)

// Trace 溯源
type Trace interface {

	//创建商品
	newGoods() protogo.Response

	//运输商品
	transferGoods() protogo.Response

	//上架商品
	uploadGoods() protogo.Response

	//销售商品
	sellGoods() protogo.Response

	//获取商品当前状态
	getGoodsStatus() protogo.Response

	//获取溯源信息
	getTraceInfo() protogo.Response
}

var _ Trace = (*TraceContract)(nil)

// TraceContract 溯源合约
type TraceContract struct {
}

func (f *TraceContract) newGoods() protogo.Response {
	args := sdk.Instance.GetArgs()
	if len(args) < 3 {
		return sdk.Error("newGoods should have arg of " + goodsIdArgKey + " and " + nameArgKey + " and " + factoryArgKey)
	}

	goodsId := string(args[goodsIdArgKey])
	if len(goodsId) == 0 {
		return sdk.Error("invalid goodsId")
	}

	name := string(args[nameArgKey])
	if len(name) == 0 {
		return sdk.Error("invalid name")
	}

	factory := string(args[factoryArgKey])
	if len(factory) == 0 {
		return sdk.Error("invalid factory")
	}

	goodsBytes, err := sdk.Instance.GetStateByte(goodsStoreMapKey, goodsId)
	if err != nil {
		return sdk.Error(fmt.Sprintf("newGoods GetStateByte error : %s", err))
	}
	if len(goodsBytes) > 0 {
		return sdk.Error("goodsId is exsits!")
	}

	goods := &Goods{
		GoodsId:    goodsId,
		Name:       name,
		Status:     statusCreate,
		TraceDatas: []*TraceData{},
	}

	id, err := sdk.Instance.Sender()
	if err != nil {
		return sdk.Error(fmt.Sprintf("newGoods GetSenderOrgId failed, err: %s", err))
	}
	time, err := sdk.Instance.GetTxTimeStamp()
	if err != nil {
		return sdk.Error(fmt.Sprintf("newGoods GetTxTimeStamp failed, err: %s", err))
	}

	goods.TraceDatas = append(goods.TraceDatas, &TraceData{
		Operator:     id,
		Status:       statusCreate,
		OperatorTime: time,
		Remark:       goodsId + ":" + factory + " created"})

	goodsBytes, err = json.Marshal(goods)
	if err != nil {
		return sdk.Error(fmt.Sprintf("newGoods Marshal failed, err: %s", err))
	}

	err = sdk.Instance.PutStateByte(goodsStoreMapKey, goodsId, goodsBytes)
	if err != nil {
		return sdk.Error(fmt.Sprintf("newGoods PutStateByte failed, err: %s", err))
	}

	return sdk.Success([]byte("newGoods success"))
}

func (f *TraceContract) transferGoods() protogo.Response {
	args := sdk.Instance.GetArgs()
	if len(args) < 3 {
		return sdk.Error("transferGoods should have arg of " + goodsIdArgKey + " and " + fromArgKey + " and " + toArgKey)
	}

	goodsId := string(args[goodsIdArgKey])
	if len(goodsId) == 0 {
		return sdk.Error("invalid goodsId")
	}

	from := string(args[fromArgKey])
	if len(from) == 0 {
		return sdk.Error("invalid from")
	}

	to := string(args[toArgKey])
	if len(to) == 0 {
		return sdk.Error("invalid to")
	}

	traceData := &TraceData{Status: statusTransfer, Remark: goodsId + ":" + from + "->" + to}

	response := updateGoodsStatus(goodsId, statusTransfer, traceData, "transferGoods")

	return response

}

func (f *TraceContract) uploadGoods() protogo.Response {
	args := sdk.Instance.GetArgs()
	if len(args) < 2 {
		return sdk.Error("uploadGoods should have arg of " + goodsIdArgKey + " and " + uploaderArgKey)
	}

	goodsId := string(args[goodsIdArgKey])
	if len(goodsId) == 0 {
		return sdk.Error("invalid goodsId")
	}

	uploader := string(args[uploaderArgKey])
	if len(uploader) == 0 {
		return sdk.Error("invalid uploader")
	}

	traceData := &TraceData{Status: statusUpload, Remark: goodsId + ":" + uploader + " upload"}

	response := updateGoodsStatus(goodsId, statusUpload, traceData, "uploadGoods")

	return response

}

func (f *TraceContract) sellGoods() protogo.Response {
	args := sdk.Instance.GetArgs()
	if len(args) < 2 {
		return sdk.Error("sellGoods should have arg of " + goodsIdArgKey + " and " + sellerArgKey)
	}

	goodsId := string(args[goodsIdArgKey])
	if len(goodsId) == 0 {
		return sdk.Error("invalid goodsId")
	}

	seller := string(args[sellerArgKey])
	if len(seller) == 0 {
		return sdk.Error("invalid seller")
	}

	traceData := &TraceData{Status: statusSelled, Remark: goodsId + ":" + "selled by " + seller}

	response := updateGoodsStatus(goodsId, statusSelled, traceData, "sellGoods")

	return response

}

func (f *TraceContract) getGoodsStatus() protogo.Response {
	args := sdk.Instance.GetArgs()
	goods, errResponse := getGoodsByArgs(args, "getGoodsStatus")
	if errResponse.GetStatus() == sdk.ERROR {
		return errResponse
	}

	var statusBytes []byte
	statusBytes, err := json.Marshal(goods.Status)
	if err != nil {
		return sdk.Error(fmt.Sprintf("getGoodsStatus Marshal goods.Status failed, err: %s", err))
	}
	return sdk.Success(statusBytes)

}

func (f *TraceContract) getTraceInfo() protogo.Response {
	args := sdk.Instance.GetArgs()
	goods, errResponse := getGoodsByArgs(args, "getTraceInfo")
	if errResponse.GetStatus() == sdk.ERROR {
		return errResponse
	}
	var traceDataBytes []byte
	traceDataBytes, err := json.Marshal(goods.TraceDatas)
	if err != nil {
		return sdk.Error(fmt.Sprintf("getTraceInfo Marshal goods.TraceDatas failed, err: %s", err))
	}
	return sdk.Success(traceDataBytes)

}

func updateGoodsStatus(goodsId string, status uint8, trace *TraceData, method string) protogo.Response {
	goodsBytes, err := sdk.Instance.GetStateByte(goodsStoreMapKey, goodsId)
	if err != nil {
		return sdk.Error(fmt.Sprintf(method+" GetStateByte error : %s", err))
	}
	if len(goodsBytes) == 0 {
		return sdk.Error("invalid goodsId, get goodsBytes len is 0")
	}

	var goods Goods
	err = json.Unmarshal(goodsBytes, &goods)
	if err != nil {
		return sdk.Error(fmt.Sprintf(method+" Unmarshal error : %s", err))
	}
	goods.Status = status

	time, err := sdk.Instance.GetTxTimeStamp()
	if err != nil {
		return sdk.Error(fmt.Sprintf(method+" GetTxTimeStamp failed, err: %s", err))
	}
	id, err := sdk.Instance.Sender()
	if err != nil {
		return sdk.Error(fmt.Sprintf(method+" GetSenderOrgId failed, err: %s", err))
	}

	trace.OperatorTime = time
	trace.Operator = id
	goods.TraceDatas = append(goods.TraceDatas, trace)

	//重新存储进去
	goodsBytes, err = json.Marshal(goods)
	if err != nil {
		return sdk.Error(fmt.Sprintf(method+" Marshal goods failed, err: %s", err))
	}
	err = sdk.Instance.PutStateByte(goodsStoreMapKey, goodsId, goodsBytes)

	if err != nil {
		return sdk.Error(fmt.Sprintf(method+" PutStateByte failed, err: %s", err))
	}
	return sdk.Success([]byte(method + " success"))
}

func getGoodsByArgs(args map[string][]byte, method string) (*Goods, protogo.Response) {
	if len(args) < 1 {
		return nil, sdk.Error(method + " should have arg of " + goodsIdArgKey)
	}

	goodsId := string(args[goodsIdArgKey])
	if len(goodsId) == 0 {
		return nil, sdk.Error("invalid goodsId")
	}

	goodsBytes, err := sdk.Instance.GetStateByte(goodsStoreMapKey, goodsId)
	if err != nil {
		return nil, sdk.Error(fmt.Sprintf(method+" GetStateByte error : %s", err))
	}
	if len(goodsBytes) == 0 {
		return nil, sdk.Error("invalid goodsId, get goodsBytes len is 0")
	}

	var goods *Goods
	err = json.Unmarshal(goodsBytes, &goods)
	if err != nil {
		return nil, sdk.Error(fmt.Sprintf(method+" Unmarshal error : %s", err))
	}
	return goods, sdk.Success([]byte("getGoodsByArgs success"))
}

// Goods 商品信息
type Goods struct {
	//商品id
	GoodsId string
	//商品名称
	Name string
	//商品状态
	Status uint8
	//商品溯源信息
	TraceDatas []*TraceData
}

// TraceData 溯源
type TraceData struct {
	//操作人
	Operator string
	//{0 :生产  1： 运输 2：上架  3：售卖}
	Status uint8
	//操作时间
	OperatorTime string
	//备注
	Remark string
}

// InitContract used to deploy and upgrade contract
func (f *TraceContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract used to upgrade contract
func (f *TraceContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

// InvokeContract used to invoke user contract
func (f *TraceContract) InvokeContract(method string) protogo.Response {
	switch method {
	case "newGoods":
		return f.newGoods()
	case "transferGoods":
		return f.transferGoods()
	case "uploadGoods":
		return f.uploadGoods()
	case "sellGoods":
		return f.sellGoods()
	case "goodsStatus":
		return f.getGoodsStatus()
	case "traceGoods":
		return f.getTraceInfo()
	default:
		return sdk.Error("invalid method")
	}
}
func main() {
	err := sandbox.Start(new(TraceContract))
	if err != nil {
		log.Fatal(err)
	}
}
