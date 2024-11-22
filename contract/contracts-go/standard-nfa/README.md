https://git.chainmaker.org.cn/contracts/standard/-/blob/master/living/CM-CS-221221-NFA.md

## The description of methods are below:
## InitContract
### args:
#### key1: categoryName(optional)
#### value1: string
#### key2: categoryURI(optional)
#### value2: string
#### key3: admin(optional)
#### value3: string
#### example:
```json
{
  "categoryName": "chainmaker",
  "categoryURI": "http://chainmaker.org.cn",
  "admin": "huanle"
}
```

## UpgradeContract
### args:
#### key1: categoryName(optional)
#### value1: string
#### key2: categoryURI(optional)
#### value2: string
#### key3: admin(optional)
#### value3: string
#### example:
```json
{
  "categoryName": "chainmaker",
  "categoryURI": "http://chainmaker.org.cn",
  "admin": "huanle"
}
```

## Standards
### response example: {"standards":["CMNFA"]}

## Mint
### args:
#### key1: "to"
#### value1: string
#### key2: "tokenId"
#### value2: string
#### key3: "categoryName"
#### value2: string
#### key4: "metadata"(optional)
#### value3: bytes
#### example:
```json
{
  "to":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
  "tokenId":"111111111111111111111112",
  "categoryName":"huanle",
  "metadata": "url:https://chainmaker.org.cn"
}
```
#### resp exampl: "Mint success"
### event:
#### topic: Mint
#### data: ZeroAddr, to, tokenId, categoryName, metadata
#### example:
```json
[
  "00000000000000000000",
  "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
  "111111111111111111111112",
  "huanle",
  "url:https://chainmaker.org.cn"
]
```

## MintBatch
### args:
#### key1: "tokens"
#### value1: json
#### example:
```json
{
  "tokens": [{
    "tokenId": "xxxxx",
    "categoryName": "111",
    "to": "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
    "metadata": "aaa"
  }, {
    "tokenId": "xxxxx1",
    "categoryName": "111",
    "to": "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
    "metadata": "aaa"
  }]
}
```
#### resp exampl: "MintBatch success"
### event:
#### topic: Mint
#### data: tokens
#### example:
```json
[
  "00000000000000000000",
  "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
  "xxxxx",
  "111",
  "url:https://chainmaker.org.cn"
]
```

## SetApproval
### args:
#### key1: "owner"
#### value1: string
#### key2: "to"
#### value2: string
#### key3: "tokenId"
#### value3: string
#### key4: "isApproval"
#### value4: string
#### example:
```json
{
  "to":"a04f7895de24f61807a729be230f03da8c0eef42", 
  "tokenId":"111111111111111111111112",
  "isApproval": "true"
}
```
#### resp exampl: "SetApproval success"
### event:
#### topic: SetApproval
#### data: owner, to, tokenId, isApproval
#### example:
```json
[
  "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
  "a04f7895de24f61807a729be230f03da8c0eef42",
  "111111111111111111111112",
  "true"
]
```

## SetApprovalForAll
### args:
#### key1: "owner"
#### value1: string
#### key2: "to"
#### value2: string
#### key3: "isApproval"
#### value3: string
#### example:
```json
{
  "to":"a04f7895de24f61807a729be230f03da8c0eef42", 
  "isApproval": "true"
}
```
#### resp exampl: "SetApprovalForAll success"
### event:
#### topic: SetApprovalForAll
#### data: owner, to, isApproval
#### example:
```json
[
  "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
  "a04f7895de24f61807a729be230f03da8c0eef42",
  "true"
]
```

## TransferFrom
### args:
#### key1: "from"
#### value1: string
#### key2: "to"
#### value2: string
#### key3: "tokenId"
#### value3: string
#### example:
```json
{
  "from":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
  "to":"a04f7895de24f61807a729be230f03da8c0eef42", 
  "tokenId": "111111111111111111111112"
}
```
#### resp exampl: "TransferFrom success"
### event:
#### topic: TransferFrom
#### data: from, to, tokenId
#### example:
```json
[
  "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
  "a04f7895de24f61807a729be230f03da8c0eef42",
  "111111111111111111111112"
]
```

