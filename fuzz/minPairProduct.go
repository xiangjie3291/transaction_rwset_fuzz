/*
	目前使用的方案为，如果当前测试函数所有输入都无法成功，将会给其中未confirm的变量置为default值

	当所有变量都被confirm后，将不会进行新的函数的尝试

	可能出现：变量A被设置为default，但包含变量A的其他函数并没有被执行，未被成功执行的函数下可能包含正确的类型

	该部分仅用于获取各参数可能的输入类型，存储于：

	GlobalParamAndCandidateTypes 下

	最终得到的结果为：

		a. 对每个变量，有自己的可用confirm value
*/

package fuzz

import (
	nodecontrol "TransactionRwset/nodeControl"
	"TransactionRwset/utils"
	"encoding/json"
	"fmt"
	"math"
	"reflect"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/google/go-cmp/cmp"
)

type Pair struct {
	KeyValuePair []*common.KeyValuePair
	KeyValue     map[string]interface{}
}

// 计算指定函数的匹配积
// 已经成功confirm的变量将不需要再进行计算
func calProduct(funcName string) {
	// 根据函数名获取输入参数
	paramNameList := utils.GlobalContractInfo.ContractFuncMap[funcName].ParamsNameList

	// 该函数不需要输入参数
	if len(paramNameList) == 0 {
		return
	}

	result := 1
	//根据输入参数计算积
	for _, paramName := range paramNameList {

		candidateTypes, ok := utils.GlobalContractInfo.ParamAndCandidateTypes[paramName]
		if !ok {
			fmt.Println("can't find this param!")
			return
		}

		if len(candidateTypes.Types) == 0 {
			fmt.Println(paramName, "this param don't have any type!")
			return
		}

		if !candidateTypes.Confirm {
			result *= len(candidateTypes.Types)
		}
	}
	utils.GlobalContractInfo.ContractFuncMap[funcName].Product = result
}

// 计算当前所有函数的匹配积，并取出最小的那个
// 当product值返回为-1时，代表所有输入的类型都已确定
func calMinProduct() (string, int) {
	for funcName, _ := range utils.GlobalContractInfo.ContractFuncMap {
		calProduct(funcName)
	}

	minProduct := math.MaxInt32
	var minKey string

	for key, val := range utils.GlobalContractInfo.ContractFuncMap {
		// 只考虑 Product 大于 1 的值
		if val.Product > 1 && val.Product < minProduct {
			minProduct = val.Product
			minKey = key
		}
	}

	// 如果没找到合适的 key，返回默认值
	if minKey == "" {
		return "", -1
	}

	return minKey, minProduct
}

// 输出所有可能性结果的组合
func generateCombinations(paramInputMap map[string][]interface{}) []map[string]interface{} {
	// 提取键以保持顺序
	keys := make([]string, 0, len(paramInputMap))
	for key := range paramInputMap {
		keys = append(keys, key)
	}
	// 初始化存储所有组合的切片
	var combinations []map[string]interface{}

	// 递归辅助函数
	var helper func(int, map[string]interface{})
	helper = func(index int, currentCombination map[string]interface{}) {
		// 基础情况：所有键都已处理
		if index == len(keys) {
			// 复制当前组合并存储
			combinationCopy := make(map[string]interface{}, len(currentCombination))
			for k, v := range currentCombination {
				combinationCopy[k] = v
			}
			combinations = append(combinations, combinationCopy)
			return
		}

		// 获取当前键及其可能的值
		key := keys[index]
		values := paramInputMap[key]

		// 遍历当前键的所有可能值
		for _, value := range values {
			// 将当前键设置为选定的值
			currentCombination[key] = value
			// 递归处理下一个键
			helper(index+1, currentCombination)
		}
	}

	// 从第一个键开始递归
	helper(0, make(map[string]interface{}))

	return combinations
}

// 根据函数，生成输入数据
func generateInput(funcName string) []Pair {
	// 获取输入参数
	paramNameList := utils.GlobalContractInfo.ContractFuncMap[funcName].ParamsNameList

	paramInputMap := make(map[string][]interface{}, len(paramNameList))

	// 对函数下的每个输入参数
	for _, paramName := range paramNameList {
		// paramName := paramNameList[i]
		candidateTypes := utils.GlobalContractInfo.ParamAndCandidateTypes[paramName]

		// 如果该param未被confirm
		if !candidateTypes.Confirm {
			paramInputMap[paramName] = make([]interface{}, 0)

			// 将每个备用类型都加入
			for _, value := range candidateTypes.Types {
				// fmt.Println(types, value)
				paramInputMap[paramName] = append(paramInputMap[paramName], value)
			}
		} else {
			// 如果param被confirm
			// 将被confirm的value第一个拿去使用即可
			paramInputMap[paramName] = make([]interface{}, 0)
			paramInputMap[paramName] = append(paramInputMap[paramName], candidateTypes.ConfirmValue[0])
		}
	}

	paramInputList := generateCombinations(paramInputMap)

	paramInputKeyValuePairList := make([]Pair, 0)
	for _, inputMap := range paramInputList {
		paramTmp := &Pair{
			KeyValuePair: make([]*common.KeyValuePair, 0),
			KeyValue:     make(map[string]interface{}),
		}

		for inputKey, inputValue := range inputMap {

			inputValueBytes, err := utils.MarshalInterfaceToBytes(inputValue)

			// fmt.Println(inputKey, inputValue, inputValueBytes)

			if err != nil {
				fmt.Println("marshal inputValue error!", err)
				return nil
			}

			paramTmp.KeyValuePair = append(paramTmp.KeyValuePair, &common.KeyValuePair{
				Key:   inputKey,
				Value: inputValueBytes,
			})

			paramTmp.KeyValue[inputKey] = inputValue
		}

		paramInputKeyValuePairList = append(paramInputKeyValuePairList, *paramTmp)

	}

	// for _, result := range paramInputList {
	// 	fmt.Println(result)
	// }

	return paramInputKeyValuePairList

}

