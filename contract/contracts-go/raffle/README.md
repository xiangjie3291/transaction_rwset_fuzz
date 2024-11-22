This is the raffle contract for test net demo.

## Test

### 安装合约
```sh
./cmc client contract user create \
--contract-name=raffle \
--version=1.0 \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--byte-code-path=./testdata/raffle.7z \
--runtime-type=DOCKER_GO \
--admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key --params="{}"
```

### 注册奖池名单
```sh
./cmc client contract user invoke \
--contract-name=raffle \
--method=registerAll \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"peoples\":{\"peoples\":[{\"num\":1,\"name\":\"Chris\"},{\"num\":2,\"name\":\"Linus\"}]}}"
```

### 抽奖
### 测试Case1: 
参数level或timestamp没有时应报错
### 测试Case2:
参数level和timestamp正确时应返回从奖池中抽出中奖的人员
### 测试Case3:
多次抽奖，中奖人员不应重复
### 测试Case4:
抽奖后查询奖池人员，应不包含已中奖人员
```sh
./cmc client contract user invoke \
--contract-name=raffle \
--method=raffle \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"level\":\"1\",\"timestamp\":\"13235432\"}"
```

### 查询奖池名单
### 测试Case1:
应返回正确的奖池人员名单
```sh
./cmc client contract user invoke \
--contract-name=raffle \
--method=query \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
```
