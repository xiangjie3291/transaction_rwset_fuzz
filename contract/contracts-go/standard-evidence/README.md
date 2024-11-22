# evidence合约

存证合约，提供存证、查验功能。

## 主要合约接口

参考： 长安链CMEVI（CM-CS-221221-Evidence）存证合约标准实现：
https://git.chainmaker.org.cn/contracts/standard/-/blob/master/living/CM-CS-221221-Evidence.md

```go

// CMEVI 长安链存证合约go接口
// https://git.chainmaker.org.cn/contracts/wx-organdard/-/blob/master/draft/CM-CS-221221-Evidence.md


// Evidence 存证
// @param id 必填，流水号
// @param hash 必填，上链哈希值
// @param metadata 可选，其他信息；比如：哈希的类型（文字，文件）、文字描述的json格式字符串，具体参考下方 Metadata 对象。
// @return error 返回错误信息
Evidence(id string, hash string, metadata string) error

// ExistsOfHash 哈希是否存在
// @param hash 必填，上链的哈希值
// @return exist 存在：true，"true"；不存在：false，"false"
// @return err 错误信息
ExistsOfHash(hash string) (exist bool, er error)

// ExistsOfId 哈希是否存在
// @param id 必填，上链的ID值
// @return exist 存在：true，"true"；不存在：false，"false"
// @return err 错误信息
ExistsOfId(id string) (exist bool, er error)

// FindByHash 根据哈希查找
// @param hash 必填，上链哈希值
// @return evidence 上链时传入的evidence信息
// @return err 返回错误信息
FindByHash(hash string) (evidence *Evidence, err error)

// FindById 根据id查找
// @param id 必填，流水号
// @return evidence 上链时传入的evidence信息
// @return err 返回错误信息
FindById(id string) (evidence *Evidence, err error)

// EvidenceBatch 批量存证
// @param evidences 必填，存证信息
// @return error 返回错误信息
EvidenceBatch(evidences []Evidence) error

// UpdateEvidence 根据ID更新存证哈希和metadata
// @param id 必填，已经在链上存证的流水号。 如果是新流水号返回错误信息不存在
// @param hash 必填，上链哈希值。必须与链上已经存储的hash不同
// @param metadata 可选，其他信息；具体参考下方 Metadata 对象。
// @return error 返回错误信息
// @desc 该方法由长安链社区志愿者@sunhuiyuan提供建议，感谢支持
UpdateEvidence(id string, hash string, metadata string) error
```

## cmc使用示例

命令行工具使用示例

```sh
echo ""
echo "create DOCKER_GO evidence"
./cmc client contract user create \
--contract-name=evidence \
--runtime-type=DOCKER_GO \
--byte-code-path=./evidence.7z \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--params="{}"

echo ""
echo "upgrade DOCKER_GO evidence"
./cmc client contract user upgrade \
--contract-name=evidence \
--runtime-type=DOCKER_GO \
--byte-code-path=./evidence.7z \
--version=2.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"

echo ""
echo "invoke evidence.Evidence(id, hash, metadata string)  hash:a2 id:02"
./cmc client contract user invoke \
--contract-name=evidence \
--method=Evidence \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"id\":\"02\",\"hash\":\"a2\",\"metadata\":\"{\\\"hashType\\\":\\\"file\\\",\\\"hashAlgorithm\\\":\\\"sha256\\\",\\\"username\\\":\\\"taifu\\\",\\\"timestamp\\\":\\\"1672048892\\\",\\\"proveTimestamp\\\":\\\"\\\"}\"}" \
--sync-result=true \
--result-to-string=true

echo ""
echo "invoke evidence.EvidenceBatch(evidences []Evidence)"
./cmc client contract user invoke \
--contract-name=evidence \
--method=EvidenceBatch \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"evidences\":\"[{\\\"id\\\":\\\"id1\\\",\\\"hash\\\":\\\"hash1\\\",\\\"txId\\\":\\\"\\\",\\\"blockHeight\\\":0,\\\"timestamp\\\":\\\"\\\",\\\"metadata\\\":\\\"11\\\"},{\\\"id\\\":\\\"id2\\\",\\\"hash\\\":\\\"hash2\\\",\\\"txId\\\":\\\"\\\",\\\"blockHeight\\\":0,\\\"timestamp\\\":\\\"\\\",\\\"metadata\\\":\\\"11\\\"}]\"}" \
--sync-result=true \
--result-to-string=true

echo ""
echo "query evidence.ExistsOfHash(hash string) hash:a2"
./cmc client contract user get \
--contract-name=evidence \
--method=ExistsOfHash \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"hash\":\"a2\"}" \
--result-to-string=true

echo ""
echo "query evidence.ExistsOfId(id string) id:02"
./cmc client contract user get \
--contract-name=evidence \
--method=ExistsOfId \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"id\":\"02\"}" \
--result-to-string=true

echo ""
echo "query evidence.FindByHash(hash string) hash:a2"
./cmc client contract user get \
--contract-name=evidence \
--method=FindByHash \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"hash\":\"a2\"}" \
--result-to-string=true

echo ""
echo "query evidence.FindById(id string) id:02"
./cmc client contract user get \
--contract-name=evidence \
--method=FindById \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"id\":\"02\"}" \
--result-to-string=true

echo
echo "查询合约 evidence.SupportStandard，是否支持合约标准CMEVI"
./cmc client contract user get \
--contract-name=evidence \
--method=SupportStandard \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"standardName\":\"CMEVI\"}" \
--result-to-string=true

echo
echo "查询合约 evidence.SupportStandard，是否支持合约标准CMBC"
./cmc client contract user get \
--contract-name=evidence \
--method=SupportStandard \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"standardName\":\"CMBC\"}" \
--result-to-string=true

echo
echo "查询合约 evidence.Standards，支持的合约标准列表"
./cmc client contract user get \
--contract-name=evidence \
--method=Standards \
--sdk-conf-path=./testdata/sdk_config.yml \
--result-to-string=true
```

