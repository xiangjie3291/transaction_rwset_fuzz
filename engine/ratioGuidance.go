package engine

// func getRelatedParamName(funcAndParams *utils.FunctionAndParams) []string {
// 	result := make([]string, 0)

// 	for _, keyValuePair := range funcAndParams.Parameters {
// 		param, exist := utils.GlobalFuncAndParamMap[funcAndParams.FunctionName][keyValuePair.Key]
// 		if exist && (param.ReadRelated || param.WriteRelated) {
// 			result = append(result, param.Name)
// 		}
// 	}

// 	return result
// }

// // 随机从FunctionOne和FunctionTwo中选择一个变量进行变异
// // 如果是无关变量则重新选取，直到选择到有关变量
// // 判断种子有效性，Mutability为true才变异
// func MutateSeed(seed *utils.Seed) (*utils.Seed, int) {
// 	round := 0
// 	if !seed.Mutability {
// 		fmt.Println("this seed dont need mutate")
// 		return nil, round
// 	}

// 	funcOneRelatedParam := getRelatedParamName(seed.FunctionOne)
// 	funcTwoRelaredParam := getRelatedParamName(seed.FunctionTwo)

// 	lenSum := len(funcOneRelatedParam) + len(funcTwoRelaredParam)

// 	if lenSum == 0 {
// 		fmt.Println("this seed dont need mutate too!")
// 		return nil, round
// 	}

// 	for seed.MaxSimilarity < 0.99 {
// 		round++
// 		paramChoice := int(generateRandomValue(reflect.Uint16).(uint16)) % lenSum
// 		var mutateParamName string
// 		var newSeed *utils.Seed

// 		if err := utils.DeepCopy(seed, &newSeed); err != nil {
// 			fmt.Println("Error:", err)
// 		}

// 		if paramChoice+1 <= len(funcOneRelatedParam) {
// 			mutateParamName = funcOneRelatedParam[paramChoice]
// 			for _, keyValuePair := range newSeed.FunctionOne.Parameters {
// 				if keyValuePair.Key == mutateParamName {
// 					keyValuePair.Value = generateDiffValue(keyValuePair.Value)
// 				}
// 			}
// 		} else {
// 			mutateParamName = funcTwoRelaredParam[paramChoice-len(funcOneRelatedParam)]
// 			for _, keyValuePair := range newSeed.FunctionTwo.Parameters {
// 				if keyValuePair.Key == mutateParamName {
// 					keyValuePair.Value = generateDiffValue(keyValuePair.Value)
// 				}
// 			}
// 		}

// 		GenerateExecutionFile(newSeed)
// 		_, _, newSeed.MaxSimilarity = utils.RunMainFile()

// 		//debug:
// 		fmt.Println("mutate computate:", newSeed.MaxSimilarity)

// 		if newSeed.MaxSimilarity > seed.MaxSimilarity {
// 			seed = newSeed
// 		}
// 	}

// 	return seed, round

// 	// if paramChoice <= len(funcOneRelatedParam)-1 {
// 	// 	mutateParamName = funcOneRelatedParam[paramChoice]
// 	// } else {
// 	// 	if len(funcOneRelatedParam) == 0 {
// 	// 		mutateParamName = funcTwoRelaredParam[paramChoice]
// 	// 	}

// 	// 	paramChoice = paramChoice - ()
// 	// }

// }