// 返回匹配数最小的FuncName，及其对应的候选输入
// 返回 "", nil 时，代表已无需进行确认
func getMinProductPairInputList() (string, []Pair) {
	funcName, product := calMinProduct()

	if product < 0 {
		return "", nil
	}

	return funcName, generateInput(funcName)

}

// 将成功的输入加入confirmParam列表中
func addToConfirmParam(confirmValue []interface{}, input interface{}) []interface{} {
	inputType := reflect.TypeOf(input).Kind()

	switch inputType {
	// 基础类型,只需要一个该类型的输入
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
		for _, value := range confirmValue {
			a, _ := json.Marshal(value)
			b, _ := json.Marshal(input)
			if cmp.Equal(a, b) {
				return confirmValue
			}
		}
	// 其他未识别类型，需要值完全相同
	default:
		for _, value := range confirmValue {
			if cmp.Equal(value, input) {
				return confirmValue
			}
		}
	}

	return append(confirmValue, input)
}

// 将对应的函数设置为已确认的数据
func ConfirmParam(KeyValue map[string]interface{}) {
	for key, value := range KeyValue {
		// fmt.Println("设置对应的输入：", key, value)
		utils.GlobalContractInfo.ParamAndCandidateTypes[key].Confirm = true
		utils.GlobalContractInfo.ParamAndCandidateTypes[key].ConfirmValue = addToConfirmParam(utils.GlobalContractInfo.ParamAndCandidateTypes[key].ConfirmValue, value)
	}
	/*
		// 根据函数名获取输入参数
		paramNameList := utils.GlobalFuncAndParamMap[funcName].Params

		// 该函数不需要输入参数
		if len(paramNameList) == 0 {
			return
		}

		//根据输入参数修改confirm
		for _, paramName := range paramNameList {
			utils.GlobalParamAndCandidateTypes[paramName].Confirm = true
			utils.GlobalParamAndCandidateTypes[paramName].ConfirmValue = value
		}
	*/
}

/*
	1. 得到匹配数最小的待测函数，获取输入候选集

	2. 依次对候选输入进行测试，获取成功返回success的输入

	3. 将对应函数设置为已确认数据
*/
// 确定所有输入的类型，如果所有候选类型都不满足，则设置为string
func ConfirmAllInputParamType() {
	Log := utils.Log
	cnt := 0

	for ; ; cnt++ {
		funcName, inputList := getMinProductPairInputList()
		// 目前已没有需要确认的函数，退出
		if funcName == "" && inputList == nil {
			Log.Log(utils.ExecutionLog, "we finished comfirm all input param type!")
			break
		}

		Log.Log(utils.ExecutionLog, fmt.Sprintf("[%d]开始对以下函数进行输入读写集测试:", cnt))
		Log.Log(utils.ExecutionLog, fmt.Sprintf("[%d]function name  :%s", cnt, funcName))
		Log.Log(utils.ExecutionLog, fmt.Sprintf("[%d]input list     :\n %v", cnt, inputList))
		Log.Log(utils.ExecutionLog, fmt.Sprintf("[%d]input list len :[%d]", cnt, len(inputList)))

		flag := false

		for i, input := range inputList {
			Log.Log(utils.ExecutionLog, fmt.Sprintf("[%d]执行第[%d]轮,     input : %v", cnt, i, input))
			_, _, _, success, err := nodecontrol.ChainmakerController.UserContractInvoke(utils.GlobalContractInfo.ContractName, funcName, input.KeyValuePair, true)
			Log.Log(utils.ExecutionLog, fmt.Sprintf("[%d]执行成功与否:[%v], err  : %v", cnt, success, err))

			if err == nil && success {
				Log.Log(utils.ExecutionLog, fmt.Sprintf("[%d]confirm success! funcName:[%s]/KeyValue:[%s]", cnt, funcName, input.KeyValuePair))
				ConfirmParam(input.KeyValue)
				flag = true
			}
		}

		// flag没有发生改变，代表没有找到符合需求的输入
		// 将对应func下的param全置为string
		if !flag {
			Log.Log(utils.ExecutionLog, fmt.Sprintf("[%d]we can't find param type, and we will set all param as string", cnt))
			params := utils.GlobalContractInfo.ContractFuncMap[funcName].ParamsNameList
			for _, param := range params {
				utils.GlobalContractInfo.ParamAndCandidateTypes[param].Confirm = true
				utils.GlobalContractInfo.ParamAndCandidateTypes[param].ConfirmValue = addToConfirmParam(utils.GlobalContractInfo.ParamAndCandidateTypes[param].ConfirmValue, "string")
			}
		}

	}
}
