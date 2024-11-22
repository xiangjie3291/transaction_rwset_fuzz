package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"chainmaker.org/chainmaker/contract-utils/safemath"
)

const (
	pName        = "name"
	pSymbol      = "symbol"
	pDecimals    = "decimals"
	pAccount     = "account"
	pAmount      = "amount"
	pTotalSupply = "totalSupply"
	pFrom        = "from"
	pTo          = "to"
	pSpender     = "spender"
	pOwner       = "owner"
)

func main() {
	err := sandbox.Start(new(CmdfaContract))
	if err != nil {
		sdk.Instance.Errorf(err.Error())
	}
}

// InitContract install contract func
func (c *CmdfaContract) InitContract() protogo.Response {
	err := c.updateErc20Info()
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract upgrade contract func
func (c *CmdfaContract) UpgradeContract() protogo.Response {
	err := c.updateErc20Info()
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte("Upgrade contract success"))
}

// UpgradeContract upgrade contract func
func (c *CmdfaContract) updateErc20Info() error {
	args := sdk.Instance.GetArgs()
	// name, symbol and decimal are optional
	name := args[pName]
	symbol := args[pSymbol]
	decimalsStr := args[pDecimals]
	totalSupply := args[pTotalSupply]

	if len(name) > 0 {
		defaultName = string(name)
	}
	err := sdk.Instance.PutState(nameKey, "", defaultName)
	if err != nil {
		return err
	}
	if len(symbol) > 0 {
		defaultSymbol = string(symbol)
	}
	err = sdk.Instance.PutState(symbolKey, "", defaultSymbol)
	if err != nil {
		return err
	}
	//decimals default to 18
	if len(decimalsStr) > 0 {
		num, err1 := strconv.Atoi(string(decimalsStr))
		if err1 != nil {
			return fmt.Errorf("param decimals err")
		}
		defaultDecimals = num
	}
	err = sdk.Instance.PutState(decimalKey, "", strconv.Itoa(defaultDecimals))
	if err != nil {
		return err
	}
	//set admin
	admin, err := sdk.Instance.Origin()
	if err != nil {
		return fmt.Errorf("get sender failed, err:%s", err)
	}
	err = sdk.Instance.PutState(adminKey, "", admin)
	if err != nil {
		return err
	}
	//set totalSupply
	if len(totalSupply) > 0 {
		totalSupplyNum, ok := safemath.ParseSafeUint256(string(totalSupply))
		if !ok {
			return errors.New("invalid totalSupply number")
		}
		defaultTotalSupply = totalSupplyNum
	}
	//Mint
	if defaultTotalSupply.ToString() != "0" {
		return c.baseMint(admin, defaultTotalSupply)
	}
	return nil
}

// InvokeContract the entry func of invoke contract func
func (c *CmdfaContract) InvokeContract(method string) protogo.Response { // nolint:gocyclo
	if len(method) == 0 {
		return sdk.Error("method of param should not be empty")
	}
	switch method {
	case "Standards":
		return ReturnJson(c.Standards())
	case "Name":
		return ReturnString(c.Name())
	case "Symbol":
		return ReturnString(c.Symbol())
	case "Decimals":
		return ReturnUint8(c.Decimals())
	case "TotalSupply":
		return ReturnUint256(c.TotalSupply())
	case "BalanceOf":
		account, err := RequireString(pAccount)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return ReturnUint256(c.BalanceOf(account))
	case "Transfer":
		to, err := RequireString(pTo)
		if err != nil {
			return sdk.Error(err.Error())
		}
		amount, err := RequireUint256(pAmount)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(c.Transfer(to, amount))
	case "TransferFrom":
		from, err := RequireString(pFrom)
		if err != nil {
			return sdk.Error(err.Error())
		}
		to, err := RequireString(pTo)
		if err != nil {
			return sdk.Error(err.Error())
		}
		amount, err := RequireUint256(pAmount)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(c.TransferFrom(from, to, amount))
	case "Approve":
		spender, err := RequireString(pSpender)
		if err != nil {
			return sdk.Error(err.Error())
		}
		amount, err := RequireUint256(pAmount)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(c.Approve(spender, amount))
	case "Allowance":
		spender, err := RequireString(pSpender)
		if err != nil {
			return sdk.Error(err.Error())
		}
		owner, err := RequireString(pOwner)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return ReturnUint256(c.Allowance(owner, spender))
	case "Mint":
		account, err := RequireString(pAccount)
		if err != nil {
			return sdk.Error(err.Error())
		}
		amount, err := RequireUint256(pAmount)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(c.Mint(account, amount))
	case "Burn":
		amount, err := RequireUint256(pAmount)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(c.Burn(amount))
	case "BurnFrom":
		account, err := RequireString(pAccount)
		if err != nil {
			return sdk.Error(err.Error())
		}
		amount, err := RequireUint256(pAmount)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(c.BurnFrom(account, amount))
	default:
		return sdk.Error("Invalid method")
	}
}

////////////////////////////////Helper//////////////////////////////////

// ReturnUint256 封装返回SafeUint256类型为Response，如果有error则忽略num，封装error
// @param num
// @param err
// @return Response
func ReturnUint256(num *safemath.SafeUint256, err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(num.ToString()))
}

// ReturnString 封装返回string类型为Response，如果有error则忽略str，封装error
// @param str
// @param err
// @return Response
func ReturnString(str string, err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(str))
}

// ReturnJson 封装返回interface类型为json string Response
// @param data
// @return Response
func ReturnJson(data interface{}) protogo.Response {
	standardsBytes, err := json.Marshal(data)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success(standardsBytes)
}

// Return 封装返回Bool类型为Response，如果有error则忽略bool，封装error
// @param err
// @return Response
func Return(err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.SuccessResponse
}

// ReturnUint8 封装返回uint8类型为Response，如果有error则忽略num，封装error
// @param num
// @param err
// @return Response
func ReturnUint8(num uint8, err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(strconv.Itoa(int(num))))
}

// RequireString 必须要有参数 string类型
// @param key
// @return string
// @return error
func RequireString(key string) (string, error) {
	args := sdk.Instance.GetArgs()
	b, ok := args[key]
	if !ok || len(b) == 0 {
		return "", fmt.Errorf("CMDFA: require parameter:'%s'", key)
	}
	return string(b), nil
}

// RequireUint256 必须要有参数 Uint256类型
// @param key
// @return *safemath.SafeUint256
// @return error
func RequireUint256(key string) (*safemath.SafeUint256, error) {
	args := sdk.Instance.GetArgs()
	b, ok := args[key]
	if !ok {
		return nil, fmt.Errorf("CMDFA: require parameter:'%s'", key)
	}
	num, ok := safemath.ParseSafeUint256(string(b))
	if !ok {
		return nil, fmt.Errorf("CMDFA: parameter:'%s' not a valid uint256", key)
	}
	return num, nil
}
