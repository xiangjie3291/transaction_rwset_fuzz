

[toc]

# 编译

```sh
make
```



# 合约调用cmc命令

to construct chainmaker with docker vm please see the document below

## 测试ERC20合约命令

部署合约
```sh
./cmc client contract user create --contract-name=erc20 --version=1.0 --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --byte-code-path=./testdata/erc20.7z --runtime-type=DOCKER_GO --admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt --admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key --params="{\"name\":\"huanletoken\", \"totalSupply\":\"10000000000000\", \"symbol\":\"hlt\", \"decimals\":\"17\"}"
```
查询name
```sh
./cmc client contract user invoke --contract-name=erc20 --method=name --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml
```
查询symbol
```sh
./cmc client contract user invoke --contract-name=erc20 --method=symbol --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml
```
查询decimals
```sh
./cmc client contract user invoke --contract-name=erc20 --method=decimals --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml
```
查询totalSupply
```sh
./cmc client contract user invoke --contract-name=erc20 --method=totalSupply --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml
```
查询账户余额
```sh
./cmc client contract user invoke --contract-name=erc20 --method=balanceOf --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"account\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"
```
增发
```sh
./cmc client contract user invoke --contract-name=erc20 --method=mint --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"account\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"amount\":\"20000\"}"
```
转账
```sh
./cmc client contract user invoke --contract-name=erc20 --method=transfer --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"to\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"amount\":\"20000\"}"
```
授权
```sh
./cmc client contract user invoke --contract-name=erc20 --method=approve --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"spender\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"amount\":\"20000\"}"
```
查询授权
```sh
./cmc client contract user invoke --contract-name=erc20 --method=allowance --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"owner\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"spender\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"
```
根据授权转账
```sh
./cmc client contract user invoke --contract-name=erc20 --method=transferFrom --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"owner\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"to\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"amount\":\"20000\"}"
```
## 测试ERC721合约命令
部署合约

```sh
./cmc client contract user create --contract-name=erc721 --version=1.0 --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --byte-code-path=./testdata/erc721.7z --runtime-type=DOCKER_GO --admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt --admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key --params="{\"name\":\"huanletoken\", \"symbol\":\"hlt\"}"
```
查询name
```sh
./cmc client contract user invoke --contract-name=erc721 --method=name --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml
```
查询symbol
```sh
./cmc client contract user invoke --contract-name=erc721 --method=symbol --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml
```
查询账户nft数量
```sh
./cmc client contract user invoke --contract-name=erc721 --method=balanceOf --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"account\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"
```
发行nft
```sh
./cmc client contract user invoke --contract-name=erc721 --method=mint --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"to\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"tokenId\":\"111111111111111111111112\"}"
```
查询nft所属账户
```sh
./cmc client contract user invoke --contract-name=erc721 --method=ownerOf --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"tokenId\":\"11111111111111111111111\"}"
```
授权
```sh
./cmc client contract user invoke --contract-name=erc721 --method=approve --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"to\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"tokenId\":\"11111111111111111111111\"}"
./cmc client contract user invoke --contract-name=erc721 --method=approve --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"tokenId\":\"11111111111111111111111\"}"
```
获取授权信息
```sh
./cmc client contract user invoke --contract-name=erc721 --method=getApprove --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"tokenId\":\"11111111111111111111111\"}"
```
根据授权转账
```sh
./cmc client contract user invoke --contract-name=erc721 --method=transferFrom --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"from\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"to\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"tokenId\":\"11111111111111111111111\"}"
```
根据授权转账
```sh
./cmc client contract user invoke --contract-name=erc721 --method=transferFrom --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"from\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"tokenId\":\"11111111111111111111111\"}"
```
发行nft
```sh
./cmc client contract user invoke --contract-name=erc721 --method=minted --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"tokenId\":\"111111111111111111111112\"}"
```
## 测试exchange合约命令

安装合约

