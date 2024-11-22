# 行程行踪追溯合约

## 一、介绍

在疫情肆意横行的当下，通信大数据行程卡，也叫行程码，是由中国信通院联合中国电信、中国移动、中国联通三家基础电信企业利用手机接收的数据，
通过用户手机所处的最近的基站位置获取 ,为全国16亿手机用户免费提供的查询服务，手机用户可通过服务，查询本人前14天到过的所有国家城市地区记录信息，
还有在国内停留4小时以上的城市信息。

疫情防控的一个基本要务是溯源，追踪，也就是可追溯性，而对于用户或个体的品质并不在意，因为其中的前提是，能够出来的，能够流通的，必须是好的，
否则相关组织或个人，就需要为此负责，当获取行程码后，行程卡会根据用户的近期行程显示颜色，一共四种，分别是绿卡、黄卡、红卡；在疫情反复不断的情况下，
行程码可以反映持有者到底是否去过危险区域，是疫情防控的一个重要手段。

## 二、特点

溯源、可追溯性

## 三、实现原理

明确了溯源可追溯性的特点之后，对于行程行踪的追溯，我们要知道溯源（也就是位置，定位），目前大家知道的定位有4种：

### 1. 定位

#### (1) 基站定位

基于基站和手机之间的通信时差，可以确定出手机持有人的大概位置；而想要通过基站确定手机的位置， 最为精准的就是“三角定位”，
即通过三个基站直接确定出手机的位置。申请的运营商的地址判断你的大致位置，注意只是大致位置。最精确可以获取到街道

#### (2) 网络IP定位

IP 地址是由一个叫互联网服务提供商，即 ISP 提供的，三层 ISP 结构分为主干 ISP，地区 ISP，本地 ISP。本地 ISP 给用户提供最直接的服务，
本地 ISP 可以连接到地区 ISP，也可以连接到主干 ISP，这样运营商判断你的大致位置。服务提供商将用户使用时的当前的坐标，跟 IP
地址数据，当前时间，定位方式，WI-FI 信息，
移动联通电信等运营商的基站信息，传送给后台服务器中，后台服务器留存了这些数据。再通过利用这些数据，就可以计算出一个 IP
曾经在哪些范围被使用过，
从而得到一个精确的范围数据，这个范围的中心点，就被认为是最接近用户的地点。

#### (3) 卫星定位

卫星可获取到目标对经纬度定位

#### (4) 手机内置定位

手机当中有一个IMEI识别码，识别码主要由IMEI、MEID或者S/N构成，这个识别码无论手机是否开机，卡槽中是否有卡，对它都没有影响。
在手机关机和没卡的时候，它依旧会被基站识别出来，识别出来后，基站再通过大数据进行提取分析，就可以轻而易举的找到手机的所在处了
（即便手机持有人拔掉了手机卡，基站也会利用手机独有的IMEI号，针对手机进行定位），

#### 结论：

从可行性和难易程度考虑，最终选择"网络IP定位"

### 2. 行程查询

行程的追溯，通过手机号查询某人到达各个区域的历史记录。
同一个区域（国家、省份、区域均相同）只需要获取最新的记录。

## 四、使用说明

### 1. 行程上报

模拟设备（手机）定时上报公网IP的位置信息，这里使用`https://ipinfo.io`提供的公网IP位置信息

#### method: `save`

#### arg1: `phone`,表示手机号

#### arg2: `itinerary`,表示公网IP的位置信息JSON结构体

#### example:

```json
{
  "phone": "18892352495",
  "itinerary": {
    "ip": "117.107.131.195",
    "city": "Beijing",
    "region": "Beijing",
    "country": "CN",
    "loc": "39.9075,116.3972",
    "org": "",
    "timezone": "Asia/Shanghai",
    "asn": {
      "asn": "AS4847",
      "name": "China Networks Inter-Exchange",
      "domain": "bta.net.cn",
      "route": "117.107.128.0/18",
      "type": "isp"
    },
    "company": {
      "name": "Beijing Sinnet Technology Co., Ltd.",
      "domain": "ghidc.net",
      "type": "business"
    },
    "privacy": {
      "vpn": false,
      "proxy": false,
      "tor": false,
      "relay": false,
      "hosting": false,
      "service": ""
    },
    "abuse": {
      "address": "Beijing, China",
      "country": "CN",
      "email": "ipas@cnnic.cn",
      "name": "Chen hao",
      "network": "117.107.128.0/17",
      "phone": "+86-13311166160"
    },
    "domains": {
      "total": 0,
      "domains": []
    }
  }
}
```

### event:

#### topic:`phone`

#### data: `itinerary`

### 2. 行程查询

通过手机号查询某人到达过所有的区域记录

#### method: `queryHistory`

#### arg1: `phone`,表示手机号

#### example:

```json
{
  "phone": "18892352495"
}
```

#### response:

```json
[
  {
    "CN-Beijing-Beijing": {
      "field": "",
      "value": {
        "ip": "117.107.131.195",
        "city": "Beijing",
        "region": "Beijing",
        "country": "CN",
        "loc": "39.9075,116.3972",
        "org": "",
        "timezone": "Asia/Shanghai",
        "asn": {
          "asn": "AS4847",
          "name": "China Networks Inter-Exchange",
          "domain": "bta.net.cn",
          "route": "117.107.128.0/18",
          "type": "isp"
        },
        "company": {
          "name": "Beijing Sinnet Technology Co., Ltd.",
          "domain": "ghidc.net",
          "type": "business"
        },
        "privacy": {
          "vpn": false,
          "proxy": false,
          "tor": false,
          "relay": false,
          "hosting": false,
          "service": ""
        },
        "abuse": {
          "address": "Beijing, China",
          "country": "CN",
          "email": "ipas@cnnic.cn",
          "name": "Chen hao",
          "network": "117.107.128.0/17",
          "phone": "+86-13311166160"
        },
        "domains": {
          "total": 0,
          "domains": [
          ]
        }
      },
      "txId": "1720ec5606b42398ca2be955cfe213fd0094da835e3140f5bd08f8e7b5f969c8",
      "timestamp": "2022-10-24 14:08:36.000",
      "blockHeight": 8630,
      "key": "18891401498"
    }
  }
]
```
