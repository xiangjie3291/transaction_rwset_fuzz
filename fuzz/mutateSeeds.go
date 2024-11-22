package fuzz

// // 已有functionName，生成对应的函数及其参数
// func initFuncParamsByName(functionName string, funcNameParamTable *utils.ContractFuncMap) *utils.FunctionAndParams {
// 	functionParamList := (*funcNameParamTable)[functionName]
// 	parameters := make([]*utils.KeyValuePair, 0)

// 	for _, param := range functionParamList {

// 		parameter := &utils.KeyValuePair{
// 			Key:   param.Name,
// 			Value: generateRandomValue(param.ActualType),
// 			// ReadRelated:  false,
// 			// WriteRealted: false,
// 		}
// 		parameters = append(parameters, parameter)
// 	}

// 	return &utils.FunctionAndParams{
// 		FunctionName: functionName,
// 		Parameters:   parameters,
// 	}
// }

// // 随机从函数名列表中选择两个函数名，并生成对应param的实际值
// func initFunctionAndParams(funcNameList []string, funcNameParamTable *utils.ContractFuncMap) *utils.FunctionAndParams {
// 	f := fuzz.New()
// 	var functionChoice uint
// 	f.Fuzz(&functionChoice)

// 	functionName := funcNameList[functionChoice%uint(len(funcNameList))]

// 	return initFuncParamsByName(functionName, funcNameParamTable)
// }

// // 生成一个初始种子，必须在函数名List、合约函数全局变量生成后才能调用
// func InitSeed() *utils.Seed {
// 	return &utils.Seed{
// 		ContractName:  "",
// 		FunctionOne:   initFunctionAndParams(utils.GlobalFuncNameList, &utils.GlobalFuncAndParamMap),
// 		FunctionTwo:   initFunctionAndParams(utils.GlobalFuncNameList, &utils.GlobalFuncAndParamMap),
// 		MaxSimilarity: 0.00,
// 		Mutability:    true,
// 	}
// }

// func InitSeedsPool() *list.List {
// 	l := list.New()

// 	// 对所有函数都构建一个种子
// 	for i := 0; i < len(utils.GlobalFuncNameList); i++ {
// 		funcNameA := utils.GlobalFuncNameList[i]

// 		for j := i; j < len(utils.GlobalFuncNameList); j++ {
// 			funcNameB := utils.GlobalFuncNameList[j]

// 			l.PushBack(&utils.Seed{
// 				ContractName:  "",
// 				FunctionOne:   initFuncParamsByName(funcNameA, &utils.GlobalFuncAndParamMap),
// 				FunctionTwo:   initFuncParamsByName(funcNameB, &utils.GlobalFuncAndParamMap),
// 				MaxSimilarity: 0.00,
// 				Mutability:    true,
// 			})
// 		}
// 	}

// 	return l
// }

// // 根据seed生成判断依赖关系的一组seeds
// func InitSeedList(initialSeed *utils.Seed) []*utils.Seed {
// 	seedList := []*utils.Seed{initialSeed}

// 	functionOneParams := initialSeed.FunctionOne.Parameters
// 	functionTwoParams := initialSeed.FunctionTwo.Parameters

// 	if len(functionOneParams) > 0 {
// 		for paramIndex := range functionOneParams {
// 			var newSeed *utils.Seed

// 			if err := utils.DeepCopy(initialSeed, &newSeed); err != nil {
// 				fmt.Println("Error:", err)
// 			}

// 			newSeed.FunctionOne.Parameters[paramIndex].Value = generateDiffValue(newSeed.FunctionOne.Parameters[paramIndex].Value)
// 			seedList = append(seedList, newSeed)
// 		}
// 	}

// 	if len(functionTwoParams) > 0 {
// 		for paramIndex := range functionTwoParams {
// 			var newSeed *utils.Seed

// 			if err := utils.DeepCopy(initialSeed, &newSeed); err != nil {
// 				fmt.Println("Error:", err)
// 			}
// 			newSeed.FunctionTwo.Parameters[paramIndex].Value = generateDiffValue(newSeed.FunctionTwo.Parameters[paramIndex].Value)
// 			seedList = append(seedList, newSeed)
// 		}
// 	}

// 	perceiveParamDependency(seedList)

// 	//debug:
// 	// for _, seed := range seedList {
// 	// 	fmt.Println("seed: ", seed.FunctionOne, seed.FunctionTwo)
// 	// }

// 	return seedList
// }
