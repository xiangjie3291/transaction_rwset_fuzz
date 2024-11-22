# identity合约
身份认证合约，提供设置权限，查询权限、公钥功能。可在实际合约中，使用跨合约调用identity合约检查该使用者的权限。


## 主要合约接口

参考： 长安链CMID（CM-CS-221221-Identity）存证合约标准实现：
https://git.chainmaker.org.cn/contracts/standard/-/blob/master/living/CM-CS-221221-Identity.md

```go
// Identities 获取该合约支持的所有认证类型
// @return metas, 所有的认证类型编号和认证类型描述
Identities() (metas []IdentityMeta)

// SetIdentity 为地址设置认证类型，管理员可调用
// @param address 必填，公钥/证书的地址。一个地址仅能绑定一个公钥和认证类型编号，重复输入则覆盖。
// @param pkPem 选填,pem格式公钥，可用于验签
// @param level 必填,认证类型编号
// @param metadata 选填,其他信息，json格式字符串，比如：地址类型，上链人身份、组织信息，上链可信时间，上链批次等等
// @return error 返回错误信息
// @event topic: setIdentity(address, level, pkPem)
SetIdentity(address, pkPem string, level int, metadata string) error

// IdentityOf 获取认证信息
// @param address 地址
// @return int 返回当前认证类型编号
// @return identity 认证信息
// @return err 返回错误信息
IdentityOf(address string) (identity IdentityMeta, err error)

// LevelOf 获取认证编号
// @param address 地址
// @return level 返回当前认证类型编号
// @return err 返回错误信息
LevelOf(address string) (level int, err error)

// EmitSetIdentityEvent 发送设置认证类型事件
// @param address 地址
// @param pkPem pem格式公钥
// @param level 认证类型编号
EmitSetIdentityEvent(address, pkPem string, level int)

// PkPemOf 获取公钥
// @param address 地址
// @return string 返回当前地址绑定的公钥
// @return error 返回错误信息
PkPemOf(address string) (string, error)

// SetIdentityBatch 设置多个认证类型，管理员可调用
// @param identities, 入参json格式字符串
// @event topic: setIdentity(address, level, pkPem)
SetIdentityBatch(identities []Identity) error

// AlterAdminAddress 修改管理员，管理员可调用
// @param adminAddresses 管理员地址，可为空，默认为创建人地址。入参为以逗号分隔的地址字符串"addr1,addr2"
// @return error 返回错误信息
// @event topic: alterAdminAddress（adminAddresses）
AlterAdminAddress(adminAddresses string) error
```

## 跨合约使用

```go
args := make(map[string][]byte)
args["address"] = []byte("55e5a6d2f6a66b6c3e6a6f55e5a6d2f6a66b6c3e")
resp := sdk.Instance.CallContract("identity", "levelOf", args)
if string(resp.Payload) != "4" {
    return Error("address not registered")
}
```

## cmc使用示例

命令行工具使用示例
```sh
echo
echo "安装合约 identity，创建者为管理员"
./cmc client contract user create \
--contract-name=identity \
--runtime-type=DOCKER_GO \
--byte-code-path=./identity.7z \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--params="{}"

echo
echo "执行合约 identity.SetIdentity(address, level, pkPem string)，设置身份"
./cmc client contract user invoke \
--contract-name=identity \
--method=SetIdentity \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"address\":\"5fa92a33364dd5ce26a9814a6aceb240bd6bf083\",\"level\":\"4\",\"pkPem\":\"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECr/ijK264TXbHfIvhJJz43z9hZLroyWUZY371pfQqaToo2by5ljMj3Ot8/XM2n5Xr/6xIwVLJ7t+C5cwcGRjqA==\n-----END PUBLIC KEY-----\",\"metadata\":\"something\"}" \
--sync-result=true \
--result-to-string=true

echo
echo "执行合约 identity.SetIdentityBatch(identities []standard.Identity)，设置身份"
./cmc client contract user invoke \
--contract-name=identity \
--method=SetIdentityBatch \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"identities\":\"[{\\\"address\\\":\\\"5fa92a33364dd5ce26a9814a6aceb240bd6bf082\\\",\\\"level\\\":4,\\\"pkPem\\\":\\\"-----BEGIN PUBLIC KEY-----\\\\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECr/ijK264TXbHfIvhJJz43z9hZLroyWUZY371pfQqaToo2by5ljMj3Ot8/XM2n5Xr/6xIwVLJ7t+C5cwcGRjqA==\\\\n-----END PUBLIC KEY-----\\\",\\\"metadata\\\":\\\"something\\\"},{\\\"address\\\":\\\"5fa92a33364dd5ce26a9814a6aceb240bd6bf081\\\",\\\"level\\\":3,\\\"pkPem\\\":\\\"-----BEGIN PUBLIC KEY-----\\\\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECr/ijK264TXbHfIvhJJz43z9hZLroyWUZY371pfQqaToo2by5ljMj3Ot8/XM2n5Xr/6xIwVLJ7t+C5cwcGRjqA==\\\\n-----END PUBLIC KEY-----\\\",\\\"metadata\\\":\\\"something\\\"}]\"}" \
--sync-result=true \
--result-to-string=true

echo
echo "执行合约 identity.AlterAdminAddress(adminAddress string)，设置管理员"
./cmc client contract user invoke \
--contract-name=identity \
--method=AlterAdminAddress \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"adminAddress\":\"e21b6661e9fa7d0056b6c3b0f97995f95ac4d540\"}" \
--sync-result=true \
--result-to-string=true

echo
echo "执行合约 identity.LevelOf(address string)，获取身份类型编号"
./cmc client contract user get \
--contract-name=identity \
--method=LevelOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"address\":\"5fa92a33364dd5ce26a9814a6aceb240bd6bf083\"}" \
--result-to-string=true

echo
echo "执行合约 identity.PkPemOf(address string)，获取pem公钥串"
./cmc client contract user get \
--contract-name=identity \
--method=PkPemOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"address\":\"5fa92a33364dd5ce26a9814a6aceb240bd6bf083\"}" \
--result-to-string=true

echo
echo "执行合约 identity.IdentityOf(address string)，获取身份信息"
./cmc client contract user get \
--contract-name=identity \
--method=IdentityOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"address\":\"5fa92a33364dd5ce26a9814a6aceb240bd6bf083\"}" \
--result-to-string=true

echo
echo "查询合约 identity.SupportStandard(standardName string)，是否支持合约标准CMID"
./cmc client contract user get \
--contract-name=identity \
--method=SupportStandard \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"standardName\":\"CMID\"}" \
--result-to-string=true

echo
echo "查询合约 identity.SupportStandard(standardName string)，是否支持合约标准CMBC"
./cmc client contract user get \
--contract-name=identity \
--method=SupportStandard \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"standardName\":\"CMBC\"}" \
--result-to-string=true

echo
echo "查询合约 identity.Standards()，支持的合约标准列表"
./cmc client contract user get \
--contract-name=identity \
--method=Standards \
--sdk-conf-path=./testdata/sdk_config.yml \
--result-to-string=true

echo
echo "查询合约 identity.Identities()，当前合约支持的身份类型列表"
./cmc client contract user get \
--contract-name=identity \
--method=Identities \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"standardName\":\"CMBC\"}" \
--result-to-string=true
```
