package getContractStaticInfo

import (
	"TransactionRwset/utils"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func MakeGlobalContractInfo(path string) *utils.ContractInfo {
	fset := token.NewFileSet()
	path, err := filepath.Abs(path)
	if err != nil {
		log.Println(err)
		return nil
	}

	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Println(err)
		return nil
	}

	info := &utils.ContractInfo{
		// 0. 获取合约Node
		Node: node,
		// 1. 获取合约路径
		ContractPath: path,
	}

	// 2. 获取合约目录和合约名
	generateContractDirAndName(info)

	// 3. 获取合约二进制文件路径
	info.ContractByteCodePath = buildContract(info.ContractDir, info.ContractName)

	// 4. 获取待测合约func和param相关信息
	generateContractFuncMap(info)

	// 5. 获取param与candidateTypes相关信息
	getParamCandidateTypeBySSA(info)

	return info
}

/*
*********************************

*	根据合约完整路径，获取：
* 	1. 合约目录
*	2. 合约名

*********************************
 */
func generateContractDirAndName(info *utils.ContractInfo) {
	path := info.ContractPath
	file := filepath.Base(path)

	ext := filepath.Ext(file)

	info.ContractDir = filepath.Dir(path)

	info.ContractName = strings.TrimSuffix(file, ext)

}

/*
*********************************

*	获取：
* 	1. 获取合约中对应的函数名
*	2. 从客户端发送交易时的调用名
*	3. 各函数对应的param函数名列表

*********************************
 */
func generateContractFuncMap(info *utils.ContractInfo) {
	info.ContractFuncMap = make(map[string]*utils.FuncAndParamsNameInfo)

	node := info.Node

	funcNames, invokeNames := GetTxsName(node)

	if len(funcNames) != len(invokeNames) {
		fmt.Println("Get Txs Name Error!")
		return
	}

	for i, funcName := range funcNames {
		// utils.GlobalFuncAndParamMap[txsNames[i]].InvokeName = invokeNames[i]
		infoTmp := &utils.FuncAndParamsNameInfo{
			InvokeName:     invokeNames[i],
			Product:        -1,
			ParamsNameList: GetFuncParams(funcName, node),
		}
		info.ContractFuncMap[funcName] = infoTmp
	}
}

/*
*********************************

*	获取：
* 	1. 合约下各个变量的candidateTypes

*********************************
 */
func getParamCandidateTypeBySSA(info *utils.ContractInfo) {
	info.ParamAndCandidateTypes = GetParamCandidateTypeBySSA(info.ContractPath)
}

/*
*********************************

*	编译合约
*	获取：
* 	1. 合约二进制文件路径

*********************************
 */
func buildContract(ContractDir string, ContractName string) string {
	scriptDir := ContractDir

	cmdStart := exec.Command("./build.sh", ContractName)
	cmdStart.Dir = scriptDir // 设置启动脚本的工作目录
	cmdStart.Stdout = os.Stdout
	cmdStart.Stderr = os.Stderr
	err := cmdStart.Run()
	if err != nil {
		fmt.Println("编译合约出错：", err)
		return ""
	}

	return ContractDir + "/" + ContractName + ".7z"
}
