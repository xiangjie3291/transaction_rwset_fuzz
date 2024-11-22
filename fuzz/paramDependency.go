/*
	本文件主要用于：

	1.将interface{}(或map[string]interface{})转化为后续可求取叶子节点路径、读写集依赖相关信息的结构
	   该结构下记录：
			a. 每个参数对应的路径
			b. 造成读集改变的参数的路径
			c. 造成写集改变的参数的路径

	2.对每个函数，有自己的种子池：
			a. 对应的param
			b. 每个param有对应的confirmValue
			c. 每个confirmValue有自己的path
		保留最终会对读写集造成影响的种子

	3.对每个函数对，有自己的种子池：
			a. 由函数A的某个种子与函数B的某个种子组成
			b. 两个种子，其中一个种子的读集必须可变，另一个种子的写集必须可变，否则没有变异靠近的必要（或者直接能发现读写集冲突）
			c. 从种子池中选取种子，对对应path进行变异
*/

package fuzz

import (
	nodecontrol "TransactionRwset/nodeControl"
	"TransactionRwset/utils"
	"container/list"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/copystructure"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

// listToSlice 将 *list.List 转换为切片
func listToSlice(l *list.List) []*FuncPairSeed {
	var result []*FuncPairSeed
	for e := l.Front(); e != nil; e = e.Next() {
		result = append(result, e.Value.(*FuncPairSeed))
	}
	return result
}

// sliceToList 将切片转换为 *list.List
func sliceToList(slice []*FuncPairSeed) *list.List {
	l := list.New()
	for _, item := range slice {
		l.PushBack(item)
	}
	return l
}

type ValuePath []string

type FuncPairSeed struct {
	SeedOne       *FuncSeed `json:"seed_one"`
	SeedTwo       *FuncSeed `json:"seed_two"`
	MaxSimilarity float64   `json:"max_similarity"`
	Mutability    bool      `json:"mutability"`
}

func (f *FuncPairSeed) String() string {
	return fmt.Sprintf(`FuncPairSeed{
		SeedOne: %s,
		SeedTwo: %s,
		MaxSimilarity: %.2f,
		Mutability: %v
	}`,
		f.SeedOne,
		f.SeedTwo,
		f.MaxSimilarity,
		f.Mutability,
	)
}

type FuncPairSeedsPool struct {
	ConflictSeeds *list.List      `json:"-"` // 不直接序列化，使用辅助字段
	MutateSeeds   *list.List      `json:"-"`
	ConflictList  []*FuncPairSeed `json:"conflict_seeds"` // 用于序列化
	MutateList    []*FuncPairSeed `json:"mutate_seeds"`
}

func (f *FuncPairSeedsPool) PrintFuncPairSeedsPool() {
	Log := utils.Log
	conflictList := f.ConflictSeeds
	testList := f.MutateSeeds

	Log.Log(utils.ConflictLog, fmt.Sprintf("当前冲突交易对: [%d], 可变异交易对: [%d]", conflictList.Len(), testList.Len()))

	Log.Log(utils.ConflictLog, "------------------conflictList------------------")
	for e := conflictList.Front(); e != nil; e = e.Next() {
		Log.Log(utils.ConflictLog, fmt.Sprint(e.Value.(*FuncPairSeed)))
	}

	Log.Log(utils.ConflictLog, "-------------------mutateList-------------------")

	for e := testList.Front(); e != nil; e = e.Next() {
		Log.Log(utils.ConflictLog, fmt.Sprint(e.Value.(*FuncPairSeed)))
	}

}

// 保存到文件
func (pool *FuncPairSeedsPool) SaveToFile(baseDir string) error {
	timestamp := time.Now().Format("20060102_150405")

	// 将 *list.List 转换为切片
	pool.ConflictList = listToSlice(pool.ConflictSeeds)
	pool.MutateList = listToSlice(pool.MutateSeeds)

	// 序列化为 JSON
	data, err := json.MarshalIndent(pool, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize FuncPairSeedsPool: %v", err)
	}
	filePath := filepath.Join(baseDir, fmt.Sprintf("func_pair_seeds_pool_%s_%s.json", utils.GlobalContractInfo.ContractName, timestamp))

	// 写入文件
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

// 从文件读取
func LoadPairSeedPoolFromFile(filename string) (*FuncPairSeedsPool, error) {
	var pool *FuncPairSeedsPool
	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	err = json.Unmarshal(data, &pool)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize FuncPairSeedsPool: %v", err)
	}

	// 将切片还原为 *list.List
	pool.ConflictSeeds = sliceToList(pool.ConflictList)
	pool.MutateSeeds = sliceToList(pool.MutateList)

	return pool, nil
}

func (f *FuncPairSeedsPool) ConflictTxsFirstSeedInPool() {
	Log := utils.Log
	e := f.ConflictSeeds.Front()

	if e == nil {
		Log.Log(utils.ConflictLog, "we don't have seed to generate conflict txs!")
		return
	}

	seed := e.Value.(*FuncPairSeed)
	Log.Log(utils.ConflictLog, fmt.Sprintf("we will start use this seed:%s", seed))

	ConflictPairSeedExperiment(seed)

	f.ConflictSeeds.Remove(e)
}

// 新增变异逻辑，对种子对下的种子字段进行随机变异
// 取出首个种子，变异10000次，每次变异后执行对比maxsimilarity
func (f *FuncPairSeedsPool) MutateFirstSeedInPool() {
	Log := utils.Log

	e := f.MutateSeeds.Front()

	if e == nil {
		Log.Log(utils.FuzzLog, "we don't have seed to mutate!")
		return
	}

	seed := e.Value.(*FuncPairSeed)
	rand.Seed(time.Now().UnixNano())

	Log.Log(utils.FuzzLog, fmt.Sprintf("we will start mutate this seed:%s", seed))
	for i := 0; i < 10000; i++ {
		// 深拷贝一个种子用于变异
		copy, err := copystructure.Copy(seed)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// 类型断言
		mutateSeed, ok := copy.(*FuncPairSeed)
		if !ok {
			fmt.Println("Type assertion failed")
			return
		}

		a := len(mutateSeed.SeedOne.ReadRelatedValuePaths)
		b := len(mutateSeed.SeedOne.WriteRelatedValuePaths)
		c := len(mutateSeed.SeedTwo.ReadRelatedValuePaths)
		d := len(mutateSeed.SeedTwo.WriteRelatedValuePaths)
		// 选取变异对象
		mutataParamLen := a + b + c + d
		// 该种子无需变异
		if mutataParamLen == 0 {
			return
		}

		r := rand.Intn(mutataParamLen)

		var path ValuePath
		switch {
		case r < a:
			path = mutateSeed.SeedOne.ReadRelatedValuePaths[r]
			mutateSeed.SeedOne.modifyField(mutateSeed.SeedOne.FunctionInput, path)
			mutateSeed.SeedOne.getRWSets()
		case r < a+b:
			path = mutateSeed.SeedOne.WriteRelatedValuePaths[r-a]
			mutateSeed.SeedOne.modifyField(mutateSeed.SeedOne.FunctionInput, path)
			mutateSeed.SeedOne.getRWSets()
		case r < a+b+c:
			path = mutateSeed.SeedTwo.ReadRelatedValuePaths[r-a-b]
			mutateSeed.SeedTwo.modifyField(mutateSeed.SeedTwo.FunctionInput, path)
			mutateSeed.SeedTwo.getRWSets()
		case r < a+b+c+d:
			path = mutateSeed.SeedTwo.WriteRelatedValuePaths[r-a-b-c]
			mutateSeed.SeedTwo.modifyField(mutateSeed.SeedTwo.FunctionInput, path)
			mutateSeed.SeedTwo.getRWSets()
		}

		mutateSeed.MaxSimilarity = calculateMaxSimilarity(mutateSeed.SeedOne, mutateSeed.SeedTwo)

		if mutateSeed.MaxSimilarity > 0.99 {
			f.MutateSeeds.Remove(e)
			f.ConflictSeeds.PushBack(mutateSeed)
			Log.Log(utils.FuzzLog, fmt.Sprintf("we find new conflict Seed!: %s", mutateSeed))
			return
		}

		if mutateSeed.MaxSimilarity > seed.MaxSimilarity {
			*seed = *mutateSeed
		}
	}
	Log.Log(utils.FuzzLog, fmt.Sprintf("we don't find confict seed in this round:%s", seed))
	f.MutateSeeds.MoveToBack(e)
}

// Helper function to check if a seed pair already exists in a list based on function names
func containsFuncPairByName(list *list.List, seed *FuncPairSeed) bool {

	seedOne := seed.SeedOne
	seedTwo := seed.SeedTwo

	// Generate a unique, order-independent identifier for the new pair
	newPair := map[string]struct{}{
		seedOne.FunctionName: {},
		seedTwo.FunctionName: {},
	}

	for e := list.Front(); e != nil; e = e.Next() {
		pair := e.Value.(*FuncPairSeed)

		// Generate the identifier for the existing pair in the list
		existingPair := map[string]struct{}{
			pair.SeedOne.FunctionName: {},
			pair.SeedTwo.FunctionName: {},
		}

		// Compare the two maps
		if len(newPair) == len(existingPair) {
			match := true
			for fn := range newPair {
				if _, exists := existingPair[fn]; !exists {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
	}
	return false
}

func NewFuncPairSeedsPool(funcSeedsPool *FuncSeedsPool) *FuncPairSeedsPool {
	Log := utils.Log
	funcSeedsList := make([]*FuncSeed, 0)
	pool := funcSeedsPool.Pool
	for _, list := range pool {
		funcSeedsList = append(funcSeedsList, list...)
	}

	funcPairSeedsPool := &FuncPairSeedsPool{
		ConflictSeeds: list.New(),
		MutateSeeds:   list.New(),
	}

	// Traverse each possible combination of funcSeed
	for i := 0; i < len(funcSeedsList); i++ {
		seedOne := funcSeedsList[i]
		for j := i; j < len(funcSeedsList); j++ {
			seedTwo := funcSeedsList[j]
			pairSeed := &FuncPairSeed{
				SeedOne:       seedOne,
				SeedTwo:       seedTwo,
				MaxSimilarity: 0.00,
				Mutability:    false,
			}

			// Calculate similarity
			pairSeed.MaxSimilarity = calculateMaxSimilarity(seedOne, seedTwo)

			// Calculate mutability
			pairSeed.Mutability = conflictPotential(seedOne, seedTwo)

			// Check if the pair already exists based on function names
			if pairSeed.MaxSimilarity > 0.99 {
				if !containsFuncPairByName(funcPairSeedsPool.ConflictSeeds, pairSeed) {
					funcPairSeedsPool.ConflictSeeds.PushBack(pairSeed)
				}
			} else if pairSeed.Mutability {
				if !containsFuncPairByName(funcPairSeedsPool.MutateSeeds, pairSeed) {
					funcPairSeedsPool.MutateSeeds.PushBack(pairSeed)
				}
			}
			Log.Log(utils.ExecutionLog, "生成交易对种子：")
			Log.Log(utils.ExecutionLog, fmt.Sprint(pairSeed))
		}
	}

	return funcPairSeedsPool
}

// 不可能存在读写冲突的种子，maxSimilarity将会被置为false
// 对当前种子进行筛选，A有读相关、B有写相关且A的读集与B的写集不能为空 || A有写相关，B有都相关且A的写集与B的读集不能为空
func conflictPotential(funcOneSeed, funcTwoSeed *FuncSeed) bool {
	functionOneReadRelated, functionOneWriteRelated, functionTwoReadRelated,
		functionTwoWriteRelated := false, false, false, false

	functionOneRead, functionOneWrite, functionTwoRead, functionTwoWrite := true, true, true, true
	// false, false, false, false

	functionOneReadRelated = len(funcOneSeed.ReadRelatedValuePaths) > 0
	functionOneWriteRelated = len(funcOneSeed.WriteRelatedValuePaths) > 0
	functionTwoReadRelated = len(funcTwoSeed.ReadRelatedValuePaths) > 0
	functionTwoWriteRelated = len(funcTwoSeed.WriteRelatedValuePaths) > 0
	// functionOneRead = len(funcOneSeed.ReadSet) > 0
	// functionOneWrite = len(funcOneSeed.WriteSet) > 0
	// functionTwoRead = len(funcTwoSeed.ReadSet) > 0
	// functionTwoWrite = len(funcTwoSeed.WriteSet) > 0

	return ((functionOneReadRelated && functionTwoWriteRelated && functionOneRead && functionTwoWrite) ||
		(functionOneWriteRelated && functionTwoReadRelated && functionOneWrite && functionTwoRead))

}

// 计算该种子见读写集最大相似度
// 不负责执行交易！
func calculateMaxSimilarity(baseTxRwSet, compareTxRwSet *FuncSeed) float64 {
	maxSimilarity := 0.00

	if len(baseTxRwSet.ReadSet) != 0 {
		if len(compareTxRwSet.WriteSet) != 0 {
			for _, read := range baseTxRwSet.ReadSet {
				for _, write := range compareTxRwSet.WriteSet {
					maxSimilarity = utils.Max(maxSimilarity, utils.CalculateSimilarity(read, write))
				}
			}
		}
	}

	if len(baseTxRwSet.WriteSet) != 0 {
		if len(compareTxRwSet.ReadSet) != 0 {
			for _, write := range baseTxRwSet.WriteSet {
				for _, read := range compareTxRwSet.ReadSet {
					maxSimilarity = utils.Max(maxSimilarity, utils.CalculateSimilarity(write, read))
				}
			}
		}
	}

	return maxSimilarity
}

type FuncSeedsPool struct {
	// 函数名 - 对应函数名下的种子
	Pool map[string][]*FuncSeed
}

func (f *FuncSeedsPool) PrintFuncSeedsPool() {
	for key, value := range f.Pool {
		utils.Log.Log(utils.FuncSeedsLog, fmt.Sprintf("-------------------------  FuncName: [%s]  -------------------------", key))
		for _, seed := range value {
			utils.Log.Log(utils.FuncSeedsLog, fmt.Sprint(seed))
		}
	}
}

// 为所有funcName都生成对应的FuncSeed
func NewFuncSeedsPool() *FuncSeedsPool {
	funcSeedsPool := &FuncSeedsPool{
		Pool: make(map[string][]*FuncSeed),
	}

	for funcName := range utils.GlobalContractInfo.ContractFuncMap {
		funcSeedsPool.Pool[funcName] = generateNewFuncSeedList(funcName)
	}

	return funcSeedsPool
}

type FuncSeed struct {
	FunctionName           string                 `json:"function_name"`
	FunctionInput          map[string]interface{} `json:"function_input"`
	ValuePaths             []ValuePath            `json:"value_paths"`
	ReadRelatedValuePaths  []ValuePath            `json:"read_related_value_paths"`
	WriteRelatedValuePaths []ValuePath            `json:"write_related_value_paths"`
	ReadSet                []string               `json:"read_set"`
	WriteSet               []string               `json:"write_set"`
}

func (f *FuncSeed) String() string {
	return fmt.Sprintf(
		`	
			FuncSeed{
				FunctionName: %s,
				FunctionInput: %v,
				ValuePaths: %v,
				ReadRelatedValuePaths: %v,
				WriteRelatedValuePaths: %v,
				ReadSet: %v,
				WriteSet: %v
			}`,
		f.FunctionName,
		f.FunctionInput,
		f.ValuePaths,
		f.ReadRelatedValuePaths,
		f.WriteRelatedValuePaths,
		f.ReadSet,
		f.WriteSet,
	)
}

// 为某个Func生成所有param对应confirmValue所生成的FuncSeed
// 1. 对每个FuncSeed生成叶子节点路径
// 2. 获取该输入下读写集
// 3. 计算每个读、写相关路径
func generateNewFuncSeedList(funcName string) []*FuncSeed {
	Log := utils.Log

	paramsName := utils.GlobalContractInfo.ContractFuncMap[funcName].ParamsNameList

	/*
		生成所有param下不同confirmValue构成的组合
	*/
	var result []map[string]interface{}
	// 递归生成组合
	var helper func(int, map[string]interface{})
	helper = func(index int, current map[string]interface{}) {
		if index == len(paramsName) {
			// 拷贝当前组合并添加到结果集
			combination := make(map[string]interface{})
			for k, v := range current {
				combination[k] = v
			}
			result = append(result, combination)
			return
		}

		// 遍历当前参数的候选值
		key := paramsName[index]
		for _, value := range utils.GlobalContractInfo.ParamAndCandidateTypes[key].ConfirmValue {
			current[key] = value
			helper(index+1, current)
		}
	}

	// 调用递归函数
	helper(0, make(map[string]interface{}))

	/*
		对每个组合构建seed
	*/

	var newFuncSeedList []*FuncSeed

	for _, functionInput := range result {
		funcSeed := &FuncSeed{
			FunctionName:           funcName,
			FunctionInput:          functionInput,
			ValuePaths:             make([]ValuePath, 0),
			ReadRelatedValuePaths:  make([]ValuePath, 0),
			WriteRelatedValuePaths: make([]ValuePath, 0),
			ReadSet:                make([]string, 0),
			WriteSet:               make([]string, 0),
		}

		// 获取每个变量的path
		funcSeed.getValuePaths(functionInput, []string{}, &funcSeed.ValuePaths)

		// 获取该funcSeed读写集
		funcSeed.getRWSets()

		// 获取该funcSeed读写集变化相关变量
		funcSeed.getRelatedValuePaths()

		// deubg:
		Log.Log(utils.ExecutionLog, "生成交易种子：")
		Log.Log(utils.ExecutionLog, fmt.Sprint(funcSeed))

		newFuncSeedList = append(newFuncSeedList, funcSeed)
	}

	return newFuncSeedList
}

// 将一个map下的叶子节点路径打印出来
func (f *FuncSeed) getValuePaths(data interface{}, parentPath ValuePath, paths *[]ValuePath) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			currentPath := append(append((make([]string, 0)), parentPath...), key)
			f.getValuePaths(value, currentPath, paths)
		}
	case []interface{}:
		for index, value := range v {
			// 使用索引作为路径的一部分
			indexKey := fmt.Sprintf("[%d]", index)
			currentPath := append(append((make([]string, 0)), parentPath...), indexKey)
			f.getValuePaths(value, currentPath, paths)
		}
	default:
		*paths = append(*paths, parentPath)
	}
}

// 将一个seed的运行结果转化为string类型的读写集合
func (f *FuncSeed) convertRwSetToStringList(txTwSet *common.TxRWSet) ([]string, []string) {

	ReadSet := []string{}
	WriteSet := []string{}

	if txTwSet != nil && txTwSet.TxReads != nil {
		for _, txRead := range txTwSet.TxReads {
			ReadSet = append(ReadSet, string(txRead.Key))
		}
	}

	if txTwSet != nil && txTwSet.TxWrites != nil {
		for _, txWrite := range txTwSet.TxWrites {
			WriteSet = append(WriteSet, string(txWrite.Key))
		}
	}

	return ReadSet, WriteSet
}

// 将一个map[string]interface{}类型转化为key value pair类型
func (f *FuncSeed) convertMapToKeyValuePair(input map[string]interface{}) []*common.KeyValuePair {
	KeyValuePair := make([]*common.KeyValuePair, 0)

	for inputKey, inputValue := range input {

		inputValueBytes, err := utils.MarshalInterfaceToBytes(inputValue)

		if err != nil {
			fmt.Println("marshal inputValue error!", err)
			return nil
		}

		KeyValuePair = append(KeyValuePair, &common.KeyValuePair{
			Key:   inputKey,
			Value: inputValueBytes,
		})
	}
	return KeyValuePair

}

// 获取读写集
// 提交交易执行请求
// 根据交易id进行查询
func (f *FuncSeed) getRWSets() {
	keyValuePair := f.convertMapToKeyValuePair(f.FunctionInput)

	txid, _, _, _, _ := nodecontrol.ChainmakerController.UserContractInvoke(utils.GlobalContractInfo.ContractName, f.FunctionName, keyValuePair, true)

	txInfo, _ := nodecontrol.ChainmakerController.Client.GetTxWithRWSetByTxId(txid)

	f.ReadSet, f.WriteSet = f.convertRwSetToStringList(txInfo.RwSet)
}

// 获取读写相关变量
/*
	1. 遍历每一条path

	2. 在原始input上根据path进行改动并执行交易

	3. 计算变异后与变异前的读写集，判断是否发生变化

	4. 将path加入对应relatedValuePaths
*/
func (f *FuncSeed) getRelatedValuePaths() {
	// 1. 遍历path
	for _, path := range f.ValuePaths {
		// 2. 根据path对input进行变异
		copy, err := copystructure.Copy(f.FunctionInput)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// 类型断言
		mutateInput, ok := copy.(map[string]interface{})
		if !ok {
			fmt.Println("Type assertion failed")
			return
		}

		f.modifyField(mutateInput, path)

		mutateKeyValuePair := f.convertMapToKeyValuePair(mutateInput)

		txid, _, _, _, _ := nodecontrol.ChainmakerController.UserContractInvoke(utils.GlobalContractInfo.ContractName, f.FunctionName, mutateKeyValuePair, true)

		txInfo, _ := nodecontrol.ChainmakerController.Client.GetTxWithRWSetByTxId(txid)

		// 3. 计算读写集差异
		ReadSet, WriteSet := f.convertRwSetToStringList(txInfo.RwSet)

		if !utils.StringArraysEqual(f.ReadSet, ReadSet) {
			// 4. 加入realatedPath
			f.ReadRelatedValuePaths = append(f.ReadRelatedValuePaths, path)
		}

		if !utils.StringArraysEqual(f.WriteSet, WriteSet) {
			// 4. 加入realatedPath
			f.WriteRelatedValuePaths = append(f.WriteRelatedValuePaths, path)
		}

	}

}

// 修改字段值的函数
func (f *FuncSeed) modifyField(data interface{}, path ValuePath) {
	if len(path) == 0 {
		return
	}
	switch v := data.(type) {
	case map[string]interface{}:
		key := path[0]
		if len(path) == 1 {
			v[key] = utils.GenerateDiffValue(v[key])
		} else {
			if nextData, exists := v[key]; exists {
				f.modifyField(nextData, path[1:])
			} else {
				// 错误处理：路径不存在
			}
		}
	case []interface{}:
		indexStr := path[0]
		if strings.HasPrefix(indexStr, "[") && strings.HasSuffix(indexStr, "]") {
			index, err := strconv.Atoi(indexStr[1 : len(indexStr)-1])
			if err != nil {
				// 错误处理：索引解析失败
				return
			}
			if index < 0 || index >= len(v) {
				// 错误处理：索引越界
				return
			}
			if len(path) == 1 {
				v[index] = utils.GenerateDiffValue(v[index])
			} else {
				f.modifyField(v[index], path[1:])
			}
		} else {
			// 错误处理：无效的索引格式
		}
	default:
		// 错误处理：非预期的类型
	}
}

// // 用于判断两个input之间是否存在不同的参数（理论上只会有一个参数不同）
// // 为空时不产生任何影响
// func inputParamsEqual(baseInput *utils.FunctionAndParams, compareInput *utils.FunctionAndParams) string {
// 	if baseInput.FunctionName != compareInput.FunctionName {
// 		fmt.Println("have different function!")
// 		return ""
// 	}

// 	if len(baseInput.Parameters) != len(compareInput.Parameters) {
// 		fmt.Println("have different params!")
// 		return ""
// 	}

// 	for i := range baseInput.Parameters {
// 		if baseInput.Parameters[i].Key != compareInput.Parameters[i].Key {
// 			fmt.Println("have different key!")
// 			return ""
// 		}

// 		if baseInput.Parameters[i].Value != compareInput.Parameters[i].Value {
// 			return baseInput.Parameters[i].Key
// 		}
// 	}
// 	return ""
// }

// // 用于运行一组seeds，获取对应读写集、实际参数列表
// // result: 交易一实际参数读写集List，交易二实际参数读写集List
// func runSeedsAndGetRwSetList(compareSeedsList []*utils.Seed) ([]*utils.TxRWSetString, []*utils.TxRWSetString) {
// 	functionOneResultList := make([]*utils.TxRWSetString, 0)
// 	functionTwoResultList := make([]*utils.TxRWSetString, 0)
// 	for seed := range compareSeedsList {
// 		GenerateExecutionFile(compareSeedsList[seed])
// 		functionOneResult, functionTwoResult, maxSimilarity := utils.RunMainFile()

// 		//debug
// 		fmt.Println(compareSeedsList[seed].FunctionOne.FunctionName)
// 		fmt.Println(compareSeedsList[seed].FunctionTwo.FunctionName)
// 		fmt.Println(functionOneResult)
// 		fmt.Println(functionTwoResult)
// 		fmt.Println(maxSimilarity)

// 		compareSeedsList[seed].MaxSimilarity = maxSimilarity
// 		functionOneResultList = append(functionOneResultList, convertRwSetToStringList(compareSeedsList[seed].FunctionOne, functionOneResult))
// 		functionTwoResultList = append(functionTwoResultList, convertRwSetToStringList(compareSeedsList[seed].FunctionTwo, functionTwoResult))
// 	}

// 	return functionOneResultList, functionTwoResultList
// }

// // 根据两个种子结果的读写集，判断变量是否因某一参数发生变化
// // 分别判断读集合、写集合的变化
// // 判断参数的变化
// // 修改对应的读写集依赖关系
// func paramDependency(base *utils.TxRWSetString, compare *utils.TxRWSetString) {
// 	readRelated := !stringArraysEqual(base.ReadSet, compare.ReadSet)

// 	// debug:
// 	// fmt.Println(base.ReadSet)
// 	// fmt.Println(compare.ReadSet)
// 	// fmt.Println(readRelated)

// 	writeRelated := !stringArraysEqual(base.WriteSet, compare.WriteSet)

// 	//debug:
// 	// fmt.Println(base.WriteSet)
// 	// fmt.Println(compare.WriteSet)
// 	// fmt.Println(writeRelated)

// 	diffParam := inputParamsEqual(base.Input, compare.Input)
// 	if diffParam != "" {
// 		for i, param := range utils.GlobalFuncAndParamMap[base.Input.FunctionName] {
// 			if param.Name == diffParam {
// 				utils.GlobalFuncAndParamMap[base.Input.FunctionName][i].ReadRelated = readRelated
// 				utils.GlobalFuncAndParamMap[base.Input.FunctionName][i].WriteRelated = writeRelated
// 			}
// 		}
// 	}
// }

// // 根据一组seeds，判断param与读写集之间是否存在依赖关系
// // 判断种子有效性，必须有functionOne & functionTwo 读与写为true才存在读写冲突
// func (compareSeedsList []*utils.Seed) {
// 	if len(compareSeedsList) <= 1 {
// 		fmt.Println("don't need compare!")
// 		return
// 	}

// 	funcOneList, funcTwoList := runSeedsAndGetRwSetList(compareSeedsList)

// 	for i := 1; i < len(funcOneList); i++ {
// 		paramDependency(funcOneList[0], funcOneList[i])
// 	}

// 	for i := 1; i < len(funcTwoList); i++ {
// 		paramDependency(funcTwoList[0], funcTwoList[i])
// 	}

// 	for _, seed := range compareSeedsList {
// 		conflictPotential(seed)

// 		//debug:
// 		//fmt.Println(seed.Mutability)
// 	}

// 	//debug:
// 	// for i := range utils.GlobalFuncAndParamMap {
// 	// 	fmt.Println("Func Name: ", i)
// 	// 	for _, param := range utils.GlobalFuncAndParamMap[i] {
// 	// 		fmt.Println(param.Name, param.ReadRelated, param.WriteRealted)
// 	// 	}
// 	// }
// }
