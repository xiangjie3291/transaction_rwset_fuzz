ContractStandardNameCMDFA ChainMaker - Contract Standard - Digital Fungible Assets

https://git.chainmaker.org.cn/contracts/standard/-/blob/master/living/CM-CS-221221-DFA.md

## CMC示例
### 安装合约
Token名：CM DFA
符号名：DFA
小数位数：8
```sh
./cmc client contract user create --contract-name=cmdfa --version=1.0 --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --byte-code-path=./testdata/cmdfa.7z --runtime-type=DOCKER_GO --admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt --admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key --params="{\"name\":\"CM DFA\", \"symbol\":\"DFA\", \"decimals\":\"8\"}"
```
### 铸造Token
铸造10个Token，因为小数位数为8，所以参数的数量是：1000000000
```sh
./cmc client contract user invoke --contract-name=cmdfa --method=Mint --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"account\":\"自己的地址\", \"amount\":\"1000000000\"}"
```
### 查询余额
```sh
./cmc client contract user get --contract-name=cmdfa --method=BalanceOf --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"account\":\"自己的地址\"}"
```
### 转账1个Token
转移1个Token，因为小数位数为8，所以参数的数量是：100000000
```sh
./cmc client contract user invoke --contract-name=cmdfa --method=Transfer --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"to\":\"UserB的地址\", \"amount\":\"100000000\"}"
```

### 销毁1个Token
```sh
./cmc client contract user invoke --contract-name=cmdfa --method=Burn --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"amount\":\"100000000\"}"
```

### 查询Token的名字
```sh
./cmc client contract user get --contract-name=cmdfa --method=Name --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{}"
```
### 查询Token的符号
```sh
./cmc client contract user get --contract-name=cmdfa --method=Symbol --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{}"
```

### 查询Token的小数位数
```sh
./cmc client contract user get --contract-name=cmdfa --method=Decimals --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{}"
```

### 查询Token的发行量
```sh
./cmc client contract user get --contract-name=cmdfa --method=TotalSupply --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{}"
```

## 授权转移操作
### 授权3个Token给用户C
```sh
./cmc client contract user invoke --contract-name=cmdfa --method=Approve --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"spender\":\"UserC的地址\", \"amount\":\"300000000\"}"
```
### C用户转移Token给用户B
因为是用的UserC的身份，所以需要切换SDK_Config到UserC的对应配置
```sh
./cmc client contract user invoke --contract-name=cmdfa --method=TransferFrom --sync-result=true --sdk-conf-path=./testdata/sdk_config_UserC.yml --params="{\"from\":\"UserA的地址\",\"to\":\"UserB的地址\", \"amount\":\"100000000\"}"
```
### 查询用户A授权给用户C的转移额度
本来授权了3Token，现在转移了1个，所以结果应该2Token
```sh
./cmc client contract user get --contract-name=cmdfa --method=Allowance --sdk-conf-path=./testdata/sdk_config_UserC.yml --params="{\"spender\":\"UserC的地址\",\"owner\":\"UserA的地址\"}"
```
### UserC从授权账户UserA销毁1Token
因为是用的UserC的身份，所以需要切换SDK_Config到UserC的对应配置
```sh
./cmc client contract user invoke --contract-name=cmdfa --method=BurnFrom --sync-result=true --sdk-conf-path=./testdata/sdk_config_UserC.yml --params="{\"account\":\"UserA的地址\", \"amount\":\"100000000\"}"
```
