This is the vote contract for test net demo.

## The description of methods are below:
## 1. issueProject
### args: 
#### key: "projectInfo"
#### value: json
#### example:
```json
{"projectInfo":{"Id":"projectId1","PicUrl":"www.sina.com","Title":"wonderful","StartTime":"1664450291","EndTime":"1665314291","Desc":"the 1","Items":[{"Id":"item1","PicUrl":"www.baidu.com","Desc":"beautiful", "Url":"www.qq.com"},{"Id":"item2","PicUrl":"www.baidu.com","Desc":"beautiful", "Url":"www.qq.com"}]}}
```
### event:
#### topic: issue project
#### data: same as args

## 2. vote
### args:
#### key1: "projectId"
#### value1: string
#### key2: "itemId"
#### value2: string
#### example:
```json
{"projectId":"projectId1","itemId":"item1"}
```
### event:
#### topic: vote
#### data: projectId, itemId, voter
#### example:
```json
["projectId1","itemId1","441224c1757ec1cc67f4b7b3ac29c76cf5799ee5"]
```

## 3. queryProjectVoters
### args:
#### key1: "projectId"
#### value1: string
#### example:
```json
{"projectId":"projectId1"}
```
#### resp exampl:
```json
{"ProjectId":"projectId1","ItemVotes":[{"ItemId":"item1","VotesCount":1,"Voters":["441224c1757ec1cc67f4b7b3ac29c76cf5799ee5"]},{"ItemId":"item2","VotesCount":1,"Voters":["441224c1757ec1cc67f4b7b3ac29c76cf5799ee5"]}]}
```

## 4. queryProjectItemVoters
### args:
#### key1: "projectId"
#### value1: string
#### key2: "itemId"
#### value2: string
#### example:
```json
{"projectId":"projectId1","itemId":"item1"}
```
#### resp example:
```json
{"ItemId":"item1","VotesCount":1,"Voters":["441224c1757ec1cc67f4b7b3ac29c76cf5799ee5"]}
```

## Test

### 安装合约
```sh
./cmc client contract user create \
--contract-name=vote \
--version=1.0 \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--byte-code-path=./testdata/vote.7z \
--runtime-type=DOCKER_GO \
--admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key --params="{}"
```

### 发布投票项目
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=issueProject \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"projectInfo\":{\"Id\":\"projectId2\",\"PicUrl\":\"www.sina.com\",\"Title\":\"wonderful\",\"StartTime\":\"1664450291\",\"EndTime\":\"1665314291\",\"Desc\":\"the 1\",\"Items\":[{\"Id\":\"item1\",\"PicUrl\":\"www.baidu.com\",\"Desc\":\"beautiful\", \"Url\":\"www.qq.com\"},{\"Id\":\"item2\",\"PicUrl\":\"www.baidu.com\",\"Desc\":\"beautiful\", \"Url\":\"www.qq.com\"}]}}"
```

### 投票
### 测试Case1: 
投票的projectId或itemId不存在时应报错
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=vote \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"projectId\":\"projectId1\",\"itemId\":\"item1\"}"
```

### 查询项目投票人员
### 测试Case1:
查询投票项目的projectId不存在时应报错
### 测试Case2:
查询投票项目的projectId存在时应返回正确的投票人员集合
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=queryProjectVoters \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"projectId\":\"projectId1\"}"
```

### 查询投票项投票人员
### 测试Case1:
查询投票项目的projectId或itemId不存在时应报错
### 测试Case2:
查询投票项目的projectId和itemId存在时应返回正确的投票人员集合
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=queryProjectItemVoters \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"projectId\":\"projectId1\",\"itemId\":\"item1\"}"
```