```sh
./cmc client contract user create --contract-name=exchange --version=1.0 --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --byte-code-path=./testdata/exchange.7z --runtime-type=DOCKER_GO --admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt --admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key --params=""
```



购买nft

```sh
./cmc client contract user invoke --contract-name=exchange --method=buyNow --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"from\":\"2a230c7ea2110446e6320a44091089f111cb5028\",\"to\":\"335338acd9f5a757a12cf4a35bf327d6051e7068\",\"tokenId\":\"1\", \"amount\":\"1\"}" --sync-result=true
```
## 测试identity合约命令

安装合约

```sh
./cmc client contract user create --contract-name=exchange --version=1.0 --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --byte-code-path=./testdata/identity.7z --runtime-type=DOCKER_GO --admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt --admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key --params=""
```

调用合约 查询地址
```sh
# org1 admin1  005521f27d745a04999c6d09f559764f9c44376a
./cmc client contract user invoke \
--contract-name=identity \
--method=callerAddress \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="" \
--sync-result=true
```
添加白名单
```sh
./cmc client contract user invoke \
--contract-name=identity \
--method=addWriteList \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"address\":\"005521f27d745a04999c6d09f559764f9c44376a,335338acd9f5a757a12cf4a35bf327d6051e7068,2a230c7ea2110446e6320a44091089f111cb5028\"}" \
--sync-result=true
```
查询是否在白名单
```sh
./cmc client contract user invoke \
--contract-name=identity \
--method=isApprovedUser \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"address\":\"005521f27d745a04999c6d09f559764f9c44376a\"}" \
--sync-result=true
```
移除白名单
```sh
./cmc client contract user invoke \
--contract-name=identity \
--method=removeWriteList \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"address\":\"005521f27d745a04999c6d09f559764f9c44376a\"}" \
--sync-result=true
```

## 测试fact合约命令

安装合约

```sh
./cmc client contract user create --contract-name=fact --version=1.0 --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --byte-code-path=./testdata/fact.7z --runtime-type=DOCKER_GO --admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt --admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key --params=""
```

存入数据

```sh
./cmc client contract user invoke \
--contract-name=fact \
--method=save \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"file_hash\":\"005521f27d745a04999c6d09f559764f9c44376a\",\"file_name\":\"aoteman.jpg\",\"time\":\"16456254\"}" \
--sync-result=true
```

查询数据

```sh
./cmc client contract user get \
--contract-name=fact \
--method=findByFileHash \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"file_hash\":\"005521f27d745a04999c6d09f559764f9c44376a\"}" \
--sync-result=true
```

## 测试投票合约

安装合约

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

发布投票项目
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=issueProject \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"projectInfo\":{\"Id\":\"projectId2\",\"PicUrl\":\"www.sina.com\",\"Title\":\"wonderful\",\"StartTime\":\"1664450291\",\"EndTime\":\"1665314291\",\"Desc\":\"the 1\",\"Items\":[{\"Id\":\"item1\",\"PicUrl\":\"www.baidu.com\",\"Desc\":\"beautiful\",\"Url\":\"www.qq.com\"},{\"Id\":\"item2\",\"PicUrl\":\"www.baidu.com\",\"Desc\":\"beautiful\",\"Url\":\"www.qq.com\"}]}}"
```

投票
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=vote \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"projectId\":\"projectId1\",\"itemId\":\"item1\"}"
```
查询项目投票人员
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=queryProjectVoters \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"projectId\":\"projectId1\"}"
```

查询投票项投票人员
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=queryProjectItemVoters \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"projectId\":\"projectId1\",\"itemId\":\"item1\"}"
```


## 其他命令

查询区块

```sh
./cmc query block-by-height 4 --sdk-conf-path=./testdata/sdk_config_solo.yml
```
生成账户地址
```sh
./cmc address cert-to-addr ./testdata/crypto-config/wx-org.chainmaker.org/user/client1/client1.tls.crt
```
