package engine

import (
	"TransactionRwset/fuzz"
	getfuncInfo "TransactionRwset/info"
	nodecontrol "TransactionRwset/nodeControl"
	"fmt"
	"strings"
	"time"

	"TransactionRwset/utils"
)

func GetContractInfoAndPrepare(contractPath string) {
	utils.GlobalContractInfo = getfuncInfo.MakeGlobalContractInfo(contractPath)

	var err error
	utils.Log, err = utils.NewLogger(utils.GlobalContractInfo.ContractName)
	if err != nil {
		fmt.Printf("Error writing conflict log: %v\n", err)
	}

	var funcNameList []string

	for key, _ := range utils.GlobalContractInfo.ContractFuncMap {
		funcNameList = append(funcNameList, key)
	}

	utils.Log.Log(utils.ExecutionLog, "==================================  获取合约信息  ======================================")
	utils.Log.Log(utils.ExecutionLog, "	ContractName          : "+utils.GlobalContractInfo.ContractName)
	utils.Log.Log(utils.ExecutionLog, "	ContractByteCodePath  : "+utils.GlobalContractInfo.ContractByteCodePath)
	utils.Log.Log(utils.ExecutionLog, "	ContractFuncList      : "+strings.Join(funcNameList, " "))
	utils.Log.Log(utils.ExecutionLog, "	ParamAndCandidateTypes:\n"+fmt.Sprint(utils.GlobalContractInfo.ParamAndCandidateTypes))
	utils.Log.Log(utils.ExecutionLog, "=======================================================================================")
}

// 启动节点，并进行合约信息初步获取工作
func Start(contractPath string) {
	GetContractInfoAndPrepare(contractPath)
	Log := utils.Log

	// 启动节点
	nodecontrol.ChainmakerController.StartChainmaker()
	Log.Log(utils.ExecutionLog, "====================================  部署合约  ========================================")
	txId, err := nodecontrol.ChainmakerController.UserContractClaimCreate(utils.GlobalContractInfo.ContractName, utils.GlobalContractInfo.ContractByteCodePath, true, false)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(5 * time.Second)

	tx := nodecontrol.TestContractGetTxByTxId(txId)
	Log.Log(utils.ExecutionLog, "Claim Txid: "+txId)
	Log.Log(utils.ExecutionLog, fmt.Sprintf("result code:%d, msg:%s\n", tx.Transaction.Result.Code, tx.Transaction.Result.Code.String()))
	Log.Log(utils.ExecutionLog, "=======================================================================================")
}

// 交易种子、交易对种子生成工作
func GenerateSeeds() *fuzz.FuncPairSeedsPool {
	Log := utils.Log

	Log.Log(utils.ExecutionLog, "===========================  确认所有参数的可能使用的类型  ==============================")
	fuzz.ConfirmAllInputParamType()
	Log.Log(utils.ExecutionLog, "=======================================================================================")
	Log.Log(utils.ExecutionLog, "===========================  打印所有参数的可能使用的类型  ==============================")
	Log.Log(utils.ExecutionLog, "	ParamAndCandidateTypes: "+fmt.Sprint(utils.GlobalContractInfo.ParamAndCandidateTypes))
	Log.Log(utils.ExecutionLog, "=======================================================================================")
	Log.Log(utils.ExecutionLog, "==================================  生成种子池  ========================================")
	funcSeedsPool := fuzz.NewFuncSeedsPool()
	// 打印种子池结果（将其打入一个文件中）
	funcSeedsPool.PrintFuncSeedsPool()
	Log.Log(utils.ExecutionLog, "=======================================================================================")
	Log.Log(utils.ExecutionLog, "===============================  生成交易对种子池  ======================================")
	funcPairSeedsPool := fuzz.NewFuncPairSeedsPool(funcSeedsPool)
	// 打印交易对种子池
	funcPairSeedsPool.PrintFuncPairSeedsPool()

	// 将对应的交易对种子存入json文件中
	err := funcPairSeedsPool.SaveToFile(Log.BaseDir)
	if err != nil {
		fmt.Println(err)
	}
	Log.Log(utils.ExecutionLog, "=======================================================================================")
	return funcPairSeedsPool
}

// 冲突交易对种子进行实验
// 可变异种子进行变异
func HandleFuncPairSeedsPool(pool *fuzz.FuncPairSeedsPool) {
	Log := utils.Log
	Log.Log(utils.ExecutionLog, "=============================  冲突交易集测试/交易对变异  ===============================")
	Log.Log(utils.ExecutionLog, "=======================================================================================")

	round := 0
	for pool.ConflictSeeds.Len() > 0 || pool.MutateSeeds.Len() > 0 {
		if pool.ConflictSeeds.Len() > 0 {
			Log.Log(utils.ExecutionLog, fmt.Sprintf("round:[%d] 使用冲突交易对种子进行测试", round))
			Log.Log(utils.ConflictLog, fmt.Sprintf("=======================================  round:[%d]  =======================================", round))
			pool.ConflictTxsFirstSeedInPool()
			Log.Log(utils.ConflictLog, "============================================================================================")
			Log.Log(utils.ExecutionLog, fmt.Sprintf("round:[%d] 冲突交易对种子测试结束！", round))
			continue
		}

		if pool.MutateSeeds.Len() > 0 {
			Log.Log(utils.ExecutionLog, fmt.Sprintf("round:[%d] 对可变异种子进行变异", round))
			Log.Log(utils.FuzzLog, fmt.Sprintf("=======================================  round:[%d]  =======================================", round))
			pool.MutateFirstSeedInPool()
			Log.Log(utils.FuzzLog, "============================================================================================")
			Log.Log(utils.ExecutionLog, fmt.Sprintf("round:[%d] 可变异种子本轮变异结束！", round))
		}
	}

}

// 程序结束
func Stop() {
	nodecontrol.ChainmakerController.StopChainmaker()
}
