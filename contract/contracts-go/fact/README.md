## Test
### 安装合约
```sh
./cmc client contract user create --contract-name=fact --version=1.0 --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --byte-code-path=./testdata/fact.7z --runtime-type=DOCKER_GO --admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt --admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key --params=""
```

### 查询数据
#### 测试Case1:
数据不存在时，查到数据应为空
```sh
./cmc client contract user get \
--contract-name=fact \
--method=findByFileHash \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"file_hash\":\"005521f27d745a04999c6d09f559764f9c44376a\"}" \
--sync-result=true
```

### 存入存证数据
#### 测试Case1:
存入存证数据后需要验证数据是否能够查询得到
```sh
./cmc client contract user invoke \
--contract-name=fact \
--method=save \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"file_hash\":\"005521f27d745a04999c6d09f559764f9c44376a\",\"file_name\":\"aoteman.jpg\",\"time\":\"16456254\"}" \
--sync-result=true
```