## TransferFromBatch
### args:
#### key1: "from"
#### value1: string
#### key2: "to"
#### value2: string
#### key3: "tokenIds"
#### value3: array of string
#### example:
```json
{
  "from":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
  "to":"a04f7895de24f61807a729be230f03da8c0eef42", 
  "tokenIds": [
    "111111111111111111111111",
    "111111111111111111111112"
  ]
}
```
#### resp exampl: "TransferFromBatch success"
### event:
#### topic: TransferFrom
#### data: from, to, tokenId
#### example:
```json
[
  "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
  "a04f7895de24f61807a729be230f03da8c0eef42",
  "111111111111111111111112"
]
```

## OwnerOf
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{"tokenId":"111111111111111111111112"}
```
### response example: "ec47ae0f0d6a0e952c240383d70ab43b19997a9f"

## TokenURI
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{"tokenId":"111111111111111111111112"}
```
#### resp exampl: "http://chainmaker.org.cn/111111111111111111111112"

## SetApprovalByCategory
### args:
#### key1: "owner"
#### value1: string
#### key2: "to"
#### value2: string
#### key3: "categoryName"
#### value3: string
#### key4: "isApproval"
#### value4: string
#### example:
```json
{
  "to":"a04f7895de24f61807a729be230f03da8c0eef42", 
  "categoryName": "1111",
  "isApproval": "true"
}
```
#### resp exampl: "SetApprovalByCategory success"
### event:
#### topic: SetApprovalByCategory
#### data: owner, to, categoryName, isApproval
#### example:
```json
[
  "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
  "a04f7895de24f61807a729be230f03da8c0eef42",
  "1111",
  "true"
]
```

## CreateOrSetCategory
### args:
#### key1: "category"
#### value1: string
#### example:
```json
{
  "category": {
    "categoryName":"1111",
    "categoryURI":"http://www.chainmaker.org.cn/"
  }
}
```
#### resp exampl: "CreateOrSetCategory success"
### event:
#### topic: CreateOrSetCategory
#### data: categoryName, categoryURI
#### example:
```json
[
  "1111",
  "http://www.chainmaker.org.cn/"
]
```

## Burn
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{
  "tokenId": "111111111111111111111112"
}
```
#### resp exampl: "Burn success"
### event:
#### topic: Burn
#### data: from, to, tokenId
#### example:
```json
[
  "111111111111111111111112"
]
```

## GetCategoryByName
### args:
#### key1: "categoryName"
#### value1: string
#### example:
```json
{
  "categoryName": "1111"
}
```
#### resp exampl:
```json
{
"categoryName":"1111",
"categoryURI":"http://www.chainmaker.org.cn/"
}
```

## GetCategoryByTokenId
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{
  "tokenId": "111111111111111111111112"
}
```
#### resp exampl:
```json
{
"categoryName":"1111",
"categoryURI":"http://www.chainmaker.org.cn/"
}
```

## TotalSupply
### args:
#### resp exampl: "0"

## TotalSupplyOfCategory
### args:
#### key1: "categoryName"
#### value1: string
#### example:
```json
{
  "categoryName": "111111111111111111111112"
}
```
#### resp exampl: "0"

## BalanceOf
### args:
#### key1: "account"
#### value1: string
#### example:
```json
{"account":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f"}
```
### response example: "0"

## AccountTokens
### args:
#### key1: "account"
#### value1: string
#### example:
```json
{"account":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f"}
```
#### resp exampl:
```json
{"account":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f","tokens":["111111111111111111111112","111111111111111111111113"]}
```

