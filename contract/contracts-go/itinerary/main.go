package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// Itinerary 行程定义，包含所在公网IP、网络运营商、国家省份城市、经纬度坐标、所在时区
type Itinerary struct {
	IP       string  `json:"ip"`
	City     string  `json:"city"`
	Region   string  `json:"region"`
	Country  string  `json:"country"`
	Loc      string  `json:"loc"`
	Org      string  `json:"org"`
	Timezone string  `json:"timezone"`
	Asn      Asn     `json:"asn"`
	Company  Company `json:"company"`
	Privacy  Privacy `json:"privacy"`
	Abuse    Abuse   `json:"abuse"`
	Domains  Domains `json:"domains"`
}

type Asn struct {
	Asn    string `json:"asn"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Route  string `json:"route"`
	Type   string `json:"type"`
}
type Company struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Type   string `json:"type"`
}
type Privacy struct {
	Vpn     bool   `json:"vpn"`
	Proxy   bool   `json:"proxy"`
	Tor     bool   `json:"tor"`
	Relay   bool   `json:"relay"`
	Hosting bool   `json:"hosting"`
	Service string `json:"service"`
}
type Abuse struct {
	Address string `json:"address"`
	Country string `json:"country"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Network string `json:"network"`
	Phone   string `json:"phone"`
}
type Domains struct {
	Total   int           `json:"total"`
	Domains []interface{} `json:"domains"`
}

// HistoryValue 针对key的历史记录定义
type HistoryValue struct {
	Field       string      `json:"field"`
	Value       interface{} `json:"value"`
	TxId        string      `json:"txId"`
	Timestamp   string      `json:"timestamp"`
	BlockHeight int         `json:"blockHeight"`
	Key         string      `json:"key"`
}

type Result struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ToResult(data interface{}) []byte {
	dataRet := &Result{
		Code: 200,
		Msg:  "",
		Data: data,
	}
	res, err := json.Marshal(dataRet)
	if err != nil {
		res = []byte("json marshal error")
	}
	return res
}

// ItineraryContract 行程行踪追溯合约
// 在疫情肆意横行的当下，通信大数据行程卡，也叫行程码，是由中国信通院联合中国电信、中国移动、中国联通三家基础电信企业利用手机接收的数据，
// 通过用户手机所处的最近的基站位置获取（基于基站和手机之间的通信时差，可以确定出手机持有人的大概位置；而想要通过基站确定手机的位置，
// 最为精准的就是“三角定位”，即通过三个基站直接确定出手机的位置。即便手机持有人拔掉了手机卡，基站也会利用手机独有的IMEI号，针对手机进行定位），
// 为全国16亿手机用户免费提供的查询服务，手机用户可通过服务，查询本人前14天到过的所有国家城市地区记录信息，还有在国内停留4小时以上的城市信息；
// 疫情防控的一个基本要务是溯源，追踪，也就是可追溯性，而对于用户或个体的品质并不在意，因为其中的前提是，能够出来的，能够流通的，必须是好的，
// 否则相关组织或个人，就需要为此负责，当获取行程码后，行程卡会根据用户的近期行程显示颜色，一共四种，分别是绿卡、黄卡、红卡；在疫情反复不断的情况下，
// 行程码可以反映持有者到底是否去过危险区域，是疫情防控的一个重要手段。
type ItineraryContract struct {
}

// InitContract install contract func
func (t *ItineraryContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract upgrade contract func
func (t *ItineraryContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

// InvokeContract the entry func of invoke contract func
func (t *ItineraryContract) InvokeContract(method string) protogo.Response {
	switch method {
	case "save":
		return t.save()
	case "queryHistory":
		return t.queryHistory()
	default:
		return sdk.Error("invalid method")
	}
}

func (t *ItineraryContract) save() protogo.Response {
	params := sdk.Instance.GetArgs()
	itineraryStr := string(params["itinerary"])
	phone := string(params["phone"])

	if isBlank(itineraryStr) {
		errMsg := "'itinerary' should not be empty!"
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}

	if isBlank(phone) {
		errMsg := "'phone' should not be empty!"
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}

	var itinerary Itinerary
	err := json.Unmarshal([]byte(itineraryStr), &itinerary)
	if err != nil {
		errMsg := fmt.Sprintf("unmarshal itinerary data failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}

	if isBlank(itinerary.IP) || isBlank(itinerary.Country) || isBlank(itinerary.City) || isBlank(itinerary.Region) {
		errMsg := fmt.Sprintf("'ip','country','city','region' should not be empty, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}

	itineraryDataBytes, err := json.Marshal(itinerary)
	if err != nil {
		errMsg := fmt.Sprintf("json.marshal new itinerary data failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}

	err = sdk.Instance.PutStateByte(phone, "", itineraryDataBytes)
	if err != nil {
		errMsg := fmt.Sprintf("put new itinerary data failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}

	// start event
	sdk.Instance.EmitEvent(phone, []string{string(itineraryDataBytes)})

	return sdk.Success(ToResult(itinerary))
}

// queryHistory 此方法需要依赖 chainmaker.yml 的 history db
func (t *ItineraryContract) queryHistory() protogo.Response {
	params := sdk.Instance.GetArgs()
	phone := string(params["phone"])

	if isBlank(phone) {
		errMsg := "'phone' should not be empty!"
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}

	iter, err := sdk.Instance.NewHistoryKvIterForKey(phone, "")
	if err != nil {
		errMsg := fmt.Sprintf("new HistoryKvIter for key=[%s] failed, %s", phone, err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}

	var itinerary Itinerary
	recordMap := make(map[string]HistoryValue, 0)
	for iter.HasNext() {
		km, err := iter.Next()
		if err != nil {
			errMsg := "iterator failed to get the next element" + "," + err.Error()
			sdk.Instance.Errorf(errMsg)
			// 避免出现EOF，暂时跳过
			continue
		}

		err = json.Unmarshal(km.Value, &itinerary)
		if err != nil {
			errMsg := "json parse element error" + "," + err.Error()
			sdk.Instance.Errorf(errMsg)
			continue
		}

		// convert
		time64, _ := strconv.ParseInt(km.Timestamp, 10, 64)
		hv := &HistoryValue{
			TxId:        km.TxId,
			Timestamp:   Time2Str(time.Unix(time64, 0), "2006-01-02 15:04:05.000"),
			BlockHeight: km.BlockHeight,
			Key:         km.Key,
			Field:       km.Field,
			Value:       itinerary,
		}

		location := itinerary.Country + "-" + itinerary.City + "-" + itinerary.Region
		if record, ok := recordMap[location]; !ok {
			recordMap[location] = *hv
		} else {
			// 只要最新行踪记录
			if record.BlockHeight < km.BlockHeight {
				recordMap[location] = *hv
			}
		}
	}

	closed, err := iter.Close()
	if !closed || err != nil {
		errMsg := fmt.Sprintf("iterator close failed, %s", err.Error())
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}

	return sdk.Success(ToResult(recordMap))
}

func isBlank(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func Time2Str(aTime time.Time, pattern string) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return aTime.In(loc).Format(pattern)
}

func main() {
	err := sandbox.Start(new(ItineraryContract))
	if err != nil {
		log.Fatal(err)
	}
}
