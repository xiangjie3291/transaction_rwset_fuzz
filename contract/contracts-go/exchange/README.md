# exchange合约
交易所合约，用于转移ERC721(NFT) token的时候同时转移ERC20的token
## 主要合约接口

```go
// 将nft tokenId从 from 转移到to
// 并将erc20 amount数量的token从to转移到from
buyNow(tokenId, from, to, amount string) protogo.Response // return "true","false"
```
## cmc使用示例

命令行工具使用示例

```sh
echo
echo "安装合约 exchange"
./cmc client contract user create \
--contract-name=exchange \
--version=1.0 \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--byte-code-path=./testdata/exchange.7z \
--runtime-type=DOCKER_GO \
--admin-crt-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org.chainmaker.org/user/admin1/admin1.sign.key \
--params=""


echo
echo "执行合约 exchange，购买nft"
./cmc client contract user invoke \
--contract-name=exchange \
--method=buyNow \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"from\":\"2a230c7ea2110446e6320a44091089f111cb5028\",\"to\":\"335338acd9f5a757a12cf4a35bf327d6051e7068\",\"tokenId\":\"1\", \"amount\":\"10\"}" \
--sync-result=true \
--result-to-string=true
```