## TokenMetadata
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{"tokenId":"111111111111111111111112"}
```
#### resp exampl: "url:http://chainmaker.org.cn/111111111111111111111112"

## Test

### 部署合约
```sh
./cmc client contract user create --contract-name=CMNFA --version=1.0 --sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml --byte-code-path=./testdata/CMNFA.7z --runtime-type=DOCKER_GO \
--admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key \
--params="{\"categoryName\":\"1111\", \"categoryURI\":\"chainmaker.org.cn/\", \"admin\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"
```

### 升级合约
```sh
./cmc client contract user upgrade --contract-name=CMNFA --version=2.0 --sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml --byte-code-path=./testdata/CMNFA.7z --runtime-type=DOCKER_GO \
--admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key \
--params="{\"categoryName\":\"1111\", \"categoryURI\":\"chainmaker.org.cn/\", \"admin\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"
```

### 发行NFA
#### 验证Case1：
发行后可以使用BalanceOf、OwnerOf以及AccountTokens进行验证是否正确发行
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=Mint --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"to\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"tokenId\":\"111111111111111111111112\", \"categoryName\":\"1111\", \"metadata\":\"aaa\"}"
```

### 批量发行NFA
#### 验证Case1：
发行后可以使用BalanceOf、OwnerOf以及AccountTokens进行验证是否正确发行
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=MintBatch --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{
  \"tokens\": [{
    \"tokenId\": \"xxxxx\",
    \"categoryName\": \"1111\",
    \"to\": \"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\",
    \"metadata\": \"YWFh\"
  }, {
    \"tokenId\": \"xxxxx1\",
    \"categoryName\": \"1111\",
    \"to\": \"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\",
    \"metadata\": \"YWFh\"
  }]
}"
```

### 授权
#### 验证Case1：
授权后需要验证授权信息是否正确
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=SetApproval --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"tokenId\":\"111111111111111111111112\", \"isApproval\":\"true\"}"
```

### 全部授权
#### 验证Case1：
授权后需要验证授权信息是否正确
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=SetApprovalForAll --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"isApproval\":\"true\"}"
```

### 转账
#### 验证Case1：
转账后需要验证Owner是否发生了变化
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=TransferFrom --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"from\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"tokenId\":\"111111111111111111111112\"}"
```

### 批量转账
#### 验证Case1：
转账后需要Owner是否发生了变化
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=TransferFromBatch --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{
  \"from\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\",
  \"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", 
  \"tokenIds\": [
    \"111111111111111111111111\",
    \"111111111111111111111112\"
  ]
}"
```

### 查询OwnerOf
#### 验证Case1：
验证返回的owner是否为安装合约时指定的tokenURI+'/'+tokenId
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=OwnerOf --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"tokenId\":\"111111111111111111111112\"}"
```

### 查询tokenURI
#### 验证Case1：
验证返回的tokenURI是否为安装合约时指定的tokenURI+'/'+tokenId
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=TokenURI --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"tokenId\":\"111111111111111111111112\"}"
```

### 按分类授权
#### 验证Case1：
授权后需要验证授权信息是否正确
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=SetApprovalByCategory --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"categoryName\":\"1111\", \"isApproval\":\"true\"}"
```

### 创建或设置分类信息
#### 验证Case1：
创建后需要验证分类信息是否正确
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=CreateOrSetCategory --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{
  \"category\": {
    \"categoryName\":\"1111\",
    \"categoryURI\":\"http://www.chainmaker.org.cn/\"
  }
}"
```

### 销毁NFA
#### 验证Case1：
销毁后需要验证token是否还存在Owner
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=Burn --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"tokenId\":\"111111111111111111111112\"}"
```

### 根据名称查询分类信息
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=GetCategoryByName --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"categoryName\":\"1111\"}"
```

### 根据TokenId查询所属分类信息
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=GetCategoryByTokenId --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"tokenId\":\"111111111111111111111112\"}"
```

### 查询已发行NFA总数量
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=TotalSupply --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml
```

### 查询某个分类下已发行NFA总数量
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=TotalSupplyOfCategory --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"categoryName\":\"1111\"}"
```

### 查询账户NFA数量
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=BalanceOf --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"account\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"
```

### 查询账户下的所有NFA
#### 验证Case1：
验证账户下是否包含了所有发行的NFA
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=AccountTokens --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
 --params="{\"account\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"
```

### 查询token metadata信息
#### 验证Case1：
这儿验证查询到的metadata是否和mint时传递的一致
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=TokenMetadata --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"tokenId\":\"111111111111111111111112\"}"
```
